package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	httpDelivery "captcha-service/internal/transport/http"
)

func main() {
	cfg, err := config.LoadBalancerProxyConfig()
	if err != nil {
		log.Fatalf("Failed to load balancer-proxy config: %v", err)
	}

	entityConfig := &entity.Config{
		MaxAttempts:      cfg.MaxAttempts,
		BlockDurationMin: cfg.BlockDurationMin,
	}
	proxy := httpDelivery.NewBalancerProxy(entityConfig)

	balancerAddr := os.Getenv("BALANCER_ADDRESS")
	if balancerAddr == "" {
		balancerAddr = "localhost:9090"
	}
	if err := proxy.ConnectToBalancer(balancerAddr); err != nil {
		log.Fatalf("Failed to connect to balancer: %v", err)
	}

	go proxy.StartServiceDiscovery()

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			proxy.CleanupSessions()
		}
	}()

	mux := http.NewServeMux()

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			h.ServeHTTP(w, r)
		})
	}

	mux.Handle("/backgrounds/", corsHandler(http.StripPrefix("/backgrounds/", http.FileServer(http.Dir("./backgrounds/")))))
	mux.HandleFunc("/ws", proxy.WebSocketHandler)

	mux.HandleFunc("/challenge", proxy.ChallengeHandler)
	mux.HandleFunc("/api/challenge", proxy.ChallengeHandler)
	mux.HandleFunc("/api/validate", proxy.ValidateChallengeHandler)
	mux.HandleFunc("/api/services", proxy.ListServicesHandler)
	mux.HandleFunc("/api/services/add", proxy.AddServiceHandler)
	mux.HandleFunc("/api/services/remove", proxy.RemoveServiceHandler)
	mux.HandleFunc("/api/health", proxy.HealthHandler)

	mux.HandleFunc("/blocked", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "User blocked", http.StatusTooManyRequests)
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Printf("Balancer proxy started on http://localhost:8081")
		log.Printf("Open http://localhost:8082/demo in your browser")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

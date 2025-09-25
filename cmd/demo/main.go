package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"captcha-service/internal/config"
	"captcha-service/internal/domain/entity"
	"captcha-service/internal/infrastructure/cache"
	"captcha-service/internal/infrastructure/persistence"
	"captcha-service/internal/infrastructure/repository"
	"captcha-service/internal/service"
	httpTransport "captcha-service/internal/transport/http"
	wsTransport "captcha-service/internal/transport/websocket"
	"captcha-service/internal/usecase"
)

func main() {
	cfg, err := config.LoadDemoConfig()
	if err != nil {
		log.Fatalf("Failed to load demo config: %v", err)
	}

	sessionRepo := repository.NewInMemorySessionRepository()
	_ = cache.NewSessionCache(cfg.MaxSessions)
	demoUsecase := usecase.NewDemoUsecase(sessionRepo, cfg)

	entityConfigFromEnv, err := config.LoadCaptchaServiceConfig()
	if err != nil {
		log.Fatalf("Failed to load captcha config: %v", err)
	}
	entityConfig := entityConfigFromEnv

	challengeRepo := persistence.NewMemoryOptimizedRepository(cfg.MaxChallenges)

	registry := service.NewGeneratorRegistry()
	sliderGenerator := service.NewSliderPuzzleGenerator(entityConfig, challengeRepo, nil)
	registry.Register(entity.ChallengeTypeSliderPuzzle, sliderGenerator)

	captchaService := service.NewCaptchaService(challengeRepo, registry, nil, entityConfig)

	tmpl := template.New("demo")

	demoHandler := httpTransport.NewDemoHandler(demoUsecase, captchaService, tmpl, cfg)
	wsHandler := wsTransport.NewDemoWebSocketHandler(cfg, sessionRepo)

	router := httpTransport.NewRouter(demoHandler, wsHandler)
	mux := router.SetupRoutes()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down demo server...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.CleanupIntervalHours) * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			sessionRepo.CleanupExpired()
		}
	}()

	log.Printf("Demo server started on http://localhost:%s", cfg.Port)
	log.Printf("Available endpoints:")
	log.Printf("  - http://localhost:%s/demo - Main demo page", cfg.Port)
	log.Printf("  - http://localhost:%s/performance - Performance test", cfg.Port)
	log.Printf("  - http://localhost:%s/health - Health check", cfg.Port)
	log.Printf("  - ws://localhost:%s/ws - WebSocket endpoint", cfg.Port)
	log.Printf("Captcha service URL: %s", cfg.CaptchaServiceURL)

	log.Fatal(server.ListenAndServe())
}

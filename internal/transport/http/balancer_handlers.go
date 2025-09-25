package http

import (
	"encoding/json"
	"net/http"
	"time"

	"captcha-service/internal/service"
)

type BalancerHandlers struct {
	balancerService *service.BalancerService
}

func NewBalancerHandlers(balancerService *service.BalancerService) *BalancerHandlers {
	return &BalancerHandlers{
		balancerService: balancerService,
	}
}

func (h *BalancerHandlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"balancer"}`))
}

func (h *BalancerHandlers) APIHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"balancer"}`))
}

func (h *BalancerHandlers) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	instances, err := h.balancerService.GetInstances()
	if err != nil {
		http.Error(w, "Failed to get instances", http.StatusInternalServerError)
		return
	}

	services := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		services[i] = map[string]interface{}{
			"id":            instance.ID,
			"type":          instance.Type,
			"host":          instance.Host,
			"port":          instance.Port,
			"status":        instance.Status,
			"last_seen":     instance.LastSeen.Format(time.RFC3339),
			"registered_at": instance.RegisteredAt.Format(time.RFC3339),
		}
	}

	response := map[string]interface{}{
		"services":  services,
		"count":     len(services),
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

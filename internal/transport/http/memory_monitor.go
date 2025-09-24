package http

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"captcha-service/internal/infrastructure/cache"
	"captcha-service/internal/infrastructure/persistence"
	"captcha-service/internal/service"
)

type MemoryMonitor struct {
	challengeRepo *persistence.MemoryOptimizedRepository
	sessionCache  *cache.SessionCache
	globalBlocker *service.GlobalUserBlocker
}

func NewMemoryMonitor(
	challengeRepo *persistence.MemoryOptimizedRepository,
	sessionCache *cache.SessionCache,
	globalBlocker *service.GlobalUserBlocker,
) *MemoryMonitor {
	return &MemoryMonitor{
		challengeRepo: challengeRepo,
		sessionCache:  sessionCache,
		globalBlocker: globalBlocker,
	}
}

func (m *MemoryMonitor) HandleMemoryStats(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	challengeStats := m.challengeRepo.GetStats()
	sessionStats := m.sessionCache.GetStats()
	blockerStats := m.globalBlocker.GetStats()

	response := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"system": map[string]interface{}{
			"alloc_mb":        memStats.Alloc / 1024 / 1024,
			"total_alloc_mb":  memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":          memStats.Sys / 1024 / 1024,
			"num_gc":          memStats.NumGC,
			"gc_cpu_fraction": memStats.GCCPUFraction,
		},
		"challenges": challengeStats,
		"sessions":   sessionStats,
		"blocker":    blockerStats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m *MemoryMonitor) HandleMemoryGC(w http.ResponseWriter, r *http.Request) {
	before := runtime.MemStats{}
	runtime.ReadMemStats(&before)

	runtime.GC()

	after := runtime.MemStats{}
	runtime.ReadMemStats(&after)

	response := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"before_mb": before.Alloc / 1024 / 1024,
		"after_mb":  after.Alloc / 1024 / 1024,
		"freed_mb":  (before.Alloc - after.Alloc) / 1024 / 1024,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

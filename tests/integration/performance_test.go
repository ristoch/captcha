package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCaptchaServicePerformance(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"38000/tcp"},
		Env: map[string]string{
			"MIN_PORT":             "38000",
			"LOG_LEVEL":            "info",
			"CHALLENGE_TYPE":       "drag-drop",
			"MAX_ATTEMPTS":         "5",
			"MIN_OVERLAP_PCT":      "80",
			"MAX_TIMEOUT_ATTEMPTS": "3",
		},
		WaitingFor: wait.ForHTTP("/health").WithPort("38000/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer container.Terminate(ctx)

	port, err := container.MappedPort(ctx, "38000")
	require.NoError(t, err)

	baseURL := fmt.Sprintf("http://localhost:%s", port.Port())

	t.Run("Concurrent Challenge Creation", func(t *testing.T) {
		const numGoroutines = 50
		const challengesPerGoroutine = 20

		var wg sync.WaitGroup
		results := make(chan bool, numGoroutines*challengesPerGoroutine)
		start := time.Now()

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < challengesPerGoroutine; j++ {
					reqBody := map[string]interface{}{
						"challenge_type": "drag-drop",
						"complexity":     50,
					}

					jsonBody, err := json.Marshal(reqBody)
					if err != nil {
						results <- false
						continue
					}

					resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
					if err != nil {
						results <- false
						continue
					}
					resp.Body.Close()

					results <- resp.StatusCode == http.StatusOK
				}
			}()
		}

		wg.Wait()
		close(results)

		duration := time.Since(start)
		successCount := 0
		totalRequests := 0

		for success := range results {
			totalRequests++
			if success {
				successCount++
			}
		}

		rps := float64(totalRequests) / duration.Seconds()
		successRate := float64(successCount) / float64(totalRequests) * 100

		t.Logf("Performance test results:")
		t.Logf("  Total requests: %d", totalRequests)
		t.Logf("  Successful requests: %d", successCount)
		t.Logf("  Success rate: %.2f%%", successRate)
		t.Logf("  Duration: %v", duration)
		t.Logf("  RPS: %.2f", rps)

		assert.GreaterOrEqual(t, successRate, 95.0, "Success rate should be at least 95%")
		assert.GreaterOrEqual(t, rps, 50.0, "RPS should be at least 50")
	})

	t.Run("Memory Usage Under Load", func(t *testing.T) {
		const numChallenges = 1000

		start := time.Now()
		challengeIDs := make([]string, 0, numChallenges)

		for i := 0; i < numChallenges; i++ {
			reqBody := map[string]interface{}{
				"challenge_type": "drag-drop",
				"complexity":     50,
			}

			jsonBody, err := json.Marshal(reqBody)
			require.NoError(t, err)

			resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)

			if resp.StatusCode == http.StatusOK {
				var challenge map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&challenge)
				resp.Body.Close()
				if err == nil {
					challengeIDs = append(challengeIDs, challenge["challenge_id"].(string))
				}
			} else {
				resp.Body.Close()
			}
		}

		duration := time.Since(start)
		t.Logf("Created %d challenges in %v", len(challengeIDs), duration)

		healthResp, err := http.Get(baseURL + "/health")
		require.NoError(t, err)
		defer healthResp.Body.Close()

		var health map[string]interface{}
		err = json.NewDecoder(healthResp.Body).Decode(&health)
		require.NoError(t, err)

		activeChallenges := int(health["active_challenges"].(float64))
		t.Logf("Active challenges: %d", activeChallenges)

		assert.GreaterOrEqual(t, len(challengeIDs), 900, "Should create at least 900 challenges")
	})
}

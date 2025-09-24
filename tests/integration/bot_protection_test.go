package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestBotProtection(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"38000/tcp"},
		Env: map[string]string{
			"MIN_PORT":               "38000",
			"LOG_LEVEL":              "debug",
			"CHALLENGE_TYPE":         "drag-drop",
			"MAX_ATTEMPTS":           "3",
			"MIN_OVERLAP_PCT":        "80",
			"MAX_TIMEOUT_ATTEMPTS":   "2",
			"BLOCK_DURATION_MINUTES": "1",
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

	t.Run("Too Many Attempts Blocking", func(t *testing.T) {
		createReq := map[string]interface{}{
			"challenge_type": "drag-drop",
			"complexity":     50,
		}

		jsonBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var challenge map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&challenge)
		require.NoError(t, err)

		challengeID := challenge["challenge_id"].(string)

		wrongAnswer := map[string]interface{}{
			"challenge_id": challengeID,
			"answer": map[string]interface{}{
				"position": map[string]float64{
					"x": 10.0,
					"y": 10.0,
				},
			},
		}

		for i := 0; i < 4; i++ {
			wrongBody, err := json.Marshal(wrongAnswer)
			require.NoError(t, err)

			validateResp, err := http.Post(baseURL+"/api/validate", "application/json", bytes.NewBuffer(wrongBody))
			require.NoError(t, err)

			if i < 3 {
				assert.Equal(t, http.StatusOK, validateResp.StatusCode)
			} else {
				assert.Equal(t, http.StatusInternalServerError, validateResp.StatusCode)
			}

			validateResp.Body.Close()
		}
	})

	t.Run("Too Fast Completion", func(t *testing.T) {
		createReq := map[string]interface{}{
			"challenge_type": "drag-drop",
			"complexity":     50,
		}

		jsonBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var challenge map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&challenge)
		require.NoError(t, err)

		challengeID := challenge["challenge_id"].(string)

		time.Sleep(100 * time.Millisecond)

		answer := map[string]interface{}{
			"challenge_id": challengeID,
			"answer": map[string]interface{}{
				"position": map[string]float64{
					"x": 200.0,
					"y": 150.0,
				},
			},
		}

		answerBody, err := json.Marshal(answer)
		require.NoError(t, err)

		validateResp, err := http.Post(baseURL+"/api/validate", "application/json", bytes.NewBuffer(answerBody))
		require.NoError(t, err)
		defer validateResp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, validateResp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(validateResp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Contains(t, result["error"], "слишком быстрое прохождение")
	})

	t.Run("Too Slow Completion", func(t *testing.T) {
		createReq := map[string]interface{}{
			"challenge_type": "drag-drop",
			"complexity":     50,
		}

		jsonBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var challenge map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&challenge)
		require.NoError(t, err)

		challengeID := challenge["challenge_id"].(string)

		time.Sleep(3 * time.Second)

		answer := map[string]interface{}{
			"challenge_id": challengeID,
			"answer": map[string]interface{}{
				"position": map[string]float64{
					"x": 200.0,
					"y": 150.0,
				},
			},
		}

		answerBody, err := json.Marshal(answer)
		require.NoError(t, err)

		validateResp, err := http.Post(baseURL+"/api/validate", "application/json", bytes.NewBuffer(answerBody))
		require.NoError(t, err)
		defer validateResp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, validateResp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(validateResp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Contains(t, result["error"], "превышено максимальное время")
	})
}

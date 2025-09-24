package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCaptchaServiceIntegration(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"38000/tcp"},
		Env: map[string]string{
			"MIN_PORT":             "38000",
			"LOG_LEVEL":            "debug",
			"CHALLENGE_TYPE":       entity.ChallengeTypeDragDrop,
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

	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var health map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)

		assert.Equal(t, "healthy", health["status"])
		assert.Equal(t, float64(38000), health["port"])
	})

	t.Run("Create Challenge", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"challenge_type": "drag-drop",
			"complexity":     50,
		}

		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var challenge map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&challenge)
		require.NoError(t, err)

		assert.Contains(t, challenge, "challenge_id")
		assert.Contains(t, challenge, "html")
		assert.Equal(t, "drag-drop", challenge["type"])
		assert.Equal(t, float64(50), challenge["complexity"])
	})

	t.Run("Validate Answer", func(t *testing.T) {
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

		validateReq := map[string]interface{}{
			"challenge_id": challengeID,
			"answer": map[string]interface{}{
				"position": map[string]float64{
					"x": 200.0,
					"y": 150.0,
				},
			},
		}

		validateBody, err := json.Marshal(validateReq)
		require.NoError(t, err)

		validateResp, err := http.Post(baseURL+"/api/validate", "application/json", bytes.NewBuffer(validateBody))
		require.NoError(t, err)
		defer validateResp.Body.Close()

		assert.Equal(t, http.StatusOK, validateResp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(validateResp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Contains(t, result, "valid")
		assert.Contains(t, result, "confidence")
	})

	t.Run("Slider Puzzle Challenge", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"challenge_type": "slider-puzzle",
			"complexity":     30,
		}

		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/api/challenge", "application/json", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var challenge map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&challenge)
		require.NoError(t, err)

		assert.Contains(t, challenge, "challenge_id")
		assert.Contains(t, challenge, "html")
		assert.Equal(t, "slider-puzzle", challenge["type"])
		assert.Equal(t, float64(30), challenge["complexity"])
	})
}

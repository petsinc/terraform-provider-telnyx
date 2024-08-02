package telnyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type TelnyxClient struct {
	apiKey  string
	baseURL string
	logger  *zap.Logger
}

func NewClient() *TelnyxClient {
	apiKey := os.Getenv("TELNYX_API_KEY")
	if apiKey == "" {
		panic("TELNYX_API_KEY environment variable must be set")
	}

	logLevel := getLogLevelFromEnv()
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)
	logger, _ := config.Build()

	return &TelnyxClient{apiKey: apiKey, baseURL: "https://api.telnyx.com/v2", logger: logger}
}

func getLogLevelFromEnv() zapcore.Level {
	level := os.Getenv("TELNYX_REST_CLIENT_LOG_LEVEL")
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (client *TelnyxClient) doRequest(method, path string, body interface{}, v interface{}) error {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			client.logger.Error("Error encoding request body", zap.Error(err))
			return err
		}
		client.logger.Info("Request Body", zap.String("body", string(bodyBytes)))
	}

	return client.retryRequest(method, path, bodyBytes, v)
}

func (client *TelnyxClient) retryRequest(method, path string, bodyBytes []byte, v interface{}) error {
	retryAttempts := 5
	var lastErr error

	for attempt := 0; attempt < retryAttempts; attempt++ {
		if attempt > 0 {
			waitTime := time.Duration(attempt*attempt) * time.Second // exponential backoff
			client.logger.Info("Retrying request", zap.Int("attempt", attempt), zap.Duration("wait_time", waitTime))
			time.Sleep(waitTime)
		}

		req, err := http.NewRequest(method, client.baseURL+path, bytes.NewReader(bodyBytes))
		if err != nil {
			client.logger.Error("Error creating request", zap.Error(err))
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+client.apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			client.logger.Error("Error making request", zap.Error(err))
			return err
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			client.logger.Error("Error reading response body", zap.Error(err))
			return err
		}

		client.logger.Info("Response Body", zap.String("response", string(respBody)))

		if resp.StatusCode == 429 {
			client.logger.Warn("Received 429 Too Many Requests", zap.Int("status_code", resp.StatusCode), zap.String("response", string(respBody)))
			lastErr = fmt.Errorf("received 429 Too Many Requests: %s", string(respBody))
			continue
		}

		if resp.StatusCode >= 400 {
			client.logger.Error("Received error response from API", zap.Int("status_code", resp.StatusCode), zap.String("response", string(respBody)))
			return fmt.Errorf("received error response from API: %s", string(respBody))
		}

		if v != nil {
			if err := json.Unmarshal(respBody, v); err != nil {
				client.logger.Error("Error unmarshaling response", zap.Error(err))
				return err
			}
		}

		return nil
	}

	return lastErr
}

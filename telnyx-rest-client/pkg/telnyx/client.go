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
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			client.logger.Error("Error encoding request body", zap.Error(err))
			return err
		}
		client.logger.Info("Request Body", zap.String("body", buf.String()))
	}

	req, err := http.NewRequest(method, client.baseURL+path, &buf)
	if err != nil {
		client.logger.Error("Error creating request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	var resp *http.Response
	var respBody []byte

	for attempts := 0; attempts < 5; attempts++ { // Maximum of 5 attempts
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			client.logger.Error("Error making request", zap.Error(err))
			return err
		}
		defer resp.Body.Close()

		respBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			client.logger.Error("Error reading response body", zap.Error(err))
			return err
		}

		client.logger.Info("Response Body", zap.String("response", string(respBody)))

		if resp.StatusCode == http.StatusTooManyRequests {
			client.logger.Warn("Rate limited, retrying...", zap.Int("attempt", attempts+1))
			time.Sleep(time.Duration(attempts+1) * time.Second) // Exponential backoff
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

	client.logger.Error("Exceeded maximum retry attempts due to rate limiting")
	return fmt.Errorf("exceeded maximum retry attempts due to rate limiting: %s", string(respBody))
}

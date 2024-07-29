package telnyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
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

	logger, _ := zap.NewProduction()
	return &TelnyxClient{apiKey: apiKey, baseURL: "https://api.telnyx.com/v2", logger: logger}
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

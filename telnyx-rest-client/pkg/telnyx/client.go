package telnyx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type TelnyxClient struct {
	ApiToken   string
	ApiBaseURL string
	logger     *zap.Logger
}

func NewClient(apiToken string, logger *zap.Logger) *TelnyxClient {
	return &TelnyxClient{ApiToken: apiToken, ApiBaseURL: "https://api.telnyx.com/v2", logger: logger}
}

func (client *TelnyxClient) doRequest(method, endpoint string, body map[string]interface{}, result interface{}) error {
	url := client.ApiBaseURL + endpoint
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			client.logger.Error("Error marshalling JSON", zap.Error(err), zap.Any("body", body))
			return err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		client.logger.Error("Error creating HTTP request", zap.Error(err), zap.String("url", url))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.ApiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		client.logger.Error("Error making HTTP request", zap.Error(err), zap.String("url", url))
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.logger.Error("Error reading HTTP response body", zap.Error(err))
		return err
	}

	if resp.StatusCode >= 400 {
		client.logger.Error("Received error response from API", zap.Int("status_code", resp.StatusCode), zap.String("response", string(bodyBytes)))
		return err
	}

	if len(bodyBytes) > 0 && result != nil {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			client.logger.Error("Error unmarshalling JSON response", zap.Error(err))
			return err
		}
	}

	return nil
}

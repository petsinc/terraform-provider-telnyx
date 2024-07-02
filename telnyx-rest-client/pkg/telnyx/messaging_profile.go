package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateMessagingProfile(name string, whitelistedDestinations []string, webhookURL, webhookFailoverURL string, enabled bool) (*MessagingProfile, error) {
	body := map[string]interface{}{
		"name":                     name,
		"whitelisted_destinations": whitelistedDestinations,
		"enabled":                  enabled,
		"webhook_url":              webhookURL,
		"webhook_failover_url":     webhookFailoverURL,
		"webhook_api_version":      "2",
	}
	var result struct {
		Data MessagingProfile `json:"data"`
	}
	err := client.doRequest("POST", "/messaging_profiles", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateMessagingProfile(profileID, name string, whitelistedDestinations []string, webhookURL, webhookFailoverURL string, enabled bool) (*MessagingProfile, error) {
	body := map[string]interface{}{
		"name":                     name,
		"whitelisted_destinations": whitelistedDestinations,
		"enabled":                  enabled,
		"webhook_url":              webhookURL,
		"webhook_failover_url":     webhookFailoverURL,
		"webhook_api_version":      "2",
	}
	var result struct {
		Data MessagingProfile `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/messaging_profiles/%s", profileID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteMessagingProfile(profileID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/messaging_profiles/%s", profileID), nil, nil)
}

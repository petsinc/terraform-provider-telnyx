package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateMessagingProfile(profile MessagingProfile) (*MessagingProfile, error) {
	var result struct {
		Data MessagingProfile `json:"data"`
	}
	err := client.doRequest("POST", "/messaging_profiles", profile, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) GetMessagingProfile(profileID string) (*MessagingProfile, error) {
	var result struct {
		Data MessagingProfile `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/messaging_profiles/%s", profileID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateMessagingProfile(profileID string, profile MessagingProfile) (*MessagingProfile, error) {
	var result struct {
		Data MessagingProfile `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/messaging_profiles/%s", profileID), profile, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteMessagingProfile(profileID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/messaging_profiles/%s", profileID), nil, nil)
}

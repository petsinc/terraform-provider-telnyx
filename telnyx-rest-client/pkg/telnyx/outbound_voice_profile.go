package telnyx

import (
	"fmt"
	"go.uber.org/zap"
)

func (client *TelnyxClient) CreateOutboundVoiceProfile(profile OutboundVoiceProfile) (*OutboundVoiceProfile, error) {
	var result struct {
		Data OutboundVoiceProfile `json:"data"`
	}
	err := client.doRequest("POST", "/outbound_voice_profiles", profile, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) GetOutboundVoiceProfile(profileID string) (*OutboundVoiceProfile, error) {
	var result struct {
		Data OutboundVoiceProfile `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/outbound_voice_profiles/%s", profileID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateOutboundVoiceProfile(profileID string, profile OutboundVoiceProfile) (*OutboundVoiceProfile, error) {
	var result struct {
		Data OutboundVoiceProfile `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/outbound_voice_profiles/%s", profileID), profile, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteOutboundVoiceProfile(profileID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/outbound_voice_profiles/%s", profileID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting outbound voice profile", zap.Error(err), zap.String("profileID", profileID))
	}
	return err
}

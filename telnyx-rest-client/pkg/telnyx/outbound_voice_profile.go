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

func (client *TelnyxClient) UpdateOutboundVoiceProfile(outboundVoiceProfileID string, profile OutboundVoiceProfile) (*OutboundVoiceProfile, error) {
	var result struct {
		Data OutboundVoiceProfile `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/outbound_voice_profiles/%s", outboundVoiceProfileID), profile, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteOutboundVoiceProfile(outboundVoiceProfileID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/outbound_voice_profiles/%s", outboundVoiceProfileID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting outbound voice profile", zap.Error(err), zap.String("outboundVoiceProfileID", outboundVoiceProfileID))
	}
	return err
}

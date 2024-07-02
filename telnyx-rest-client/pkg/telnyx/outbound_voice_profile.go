package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateOutboundVoiceProfile(profile OutboundVoiceProfile) (*OutboundVoiceProfile, error) {
	body := map[string]interface{}{
		"name":                      profile.Name,
		"traffic_type":              profile.TrafficType,
		"service_plan":              profile.ServicePlan,
		"concurrent_call_limit":     profile.ConcurrentCallLimit,
		"enabled":                   profile.Enabled,
		"tags":                      profile.Tags,
		"usage_payment_method":      profile.UsagePaymentMethod,
		"whitelisted_destinations":  profile.WhitelistedDestinations,
		"max_destination_rate":      profile.MaxDestinationRate,
		"daily_spend_limit":         profile.DailySpendLimit,
		"daily_spend_limit_enabled": profile.DailySpendLimitEnabled,
		"billing_group_id":          profile.BillingGroupID,
	}

	if profile.CallRecording.Type != "" {
		body["call_recording"] = map[string]interface{}{
			"call_recording_type":                 profile.CallRecording.Type,
			"call_recording_caller_phone_numbers": profile.CallRecording.CallerPhoneNumbers,
			"call_recording_channels":             profile.CallRecording.Channels,
			"call_recording_format":               profile.CallRecording.Format,
		}
	}
	var result struct {
		Data OutboundVoiceProfile `json:"data"`
	}
	err := client.doRequest("POST", "/outbound_voice_profiles", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateOutboundVoiceProfile(outboundVoiceProfileID string, profile OutboundVoiceProfile) (*OutboundVoiceProfile, error) {
	body := map[string]interface{}{
		"name":                      profile.Name,
		"traffic_type":              profile.TrafficType,
		"service_plan":              profile.ServicePlan,
		"concurrent_call_limit":     profile.ConcurrentCallLimit,
		"enabled":                   profile.Enabled,
		"tags":                      profile.Tags,
		"usage_payment_method":      profile.UsagePaymentMethod,
		"whitelisted_destinations":  profile.WhitelistedDestinations,
		"max_destination_rate":      profile.MaxDestinationRate,
		"daily_spend_limit":         profile.DailySpendLimit,
		"daily_spend_limit_enabled": profile.DailySpendLimitEnabled,
		"billing_group_id":          profile.BillingGroupID,
	}

	if profile.CallRecording.Type != "" {
		body["call_recording"] = map[string]interface{}{
			"call_recording_type":                 profile.CallRecording.Type,
			"call_recording_caller_phone_numbers": profile.CallRecording.CallerPhoneNumbers,
			"call_recording_channels":             profile.CallRecording.Channels,
			"call_recording_format":               profile.CallRecording.Format,
		}
	}

	var result struct {
		Data OutboundVoiceProfile `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/outbound_voice_profiles/%s", outboundVoiceProfileID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteOutboundVoiceProfile(outboundVoiceProfileID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/outbound_voice_profiles/%s", outboundVoiceProfileID), nil, nil)
}

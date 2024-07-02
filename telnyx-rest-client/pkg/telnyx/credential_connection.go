package telnyx

import (
	"fmt"

	"go.uber.org/zap"
)

// Create Credential Connection
func (client *TelnyxClient) CreateCredentialConnection(profile FQDNConnection) (*FQDNConnection, error) {
	body := map[string]interface{}{
		"active":                                profile.Active,
		"anchorsite_override":                   profile.AnchorsiteOverride,
		"connection_name":                       profile.ConnectionName,
		"transport_protocol":                    profile.TransportProtocol,
		"default_on_hold_comfort_noise_enabled": profile.DefaultOnHoldComfortNoiseEnabled,
		"dtmf_type":                             profile.DTMFType,
		"encode_contact_header_enabled":         profile.EncodeContactHeaderEnabled,
		"encrypted_media":                       profile.EncryptedMedia,
		"onnet_t38_passthrough_enabled":         profile.OnnetT38PassthroughEnabled,
		"user_name":                             profile.Username,
		"password":                              profile.Password,
		"sip_uri_calling_preference":            profile.SipUriCallingPreference,
		"webhook_event_url":                     profile.WebhookEventURL,
		"webhook_event_failover_url":            profile.WebhookEventFailoverURL,
		"webhook_api_version":                   profile.WebhookAPIVersion,
		"webhook_timeout_secs":                  profile.WebhookTimeoutSecs,
		"rtcp_settings": map[string]interface{}{
			"port":                  profile.RTCPSettings.Port,
			"capture_enabled":       profile.RTCPSettings.CaptureEnabled,
			"report_frequency_secs": profile.RTCPSettings.ReportFrequencySecs,
		},
		"inbound": map[string]interface{}{
			"ani_number_format":              profile.Inbound.ANINumberFormat,
			"dnis_number_format":             profile.Inbound.DNISNumberFormat,
			"codecs":                         profile.Inbound.Codecs,
			"default_routing_method":         profile.Inbound.DefaultRoutingMethod,
			"channel_limit":                  profile.Inbound.ChannelLimit,
			"generate_ringback_tone":         profile.Inbound.GenerateRingbackTone,
			"isup_headers_enabled":           profile.Inbound.ISUPHeadersEnabled,
			"prack_enabled":                  profile.Inbound.PRACKEnabled,
			"privacy_zone_enabled":           profile.Inbound.PrivacyZoneEnabled,
			"sip_compact_headers_enabled":    profile.Inbound.SIPCompactHeadersEnabled,
			"sip_region":                     profile.Inbound.SIPRegion,
			"sip_subdomain":                  profile.Inbound.SIPSubdomain,
			"sip_subdomain_receive_settings": profile.Inbound.SIPSubdomainReceiveSettings,
			"timeout_1xx_secs":               profile.Inbound.Timeout1xxSecs,
			"timeout_2xx_secs":               profile.Inbound.Timeout2xxSecs,
			"shaken_stir_enabled":            profile.Inbound.ShakenSTIREnabled,
		},
		"outbound": map[string]interface{}{
			"ani_override":              profile.Outbound.ANIOverride,
			"ani_override_type":         profile.Outbound.ANIOverrideType,
			"call_parking_enabled":      profile.Outbound.CallParkingEnabled,
			"channel_limit":             profile.Outbound.ChannelLimit,
			"generate_ringback_tone":    profile.Outbound.GenerateRingbackTone,
			"instant_ringback_enabled":  profile.Outbound.InstantRingbackEnabled,
			"ip_authentication_method":  profile.Outbound.IPAuthenticationMethod,
			"ip_authentication_token":   profile.Outbound.IPAuthenticationToken,
			"localization":              profile.Outbound.Localization,
			"outbound_voice_profile_id": profile.Outbound.OutboundVoiceProfileID,
			"t38_reinvite_source":       profile.Outbound.T38ReinviteSource,
			"tech_prefix":               profile.Outbound.TechPrefix,
			"encrypted_media":           profile.Outbound.EncryptedMedia,
			"timeout_1xx_secs":          profile.Outbound.Timeout1xxSecs,
			"timeout_2xx_secs":          profile.Outbound.Timeout2xxSecs,
		},
	}

	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("POST", "/credential_connections", body, &result)
	if err != nil {
		client.logger.Error("Error creating credential connection", zap.Error(err))
		return nil, err
	}
	return &result.Data, nil
}

// Update Credential Connection
func (client *TelnyxClient) UpdateCredentialConnection(credentialConnectionID string, profile FQDNConnection) (*FQDNConnection, error) {
	body := map[string]interface{}{
		"user_name":                             profile.Username,
		"password":                              profile.Password,
		"active":                                profile.Active,
		"anchorsite_override":                   profile.AnchorsiteOverride,
		"connection_name":                       profile.ConnectionName,
		"transport_protocol":                    profile.TransportProtocol,
		"default_on_hold_comfort_noise_enabled": profile.DefaultOnHoldComfortNoiseEnabled,
		"dtmf_type":                             profile.DTMFType,
		"encode_contact_header_enabled":         profile.EncodeContactHeaderEnabled,
		"encrypted_media":                       profile.EncryptedMedia,
		"onnet_t38_passthrough_enabled":         profile.OnnetT38PassthroughEnabled,
		"ios_push_credential_id":                profile.IosPushCredentialID,
		"android_push_credential_id":            profile.AndroidPushCredentialID,
		"microsoft_teams_sbc":                   profile.MicrosoftTeamsSbc,
		"webhook_event_url":                     profile.WebhookEventURL,
		"webhook_event_failover_url":            profile.WebhookEventFailoverURL,
		"webhook_api_version":                   profile.WebhookAPIVersion,
		"webhook_timeout_secs":                  profile.WebhookTimeoutSecs,
		"fqdn_outbound_authentication":          profile.FQDNOutboundAuthentication,
		"rtcp_settings": map[string]interface{}{
			"port":                  profile.RTCPSettings.Port,
			"capture_enabled":       profile.RTCPSettings.CaptureEnabled,
			"report_frequency_secs": profile.RTCPSettings.ReportFrequencySecs,
		},
		"inbound": map[string]interface{}{
			"ani_number_format":              profile.Inbound.ANINumberFormat,
			"dnis_number_format":             profile.Inbound.DNISNumberFormat,
			"codecs":                         profile.Inbound.Codecs,
			"default_routing_method":         profile.Inbound.DefaultRoutingMethod,
			"channel_limit":                  profile.Inbound.ChannelLimit,
			"generate_ringback_tone":         profile.Inbound.GenerateRingbackTone,
			"isup_headers_enabled":           profile.Inbound.ISUPHeadersEnabled,
			"prack_enabled":                  profile.Inbound.PRACKEnabled,
			"privacy_zone_enabled":           profile.Inbound.PrivacyZoneEnabled,
			"sip_compact_headers_enabled":    profile.Inbound.SIPCompactHeadersEnabled,
			"sip_region":                     profile.Inbound.SIPRegion,
			"sip_subdomain":                  profile.Inbound.SIPSubdomain,
			"sip_subdomain_receive_settings": profile.Inbound.SIPSubdomainReceiveSettings,
			"timeout_1xx_secs":               profile.Inbound.Timeout1xxSecs,
			"timeout_2xx_secs":               profile.Inbound.Timeout2xxSecs,
			"shaken_stir_enabled":            profile.Inbound.ShakenSTIREnabled,
		},
		"outbound": map[string]interface{}{
			"ani_override":              profile.Outbound.ANIOverride,
			"ani_override_type":         profile.Outbound.ANIOverrideType,
			"call_parking_enabled":      profile.Outbound.CallParkingEnabled,
			"channel_limit":             profile.Outbound.ChannelLimit,
			"generate_ringback_tone":    profile.Outbound.GenerateRingbackTone,
			"instant_ringback_enabled":  profile.Outbound.InstantRingbackEnabled,
			"ip_authentication_method":  profile.Outbound.IPAuthenticationMethod,
			"ip_authentication_token":   profile.Outbound.IPAuthenticationToken,
			"localization":              profile.Outbound.Localization,
			"outbound_voice_profile_id": profile.Outbound.OutboundVoiceProfileID,
			"t38_reinvite_source":       profile.Outbound.T38ReinviteSource,
			"tech_prefix":               profile.Outbound.TechPrefix,
			"encrypted_media":           profile.Outbound.EncryptedMedia,
			"timeout_1xx_secs":          profile.Outbound.Timeout1xxSecs,
			"timeout_2xx_secs":          profile.Outbound.Timeout2xxSecs,
		},
	}

	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), body, &result)
	if err != nil {
		client.logger.Error("Error updating credential connection", zap.Error(err), zap.String("credentialConnectionID", credentialConnectionID))
		return nil, err
	}
	return &result.Data, nil
}

// Delete Credential Connection
func (client *TelnyxClient) DeleteCredentialConnection(credentialConnectionID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting credential connection", zap.Error(err), zap.String("credentialConnectionID", credentialConnectionID))
	}
	return err
}

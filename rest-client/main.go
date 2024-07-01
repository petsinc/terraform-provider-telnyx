package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type TelnyxClient struct {
	ApiToken   string
	ApiBaseURL string
}

func NewClient(apiToken string) *TelnyxClient {
	return &TelnyxClient{ApiToken: apiToken, ApiBaseURL: "https://api.telnyx.com/v2"}
}

// Struct Definitions

type BillingGroup struct {
	RecordType     string    `json:"record_type"`
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at,omitempty"`
}

type MessagingProfile struct {
	ID                      string                `json:"id"`
	Name                    string                `json:"name"`
	Enabled                 bool                  `json:"enabled"`
	WebhookURL              string                `json:"webhook_url"`
	WebhookFailoverURL      string                `json:"webhook_failover_url"`
	WebhookAPIVersion       string                `json:"webhook_api_version"`
	WhitelistedDestinations []string              `json:"whitelisted_destinations"`
	NumberPoolSettings      *NumberPoolSettings   `json:"number_pool_settings,omitempty"`
	URLShortenerSettings    *URLShortenerSettings `json:"url_shortener_settings,omitempty"`
	AlphaSender             *string               `json:"alpha_sender,omitempty"`
	CreatedAt               time.Time             `json:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at"`
	V1Secret                string                `json:"v1_secret"`
}

type NumberPoolSettings struct {
	TollFreeWeight float64 `json:"toll_free_weight"`
	LongCodeWeight float64 `json:"long_code_weight"`
	SkipUnhealthy  bool    `json:"skip_unhealthy"`
	StickySender   bool    `json:"sticky_sender"`
	Geomatch       bool    `json:"geomatch"`
}

type URLShortenerSettings struct {
	Domain               string `json:"domain"`
	Prefix               string `json:"prefix"`
	ReplaceBlacklistOnly bool   `json:"replace_blacklist_only"`
	SendWebhooks         bool   `json:"send_webhooks"`
}

type OutboundVoiceProfile struct {
	ID                      string        `json:"id,omitempty"`
	Name                    string        `json:"name"`
	ConnectionsCount        int           `json:"connections_count,omitempty"`
	TrafficType             string        `json:"traffic_type"`
	ServicePlan             string        `json:"service_plan"`
	ConcurrentCallLimit     int           `json:"concurrent_call_limit"`
	Enabled                 bool          `json:"enabled"`
	Tags                    []string      `json:"tags"`
	UsagePaymentMethod      string        `json:"usage_payment_method"`
	WhitelistedDestinations []string      `json:"whitelisted_destinations"`
	MaxDestinationRate      float64       `json:"max_destination_rate"`
	DailySpendLimit         string        `json:"daily_spend_limit"`
	DailySpendLimitEnabled  bool          `json:"daily_spend_limit_enabled"`
	CallRecording           CallRecording `json:"call_recording"`
	BillingGroupID          string        `json:"billing_group_id"`
	CreatedAt               time.Time     `json:"created_at,omitempty"`
	UpdatedAt               time.Time     `json:"updated_at,omitempty"`
}

type CallRecording struct {
	Type               string   `json:"call_recording_type"`
	CallerPhoneNumbers []string `json:"call_recording_caller_phone_numbers"`
	Channels           string   `json:"call_recording_channels"`
	Format             string   `json:"call_recording_format"`
}

type PhoneNumberReservation struct {
	ID           string `json:"id"`
	RecordType   string `json:"record_type"`
	PhoneNumbers []struct {
		ID          string    `json:"id"`
		RecordType  string    `json:"record_type"`
		PhoneNumber string    `json:"phone_number"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		ExpiredAt   time.Time `json:"expired_at"`
	} `json:"phone_numbers"`
	Status            string    `json:"status"`
	CustomerReference string    `json:"customer_reference"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PhoneNumberOrder struct {
	ID           string `json:"id"`
	RecordType   string `json:"record_type"`
	PhoneNumbers []struct {
		ID                     string `json:"id"`
		RecordType             string `json:"record_type"`
		PhoneNumber            string `json:"phone_number"`
		BundleID               string `json:"bundle_id"`
		PhoneNumberType        string `json:"phone_number_type"`
		RegulatoryRequirements []struct {
			RecordType    string `json:"record_type"`
			RequirementID string `json:"requirement_id"`
			FieldValue    string `json:"field_value"`
			FieldType     string `json:"field_type"`
		} `json:"regulatory_requirements"`
		CountryCode     string `json:"country_code"`
		RequirementsMet bool   `json:"requirements_met"`
		Status          string `json:"status"`
	} `json:"phone_numbers"`
	PhoneNumbersCount  int       `json:"phone_numbers_count"`
	ConnectionID       string    `json:"connection_id"`
	MessagingProfileID string    `json:"messaging_profile_id"`
	BillingGroupID     string    `json:"billing_group_id"`
	Status             string    `json:"status"`
	CustomerReference  string    `json:"customer_reference"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RTCPSettings struct {
	Port                string `json:"port"`
	CaptureEnabled      bool   `json:"capture_enabled"`
	ReportFrequencySecs int    `json:"report_frequency_secs"`
}

type InboundSettings struct {
	ANINumberFormat             string   `json:"ani_number_format"`
	DNISNumberFormat            string   `json:"dnis_number_format"`
	Codecs                      []string `json:"codecs"`
	DefaultRoutingMethod        string   `json:"default_routing_method"`
	ChannelLimit                int      `json:"channel_limit"`
	GenerateRingbackTone        bool     `json:"generate_ringback_tone"`
	ISUPHeadersEnabled          bool     `json:"isup_headers_enabled"`
	PRACKEnabled                bool     `json:"prack_enabled"`
	PrivacyZoneEnabled          bool     `json:"privacy_zone_enabled"`
	SIPCompactHeadersEnabled    bool     `json:"sip_compact_headers_enabled"`
	SIPRegion                   string   `json:"sip_region"`
	SIPSubdomain                string   `json:"sip_subdomain"`
	SIPSubdomainReceiveSettings string   `json:"sip_subdomain_receive_settings"`
	Timeout1xxSecs              int      `json:"timeout_1xx_secs"`
	Timeout2xxSecs              int      `json:"timeout_2xx_secs"`
	ShakenSTIREnabled           bool     `json:"shaken_stir_enabled"`
}

type OutboundSettings struct {
	ANIOverride            string `json:"ani_override"`
	ANIOverrideType        string `json:"ani_override_type"`
	CallParkingEnabled     bool   `json:"call_parking_enabled"`
	ChannelLimit           int    `json:"channel_limit"`
	GenerateRingbackTone   bool   `json:"generate_ringback_tone"`
	InstantRingbackEnabled bool   `json:"instant_ringback_enabled"`
	IPAuthenticationMethod string `json:"ip_authentication_method"`
	IPAuthenticationToken  string `json:"ip_authentication_token"`
	Localization           string `json:"localization"`
	OutboundVoiceProfileID string `json:"outbound_voice_profile_id"`
	T38ReinviteSource      string `json:"t38_reinvite_source"`
	TechPrefix             string `json:"tech_prefix"`
	EncryptedMedia         string `json:"encrypted_media"`
	Timeout1xxSecs         int    `json:"timeout_1xx_secs"`
	Timeout2xxSecs         int    `json:"timeout_2xx_secs"`
}

// New FQDN Struct
type FQDN struct {
	ID            string    `json:"id"`
	RecordType    string    `json:"record_type"`
	ConnectionID  string    `json:"connection_id"`
	FQDN          string    `json:"fqdn"`
	Port          int       `json:"port"`
	DNSRecordType string    `json:"dns_record_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FQDNConnection Struct with Credential Fields
type FQDNConnection struct {
	ID                               string           `json:"id"`
	RecordType                       string           `json:"record_type"`
	Active                           bool             `json:"active"`
	AnchorsiteOverride               string           `json:"anchorsite_override"`
	ConnectionName                   string           `json:"connection_name"`
	TransportProtocol                string           `json:"transport_protocol"`
	DefaultOnHoldComfortNoiseEnabled bool             `json:"default_on_hold_comfort_noise_enabled"`
	DTMFType                         string           `json:"dtmf_type"`
	EncodeContactHeaderEnabled       bool             `json:"encode_contact_header_enabled"`
	EncryptedMedia                   *string          `json:"encrypted_media,omitempty"`
	OnnetT38PassthroughEnabled       bool             `json:"onnet_t38_passthrough_enabled"`
	IosPushCredentialID              *string          `json:"ios_push_credential_id,omitempty"`
	AndroidPushCredentialID          *string          `json:"android_push_credential_id,omitempty"`
	MicrosoftTeamsSbc				 bool			   `json:"microsoft_teams_sbc`
	WebhookEventURL                  string           `json:"webhook_event_url"`
	WebhookEventFailoverURL          string          `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion                string           `json:"webhook_api_version"`
	WebhookTimeoutSecs               int             `json:"webhook_timeout_secs,omitempty"`
	RTCPSettings                     RTCPSettings     `json:"rtcp_settings"`
	Inbound                          InboundSettings  `json:"inbound"`
	Outbound                         OutboundSettings `json:"outbound"`
	CreatedAt                        time.Time        `json:"created_at"`
	UpdatedAt                        time.Time        `json:"updated_at"`
	Username                         *string          `json:"user_name,omitempty"`
	Password                         *string          `json:"password,omitempty"`
	SipUriCallingPreference          *string          `json:"sip_uri_calling_preference,omitempty"`
}

type PhoneNumber struct {
	ID                    string    `json:"id"`
	RecordType            string    `json:"record_type"`
	PhoneNumber           string    `json:"phone_number"`
	Status                string    `json:"status"`
	Tags                  []string  `json:"tags"`
	ExternalPin           string    `json:"external_pin"`
	ConnectionID          string    `json:"connection_id"`
	ConnectionName        string    `json:"connection_name"`
	CustomerReference     string    `json:"customer_reference"`
	MessagingProfileID    string    `json:"messaging_profile_id"`
	MessagingProfileName  string    `json:"messaging_profile_name"`
	BillingGroupID        string    `json:"billing_group_id"`
	EmergencyEnabled      bool      `json:"emergency_enabled"`
	EmergencyAddressID    string    `json:"emergency_address_id"`
	CallForwardingEnabled bool      `json:"call_forwarding_enabled"`
	CNAMListingEnabled    bool      `json:"cnam_listing_enabled"`
	CallerIDNameEnabled   bool      `json:"caller_id_name_enabled"`
	CallRecordingEnabled  bool      `json:"call_recording_enabled"`
	T38FaxGatewayEnabled  bool      `json:"t38_fax_gateway_enabled"`
	NumberLevelRouting    string    `json:"number_level_routing"`
	PhoneNumberType       string    `json:"phone_number_type"`
	PurchasedAt           time.Time `json:"purchased_at"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	HDVoiceEnabled        bool      `json:"hd_voice_enabled"`
}

// Helper Method

func (client *TelnyxClient) doRequest(method, endpoint string, body map[string]interface{}, result interface{}) error {
	url := client.ApiBaseURL + endpoint
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.ApiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if len(bodyBytes) > 0 && result != nil {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("error unmarshalling JSON response: %w", err)
		}
	}

	return nil
}

// TelnyxClient Methods

// Billing Group Operations

func (client *TelnyxClient) CreateBillingGroup(name string) (*BillingGroup, error) {
	body := map[string]interface{}{
		"name": name,
	}
	var result struct {
		Data BillingGroup `json:"data"`
	}
	err := client.doRequest("POST", "/billing_groups", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateBillingGroup(billingGroupID, name string) (*BillingGroup, error) {
	body := map[string]interface{}{
		"name": name,
	}
	var result struct {
		Data BillingGroup `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/billing_groups/%s", billingGroupID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteBillingGroup(billingGroupID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/billing_groups/%s", billingGroupID), nil, nil)
}

// Outbound Voice Profile Operations

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

// Messaging Profile Operations

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

// FQDN Connection Operations

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
		return nil, err
	}
	return &result.Data, nil
}

// Update Credential Connection
func (client *TelnyxClient) UpdateCredentialConnection(credentialConnectionID string, profile FQDNConnection) (*FQDNConnection, error) {
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
		"username_name":                         profile.Username,
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
	err := client.doRequest("PATCH", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Delete Credential Connection
func (client *TelnyxClient) DeleteCredentialConnection(credentialConnectionID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), nil, nil)
}

func (client *TelnyxClient) CreateFQDNConnection(profile FQDNConnection) (*FQDNConnection, error) {
	body := map[string]interface{}{
		"user_name":							 "hptest12345",
		"password":								"hptest54321",
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
		"webhook_event_url":          profile.WebhookEventURL,
		"webhook_event_failover_url": profile.WebhookEventFailoverURL,
		"webhook_api_version":        profile.WebhookAPIVersion,
		"webhook_timeout_secs":       profile.WebhookTimeoutSecs,
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

	
	// Pretty print the body as JSON
	prettyJSON, lol := json.MarshalIndent(body, "", "  ")
	if lol != nil {
		fmt.Println("Failed to generate JSON:", lol)
		return nil, lol
	}
	fmt.Println("Request Body JSON:")
	fmt.Println(string(prettyJSON))

	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("POST", "/fqdn_connections", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) CreateFQDN(connectionID, fqdn, dnsRecordType string, port int) (*FQDN, error) {
	body := map[string]interface{}{
		"connection_id":   connectionID,
		"fqdn":            fqdn,
		"dns_record_type": dnsRecordType,
		"port":            port,
	}
	var result struct {
		Data FQDN `json:"data"`
	}
	err := client.doRequest("POST", "/fqdns", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateFQDN(fqdnID string, fqdn, dnsRecordType string, port int) (*FQDN, error) {
	body := map[string]interface{}{
		"fqdn":            fqdn,
		"dns_record_type": dnsRecordType,
		"port":            port,
	}
	var result struct {
		Data FQDN `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/fqdns/%d", fqdnID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteFQDN(fqdnID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/fqdns/%d", fqdnID), nil, nil)
}

func (client *TelnyxClient) UpdateFQDNConnection(fqdnConnectionID string, profile FQDNConnection) (*FQDNConnection, error) {
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
		"ios_push_credential_id":                profile.IosPushCredentialID,
		"android_push_credential_id":            profile.AndroidPushCredentialID,
		"microsoft_teams_sbc":                   profile.MicrosoftTeamsSbc,
		"webhook_event_url":          profile.WebhookEventURL,
		"webhook_event_failover_url": profile.WebhookEventFailoverURL,
		"webhook_api_version":        profile.WebhookAPIVersion,
		"webhook_timeout_secs":       profile.WebhookTimeoutSecs,
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
	err := client.doRequest("PATCH", fmt.Sprintf("/fqdn_connections/%s", fqdnConnectionID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteFQDNConnection(connectionID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/fqdn_connections/%s", connectionID), nil, nil)
}

// Phone Number Operations

func (client *TelnyxClient) UpdatePhoneNumber(phoneNumberID, customerReference, connectionID, billingGroupID string, tags []string, hdVoiceEnabled bool) (*PhoneNumber, error) {
	body := map[string]interface{}{
		"customer_reference": customerReference,
		"connection_id":      connectionID,
		"billing_group_id":   billingGroupID,
		"tags":               tags,
		"hd_voice_enabled":   hdVoiceEnabled,
	}
	var result struct {
		Data PhoneNumber `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/phone_numbers/%s", phoneNumberID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeletePhoneNumber(phoneNumberID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/phone_numbers/%s", phoneNumberID), nil, nil)
}

// Number Order Operations

func (client *TelnyxClient) CreateNumberOrder(phoneNumbers []string, connectionID, messagingProfileID, billingGroupID, customerReference string) (*PhoneNumberOrder, error) {
	var phoneNumbersMap []map[string]string
	for _, phoneNumber := range phoneNumbers {
		phoneNumbersMap = append(phoneNumbersMap, map[string]string{"phone_number": phoneNumber})
	}
	body := map[string]interface{}{
		"phone_numbers":        phoneNumbersMap,
		"connection_id":        connectionID,
		"messaging_profile_id": messagingProfileID,
		"billing_group_id":     billingGroupID,
		"customer_reference":   customerReference,
	}
	var result struct {
		Data PhoneNumberOrder `json:"data"`
	}
	err := client.doRequest("POST", "/number_orders", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateNumberOrder(numberOrderID, customerReference string, regulatoryRequirements []map[string]string) (*PhoneNumberOrder, error) {
	body := map[string]interface{}{
		"customer_reference":      customerReference,
		"regulatory_requirements": regulatoryRequirements,
	}
	var result struct {
		Data PhoneNumberOrder `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/number_orders/%s", numberOrderID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Number Reservation Operations

func (client *TelnyxClient) CreateNumberReservation(phoneNumbers []string, customerReference string) (*PhoneNumberReservation, error) {
	var phoneNumbersMap []map[string]string
	for _, phoneNumber := range phoneNumbers {
		phoneNumbersMap = append(phoneNumbersMap, map[string]string{"phone_number": phoneNumber})
	}
	body := map[string]interface{}{
		"phone_numbers":      phoneNumbersMap,
		"customer_reference": customerReference,
	}
	var result struct {
		Data PhoneNumberReservation `json:"data"`
	}
	err := client.doRequest("POST", "/number_reservations", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) ExtendPhoneNumberReservation(reservationID string) (*PhoneNumberReservation, error) {
	var result struct {
		Data PhoneNumberReservation `json:"data"`
	}
	err := client.doRequest("POST", fmt.Sprintf("/number_reservations/%s/actions/extend", reservationID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// TestRunner Struct

type TestRunner struct {
	client                 *TelnyxClient
	billingGroupID         string
	outboundVoiceProfileID string
	messagingProfileID     string
	fqdnConnectionID       string
	fqdnID                 string
	phoneNumberID          string
	numberReservationID    string
	numberOrderID          string
}

func NewTestRunner(client *TelnyxClient) *TestRunner {
	return &TestRunner{client: client}
}

func (runner *TestRunner) PerformCreates() {
	// // Create a Billing Group
	// billingGroup, err := runner.client.CreateBillingGroup("Test Billing Group")
	// if err != nil {
	// 	fmt.Printf("Error creating billing group: %v\n", err)
	// 	os.Exit(1)
	// }
	// runner.billingGroupID = billingGroup.ID
	// fmt.Printf("Created Billing Group:\nID: %s\nName: %s\nCreated At: %s\n",
	// 	billingGroup.ID, billingGroup.Name, billingGroup.CreatedAt)

	// // Create an Outbound Voice Profile
	// outboundVoiceProfile, err := runner.client.CreateOutboundVoiceProfile(OutboundVoiceProfile{
	// 	Name:                    "Test Outbound Profile",
	// 	TrafficType:             "conversational",
	// 	ServicePlan:             "global",
	// 	ConcurrentCallLimit:     10,
	// 	Enabled:                 true,
	// 	Tags:                    []string{"test-profile"},
	// 	UsagePaymentMethod:      "rate-deck",
	// 	WhitelistedDestinations: []string{"US"},
	// 	MaxDestinationRate:      10.0,
	// 	DailySpendLimit:         "100.00",
	// 	DailySpendLimitEnabled:  true,
	// 	BillingGroupID:          runner.billingGroupID,
	// 	CallRecording: CallRecording{
	// 		Type:               "all",
	// 		CallerPhoneNumbers: []string{},
	// 		Channels:           "single",
	// 		Format:             "wav",
	// 	},
	// })

	// if err != nil {
	// 	fmt.Printf("Error creating outbound voice profile: %v\n", err)
	// 	os.Exit(1)
	// }
	runner.outboundVoiceProfileID = "906df5d3-d34e-4858-8cb2-230c99488e95"
	// fmt.Printf("Created Outbound Voice Profile:\nID: %s\nName: %s\nCreated At: %s\n",
	// 	outboundVoiceProfile.ID, outboundVoiceProfile.Name, outboundVoiceProfile.CreatedAt)

	// // Create a Messaging Profile
	// messagingProfile, err := runner.client.CreateMessagingProfile("Test Profile", []string{"US"}, "https://www.example.com/hooks", "https://backup.example.com/hooks", true)
	// if err != nil {
	// 	fmt.Printf("Error creating messaging profile: %v\n", err)
	// 	os.Exit(1)
	// }
	runner.messagingProfileID = "4001906e-f677-4a2e-a8e0-2a29581e4b68"
	// fmt.Printf("Created Messaging Profile:\nID: %s\nName: %s\nEnabled: %t\nWebhook URL: %s\nWebhook Failover URL: %s\nWebhook API Version: %s\nWhitelisted Destinations: %v\nCreated At: %s\nUpdated At: %s\n",
	// 	messagingProfile.ID, messagingProfile.Name, messagingProfile.Enabled, messagingProfile.WebhookURL, messagingProfile.WebhookFailoverURL, messagingProfile.WebhookAPIVersion, messagingProfile.WhitelistedDestinations, messagingProfile.CreatedAt, messagingProfile.UpdatedAt)

	// Create an FQDN Connection
	fqdnConnection, err := runner.client.CreateFQDNConnection(FQDNConnection{
		Active:                           true,
		AnchorsiteOverride:               "Latency",
		ConnectionName:                   "Test FQDN Connection",
		TransportProtocol:                "UDP",
		Username:						stringPtr("hp-test-12345"),
		Password:						stringPtr("hp-test-54321"),
		DefaultOnHoldComfortNoiseEnabled: true,
		DTMFType:                         "RFC 2833",
		EncodeContactHeaderEnabled:       false,
		EncryptedMedia:                   nil,
		OnnetT38PassthroughEnabled:       false,
		MicrosoftTeamsSbc:				  false,
		WebhookEventURL:                  "https://www.example.com/hooks",
		WebhookEventFailoverURL:          "https://failover.example.com/hooks",
		WebhookAPIVersion:                "1",
		WebhookTimeoutSecs:               25,
		RTCPSettings: RTCPSettings{
			Port:                "rtp+1",
			CaptureEnabled:      false,
			ReportFrequencySecs: 5,
		},
		Inbound: InboundSettings{
			ANINumberFormat:             "E.164-national",
			DNISNumberFormat:            "e164",
			Codecs:                      []string{"G722", "G711U", "G711A", "G729", "OPUS", "H.264"},
			DefaultRoutingMethod:        "sequential",
			ChannelLimit:                10,
			GenerateRingbackTone:        true,
			ISUPHeadersEnabled:          true,
			PRACKEnabled:                true,
			PrivacyZoneEnabled:          true,
			SIPCompactHeadersEnabled:    true,
			SIPRegion:                   "US",
			SIPSubdomain:                "uniqueexample.sip.telnyx.com",
			SIPSubdomainReceiveSettings: "only_my_connections",
			Timeout1xxSecs:              3,
			Timeout2xxSecs:              90,
			ShakenSTIREnabled:           true,
		},
		Outbound: OutboundSettings{
			ANIOverride:            "+12345678901",
			ANIOverrideType:        "always",
			CallParkingEnabled:     true,
			ChannelLimit:           10,
			GenerateRingbackTone:   true,
			InstantRingbackEnabled: false, // Ensure only one ringback setting is enabled
			IPAuthenticationMethod: "credentials-connection",
			// IPAuthenticationToken:    "aBcD1234aBcD1234", // Ensure token is at least 12 characters
			Localization: "US",
			// OutboundVoiceProfileID:   runner.outboundVoiceProfileID,
			T38ReinviteSource: "customer",
			EncryptedMedia:    "SRTP",
			Timeout1xxSecs:    3,
			Timeout2xxSecs:    90,
		},
	})
	if err != nil {
		fmt.Printf("Error creating FQDN connection: %v\n", err)
		os.Exit(1)
	}
	runner.fqdnConnectionID = fqdnConnection.ID
	fmt.Printf("Created FQDN Connection:\nID: %s\nName: %s\nCreated At: %s\n",
		fqdnConnection.ID, fqdnConnection.ConnectionName, fqdnConnection.CreatedAt)

	// Create an FQDN and bind it to the connection
	fqdn, err := runner.client.CreateFQDN(runner.fqdnConnectionID, "test.sip.livekit.cloud", "a", 5060)
	if err != nil {
		fmt.Printf("Error creating FQDN: %v\n", err)
		os.Exit(1)
	}
	runner.fqdnID = fqdn.ID
	fmt.Printf("Created FQDN:\nID: %d\nFQDN: %s\nCreated At: %s\n",
		fqdn.ID, fqdn.FQDN, fqdn.CreatedAt)
}

func (runner *TestRunner) PerformUpdates() {
	// Update the Messaging Profile
	updatedMessagingProfile, err := runner.client.UpdateMessagingProfile(runner.messagingProfileID, "Updated Profile for Messages", []string{"US"}, "https://www.example.com/hooks", "https://backup.example.com/hooks", true)
	if err != nil {
		fmt.Printf("Error updating messaging profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated Messaging Profile:\nID: %s\nName: %s\nEnabled: %t\nWebhook URL: %s\nWebhook Failover URL: %s\nWebhook API Version: %s\nWhitelisted Destinations: %v\nCreated At: %s\nUpdated At: %s\n",
		updatedMessagingProfile.ID, updatedMessagingProfile.Name, updatedMessagingProfile.Enabled, updatedMessagingProfile.WebhookURL, updatedMessagingProfile.WebhookFailoverURL, updatedMessagingProfile.WebhookAPIVersion, updatedMessagingProfile.WhitelistedDestinations, updatedMessagingProfile.CreatedAt, updatedMessagingProfile.UpdatedAt)

	// Update the Outbound Voice Profile
	updatedOutboundVoiceProfile, err := runner.client.UpdateOutboundVoiceProfile(runner.outboundVoiceProfileID, OutboundVoiceProfile{
		Name:                    "Test Outbound Profile Updated",
		TrafficType:             "conversational",
		ServicePlan:             "global",
		ConcurrentCallLimit:     10,
		Enabled:                 true,
		Tags:                    []string{"test-profile"},
		UsagePaymentMethod:      "rate-deck",
		WhitelistedDestinations: []string{"US"},
		MaxDestinationRate:      10.0,
		DailySpendLimit:         "100.00",
		DailySpendLimitEnabled:  true,
		BillingGroupID:          runner.billingGroupID,
		CallRecording: CallRecording{
			Type:               "all",
			CallerPhoneNumbers: []string{},
			Channels:           "single",
			Format:             "wav",
		},
	})

	if err != nil {
		fmt.Printf("Error updating outbound voice profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated Outbound Voice Profile:\nID: %s\nName: %s\nUpdated At: %s\n",
		updatedOutboundVoiceProfile.ID, updatedOutboundVoiceProfile.Name, updatedOutboundVoiceProfile.UpdatedAt)

	// Update the FQDN Connection
	updatedFQDNConnection, err := runner.client.UpdateFQDNConnection(runner.fqdnConnectionID, FQDNConnection{
		Active:                           true,
		AnchorsiteOverride:               "Latency",
		ConnectionName:                   "Updated Test FQDN Connection",
		TransportProtocol:                "UDP",
		DefaultOnHoldComfortNoiseEnabled: true,
		DTMFType:                         "RFC 2833",
		EncodeContactHeaderEnabled:       false,
		EncryptedMedia:                   nil,
		OnnetT38PassthroughEnabled:       false,
		MicrosoftTeamsSbc:				  false,
		WebhookEventURL:                  "https://www.example.com/hooks",
		WebhookEventFailoverURL:          "https://failover.example.com/hooks",
		WebhookAPIVersion:                "1",
		WebhookTimeoutSecs:               25,
		RTCPSettings: RTCPSettings{
			Port:                "rtp+1",
			CaptureEnabled:      false,
			ReportFrequencySecs: 5,
		},
		Inbound: InboundSettings{
			ANINumberFormat:             "E.164-national",
			DNISNumberFormat:            "e164",
			Codecs:                      []string{"G722", "G711U", "G711A", "G729", "OPUS", "H.264"},
			DefaultRoutingMethod:        "sequential",
			ChannelLimit:                10,
			GenerateRingbackTone:        true,
			ISUPHeadersEnabled:          true,
			PRACKEnabled:                true,
			PrivacyZoneEnabled:          true,
			SIPCompactHeadersEnabled:    true,
			SIPRegion:                   "US",
			SIPSubdomain:                "uniqueexample.sip.telnyx.com",
			SIPSubdomainReceiveSettings: "only_my_connections",
			Timeout1xxSecs:              3,
			Timeout2xxSecs:              90,
			ShakenSTIREnabled:           true,
		},
		Outbound: OutboundSettings{
			ANIOverride:            "+12345678901",
			ANIOverrideType:        "always",
			CallParkingEnabled:     true,
			ChannelLimit:           10,
			GenerateRingbackTone:   true,
			InstantRingbackEnabled: false, // Ensure only one ringback setting is enabled
			IPAuthenticationMethod: "credentials-connection",
			// IPAuthenticationToken:    "aBcD1234aBcD1234", // Ensure token is at least 12 characters
			Localization:           "US",
			OutboundVoiceProfileID: runner.outboundVoiceProfileID,
			T38ReinviteSource:      "customer",
			EncryptedMedia:         "SRTP",
			Timeout1xxSecs:         3,
			Timeout2xxSecs:         90,
		},
	})
	if err != nil {
		fmt.Printf("Error updating FQDN connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated FQDN Connection:\nID: %s\nName: %s\nUpdated At: %s\n",
		updatedFQDNConnection.ID, updatedFQDNConnection.ConnectionName, updatedFQDNConnection.UpdatedAt)

	// Update the FQDN
	updatedFQDN, err := runner.client.UpdateFQDN(runner.fqdnID, "updated.test.sip.livekit.cloud", "a", 5060)
	if err != nil {
		fmt.Printf("Error updating FQDN: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated FQDN:\nID: %d\nFQDN: %s\nUpdated At: %s\n",
		updatedFQDN.ID, updatedFQDN.FQDN, updatedFQDN.UpdatedAt)
}

func (runner *TestRunner) PerformCascadingDeletes() {
	// Delete the FQDN
	err := runner.client.DeleteFQDN(runner.fqdnID)
	if err != nil {
		fmt.Printf("Error deleting FQDN: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Deleted FQDN")

	// Delete the FQDN Connection
	err = runner.client.DeleteFQDNConnection(runner.fqdnConnectionID)
	if err != nil {
		fmt.Printf("Error deleting FQDN connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Deleted FQDN Connection")

	// Delete the Messaging Profile
	err = runner.client.DeleteMessagingProfile(runner.messagingProfileID)
	if err != nil {
		fmt.Printf("Error deleting messaging profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Deleted Messaging Profile")

	// Delete the Outbound Voice Profile
	err = runner.client.DeleteOutboundVoiceProfile(runner.outboundVoiceProfileID)
	if err != nil {
		fmt.Printf("Error deleting outbound voice profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Deleted Outbound Voice Profile")

	// Delete the Billing Group
	err = runner.client.DeleteBillingGroup(runner.billingGroupID)
	if err != nil {
		fmt.Printf("Error deleting billing group: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Deleted Billing Group")
}

// Example Usage

func main() {
	client := NewClient("lol")
	runner := NewTestRunner(client)

	// Perform create operations
	runner.PerformCreates()

	// Perform update operations
	runner.PerformUpdates()

	// Perform cascading delete operations
	runner.PerformCascadingDeletes()
}

// cuz go is cool

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

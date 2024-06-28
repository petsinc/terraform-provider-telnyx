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
	RecordType     string     `json:"record_type"`
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Name           string     `json:"name"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type MessagingProfile struct {
	ID                      string              `json:"id"`
	Name                    string              `json:"name"`
	Enabled                 bool                `json:"enabled"`
	WebhookURL              string              `json:"webhook_url"`
	WebhookFailoverURL      string              `json:"webhook_failover_url"`
	WebhookAPIVersion       string              `json:"webhook_api_version"`
	WhitelistedDestinations []string            `json:"whitelisted_destinations"`
	NumberPoolSettings      *NumberPoolSettings `json:"number_pool_settings,omitempty"`
	URLShortenerSettings    *URLShortenerSettings `json:"url_shortener_settings,omitempty"`
	AlphaSender             *string             `json:"alpha_sender,omitempty"`
	CreatedAt               time.Time           `json:"created_at"`
	UpdatedAt               time.Time           `json:"updated_at"`
	V1Secret                string              `json:"v1_secret"`
}


type NumberPoolSettings struct {
	TollFreeWeight   float64 `json:"toll_free_weight"`
	LongCodeWeight   float64 `json:"long_code_weight"`
	SkipUnhealthy    bool    `json:"skip_unhealthy"`
	StickySender     bool    `json:"sticky_sender"`
	Geomatch         bool    `json:"geomatch"`
}

type URLShortenerSettings struct {
	Domain               string `json:"domain"`
	Prefix               string `json:"prefix"`
	ReplaceBlacklistOnly bool   `json:"replace_blacklist_only"`
	SendWebhooks         bool   `json:"send_webhooks"`
}

type OutboundVoiceProfile struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	ConnectionsCount        int       `json:"connections_count"`
	TrafficType             string    `json:"traffic_type"`
	ServicePlan             string    `json:"service_plan"`
	ConcurrentCallLimit     int       `json:"concurrent_call_limit"`
	Enabled                 bool      `json:"enabled"`
	Tags                    []string  `json:"tags"`
	UsagePaymentMethod      string    `json:"usage_payment_method"`
	WhitelistedDestinations []string  `json:"whitelisted_destinations"`
	MaxDestinationRate      float64   `json:"max_destination_rate"`
	DailySpendLimit         string    `json:"daily_spend_limit"`
	DailySpendLimitEnabled  bool      `json:"daily_spend_limit_enabled"`
	CallRecording           struct {
		Type                string   `json:"call_recording_type"`
		CallerPhoneNumbers  []string `json:"call_recording_caller_phone_numbers"`
		Channels            string   `json:"call_recording_channels"`
		Format              string   `json:"call_recording_format"`
	} `json:"call_recording"`
	BillingGroupID          string    `json:"billing_group_id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type PhoneNumberReservation struct {
	ID                string    `json:"id"`
	RecordType        string    `json:"record_type"`
	PhoneNumbers      []struct {
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
	ID                 string    `json:"id"`
	RecordType         string    `json:"record_type"`
	PhoneNumbers       []struct {
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

type FQDNConnection struct {
	ID                               string    `json:"id"`
	RecordType                       string    `json:"record_type"`
	Active                           bool      `json:"active"`
	AnchorsiteOverride               string    `json:"anchorsite_override"`
	ConnectionName                   string    `json:"connection_name"`
	TransportProtocol                string    `json:"transport_protocol"`
	DefaultOnHoldComfortNoiseEnabled bool      `json:"default_on_hold_comfort_noise_enabled"`
	DTMFType                         string    `json:"dtmf_type"`
	EncodeContactHeaderEnabled       bool      `json:"encode_contact_header_enabled"`
	EncryptedMedia                   string    `json:"encrypted_media"`
	OnnetT38PassthroughEnabled       bool      `json:"onnet_t38_passthrough_enabled"`
	IosPushCredentialID              string    `json:"ios_push_credential_id"`
	AndroidPushCredentialID          string    `json:"android_push_credential_id"`
	WebhookEventURL                  string    `json:"webhook_event_url"`
	WebhookEventFailoverURL          string    `json:"webhook_event_failover_url"`
	WebhookAPIVersion                string    `json:"webhook_api_version"`
	WebhookTimeoutSecs               int       `json:"webhook_timeout_secs"`
	RTCPSettings                     struct {
		Port                string `json:"port"`
		CaptureEnabled      bool   `json:"capture_enabled"`
		ReportFrequencySecs int    `json:"report_frequency_secs"`
	} `json:"rtcp_settings"`
	Inbound struct {
		ANINumberFormat              string   `json:"ani_number_format"`
		DNISNumberFormat             string   `json:"dnis_number_format"`
		Codecs                       []string `json:"codecs"`
		DefaultRoutingMethod         string   `json:"default_routing_method"`
		ChannelLimit                 int      `json:"channel_limit"`
		GenerateRingbackTone         bool     `json:"generate_ringback_tone"`
		ISUPHeadersEnabled           bool     `json:"isup_headers_enabled"`
		PRACKEnabled                 bool     `json:"prack_enabled"`
		PrivacyZoneEnabled           bool     `json:"privacy_zone_enabled"`
		SIPCompactHeadersEnabled     bool     `json:"sip_compact_headers_enabled"`
		SIPRegion                    string   `json:"sip_region"`
		SIPSubdomain                 string   `json:"sip_subdomain"`
		SIPSubdomainReceiveSettings  string   `json:"sip_subdomain_receive_settings"`
		Timeout1xxSecs               int      `json:"timeout_1xx_secs"`
		Timeout2xxSecs               int      `json:"timeout_2xx_secs"`
		ShakenSTIREnabled            bool     `json:"shaken_stir_enabled"`
	} `json:"inbound"`
	Outbound struct {
		ANIOverride              string `json:"ani_override"`
		ANIOverrideType          string `json:"ani_override_type"`
		CallParkingEnabled       bool   `json:"call_parking_enabled"`
		ChannelLimit             int    `json:"channel_limit"`
		GenerateRingbackTone     bool   `json:"generate_ringback_tone"`
		InstantRingbackEnabled   bool   `json:"instant_ringback_enabled"`
		IPAuthenticationMethod   string `json:"ip_authentication_method"`
		IPAuthenticationToken    string `json:"ip_authentication_token"`
		Localization             string `json:"localization"`
		OutboundVoiceProfileID   string `json:"outbound_voice_profile_id"`
		T38ReinviteSource        string `json:"t38_reinvite_source"`
		TechPrefix               string `json:"tech_prefix"`
		EncryptedMedia           string `json:"encrypted_media"`
		Timeout1xxSecs           int    `json:"timeout_1xx_secs"`
		Timeout2xxSecs           int    `json:"timeout_2xx_secs"`
	} `json:"outbound"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type PhoneNumber struct {
	ID                   string    `json:"id"`
	RecordType           string    `json:"record_type"`
	PhoneNumber          string    `json:"phone_number"`
	Status               string    `json:"status"`
	Tags                 []string  `json:"tags"`
	ExternalPin          string    `json:"external_pin"`
	ConnectionID         string    `json:"connection_id"`
	ConnectionName       string    `json:"connection_name"`
	CustomerReference    string    `json:"customer_reference"`
	MessagingProfileID   string    `json:"messaging_profile_id"`
	MessagingProfileName string    `json:"messaging_profile_name"`
	BillingGroupID       string    `json:"billing_group_id"`
	EmergencyEnabled     bool      `json:"emergency_enabled"`
	EmergencyAddressID   string    `json:"emergency_address_id"`
	CallForwardingEnabled bool     `json:"call_forwarding_enabled"`
	CNAMListingEnabled   bool      `json:"cnam_listing_enabled"`
	CallerIDNameEnabled  bool      `json:"caller_id_name_enabled"`
	CallRecordingEnabled bool      `json:"call_recording_enabled"`
	T38FaxGatewayEnabled bool      `json:"t38_fax_gateway_enabled"`
	NumberLevelRouting   string    `json:"number_level_routing"`
	PhoneNumberType      string    `json:"phone_number_type"`
	PurchasedAt          time.Time `json:"purchased_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	HDVoiceEnabled       bool      `json:"hd_voice_enabled"`
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

func (client *TelnyxClient) CreateOutboundVoiceProfile(name, trafficType, servicePlan string, concurrentCallLimit int, enabled bool, tags []string, usagePaymentMethod string, whitelistedDestinations []string, maxDestinationRate float64, dailySpendLimit string, dailySpendLimitEnabled bool, billingGroupID string) (*OutboundVoiceProfile, error) {
	body := map[string]interface{}{
		"name":                      name,
		"traffic_type":              trafficType,
		"service_plan":              servicePlan,
		"concurrent_call_limit":     concurrentCallLimit,
		"enabled":                   enabled,
		"tags":                      tags,
		"usage_payment_method":      usagePaymentMethod,
		"whitelisted_destinations":  whitelistedDestinations,
		"max_destination_rate":      maxDestinationRate,
		"daily_spend_limit":         dailySpendLimit,
		"daily_spend_limit_enabled": dailySpendLimitEnabled,
		"billing_group_id":          billingGroupID,
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

func (client *TelnyxClient) UpdateOutboundVoiceProfile(outboundVoiceProfileID, name, trafficType, servicePlan string, concurrentCallLimit int, enabled bool, tags []string, usagePaymentMethod string, whitelistedDestinations []string, maxDestinationRate float64, dailySpendLimit string, dailySpendLimitEnabled bool, billingGroupID string) (*OutboundVoiceProfile, error) {
	body := map[string]interface{}{
		"name":                      name,
		"traffic_type":              trafficType,
		"service_plan":              servicePlan,
		"concurrent_call_limit":     concurrentCallLimit,
		"enabled":                   enabled,
		"tags":                      tags,
		"usage_payment_method":      usagePaymentMethod,
		"whitelisted_destinations":  whitelistedDestinations,
		"max_destination_rate":      maxDestinationRate,
		"daily_spend_limit":         dailySpendLimit,
		"daily_spend_limit_enabled": dailySpendLimitEnabled,
		"billing_group_id":          billingGroupID,
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

func (client *TelnyxClient) CreateFQDNConnection(connectionName, transportProtocol, webhookEventURL, webhookEventFailoverURL, outboundVoiceProfileID string) (*FQDNConnection, error) {
	body := map[string]interface{}{
		"active":                             true,
		"anchorsite_override":                "Latency",
		"connection_name":                    connectionName,
		"transport_protocol":                 transportProtocol,
		"default_on_hold_comfort_noise_enabled": true,
		"dtmf_type":                          "RFC 2833",
		"encode_contact_header_enabled":      true,
		"encrypted_media":                    "SRTP",
		"onnet_t38_passthrough_enabled":      true,
		"webhook_event_url":                  webhookEventURL,
		"webhook_event_failover_url":         webhookEventFailoverURL,
		"webhook_api_version":                "1",
		"webhook_timeout_secs":               25,
		"outbound_voice_profile_id":          outboundVoiceProfileID,
	}
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("POST", "/fqdn_connections", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateFQDNConnection(connectionID, connectionName, transportProtocol, webhookEventURL, webhookEventFailoverURL, outboundVoiceProfileID string) (*FQDNConnection, error) {
	body := map[string]interface{}{
		"active":                             true,
		"anchorsite_override":                "Latency",
		"connection_name":                    connectionName,
		"transport_protocol":                 transportProtocol,
		"default_on_hold_comfort_noise_enabled": true,
		"dtmf_type":                          "RFC 2833",
		"encode_contact_header_enabled":      true,
		"encrypted_media":                    "SRTP",
		"onnet_t38_passthrough_enabled":      true,
		"webhook_event_url":                  webhookEventURL,
		"webhook_event_failover_url":         webhookEventFailoverURL,
		"webhook_api_version":                "1",
		"webhook_timeout_secs":               25,
		"outbound_voice_profile_id":          outboundVoiceProfileID,
	}
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/fqdn_connections/%s", connectionID), body, &result)
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
		"customer_reference":     customerReference,
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
		"phone_numbers":       phoneNumbersMap,
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
	phoneNumberID          string
	numberReservationID    string
	numberOrderID          string
}

func NewTestRunner(client *TelnyxClient) *TestRunner {
	return &TestRunner{client: client}
}

func (runner *TestRunner) PerformCreates() {
	// Create a Billing Group
	billingGroup, err := runner.client.CreateBillingGroup("Test Billing Group")
	if err != nil {
		fmt.Printf("Error creating billing group: %v\n", err)
		os.Exit(1)
	}
	runner.billingGroupID = billingGroup.ID
	fmt.Printf("Created Billing Group:\nID: %s\nName: %s\nCreated At: %s\n",
		billingGroup.ID, billingGroup.Name, billingGroup.CreatedAt)

	// Create an Outbound Voice Profile
	outboundVoiceProfile, err := runner.client.CreateOutboundVoiceProfile("test", "conversational", "global", 10, true, []string{"test-profile"}, "rate-deck", []string{"US"}, 10.0, "100.00", true, billingGroup.ID)
	if err != nil {
		fmt.Printf("Error creating outbound voice profile: %v\n", err)
		os.Exit(1)
	}
	runner.outboundVoiceProfileID = outboundVoiceProfile.ID
	fmt.Printf("Created Outbound Voice Profile:\nID: %s\nName: %s\nCreated At: %s\n",
		outboundVoiceProfile.ID, outboundVoiceProfile.Name, outboundVoiceProfile.CreatedAt)

	// Create a Messaging Profile
	messagingProfile, err := runner.client.CreateMessagingProfile("Test Profile", []string{"US"}, "https://www.example.com/hooks", "https://backup.example.com/hooks", true)
	if err != nil {
		fmt.Printf("Error creating messaging profile: %v\n", err)
		os.Exit(1)
	}
	runner.messagingProfileID = messagingProfile.ID
	fmt.Printf("Created Messaging Profile:\nID: %s\nName: %s\nEnabled: %t\nWebhook URL: %s\nWebhook Failover URL: %s\nWebhook API Version: %s\nWhitelisted Destinations: %v\nCreated At: %s\nUpdated At: %s\n",
		messagingProfile.ID, messagingProfile.Name, messagingProfile.Enabled, messagingProfile.WebhookURL, messagingProfile.WebhookFailoverURL, messagingProfile.WebhookAPIVersion, messagingProfile.WhitelistedDestinations, messagingProfile.CreatedAt, messagingProfile.UpdatedAt)

	// Create an FQDN Connection
	fqdnConnection, err := runner.client.CreateFQDNConnection("My FQDN Connection", "UDP", "https://www.example.com/hooks", "https://failover.example.com/hooks", outboundVoiceProfile.ID)
	if err != nil {
		fmt.Printf("Error creating FQDN connection: %v\n", err)
		os.Exit(1)
	}
	runner.fqdnConnectionID = fqdnConnection.ID
	fmt.Printf("Created FQDN Connection:\nID: %s\nName: %s\nCreated At: %s\n",
		fqdnConnection.ID, fqdnConnection.ConnectionName, fqdnConnection.CreatedAt)
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
	updatedOutboundVoiceProfile, err := runner.client.UpdateOutboundVoiceProfile(runner.outboundVoiceProfileID, "office-updated", "conversational", "global", 10, true, []string{"office-profile-updated"}, "rate-deck", []string{"US"}, 10.0, "200.00", true, runner.billingGroupID)
	if err != nil {
		fmt.Printf("Error updating outbound voice profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated Outbound Voice Profile:\nID: %s\nName: %s\nUpdated At: %s\n",
		updatedOutboundVoiceProfile.ID, updatedOutboundVoiceProfile.Name, updatedOutboundVoiceProfile.UpdatedAt)

	// Update the FQDN Connection
	updatedFQDNConnection, err := runner.client.UpdateFQDNConnection(runner.fqdnConnectionID, "Updated FQDN Connection", "TCP", "https://www.example.com/hooks", "https://failover.example.com/hooks", runner.outboundVoiceProfileID)
	if err != nil {
		fmt.Printf("Error updating FQDN connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated FQDN Connection:\nID: %s\nName: %s\nUpdated At: %s\n",
		updatedFQDNConnection.ID, updatedFQDNConnection.ConnectionName, updatedFQDNConnection.UpdatedAt)
}

func (runner *TestRunner) PerformCascadingDeletes() {
	// Delete the FQDN Connection
	err := runner.client.DeleteFQDNConnection(runner.fqdnConnectionID)
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

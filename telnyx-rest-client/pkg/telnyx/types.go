package telnyx

import "time"

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

type CreateNumberOrderRequest struct {
	PhoneNumbers       []PhoneNumberRequest `json:"phone_numbers"`
	ConnectionID       string               `json:"connection_id"`
	MessagingProfileID string               `json:"messaging_profile_id"`
	BillingGroupID     string               `json:"billing_group_id"`
	CustomerReference  string               `json:"customer_reference"`
	SubNumberOrderIDs  []string             `json:"sub_number_orders_ids,omitempty"`
}

type PhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// CreateBillingGroupRequest represents the request payload for creating a billing group.
type CreateBillingGroupRequest struct {
	Name string `json:"name"`
}

// UpdateBillingGroupRequest represents the request payload for updating a billing group.
type UpdateBillingGroupRequest struct {
	Name string `json:"name"`
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
	MaxDestinationRate      *float64       `json:"max_destination_rate"`
	DailySpendLimit         *string        `json:"daily_spend_limit"`
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

// SubNumberOrderRegulatoryRequirement represents a regulatory requirement for a sub number order
type SubNumberOrderRegulatoryRequirement struct {
	RecordType    string `json:"record_type"`
	RequirementID string `json:"requirement_id"`
	FieldType     string `json:"field_type"`
}

// SubNumberOrderResponse represents the details of a sub number order
type SubNumberOrderResponse struct {
	ID                     string                                `json:"id"`
	OrderRequestID         string                                `json:"order_request_id"`
	CountryCode            string                                `json:"country_code"`
	PhoneNumberType        string                                `json:"phone_number_type"`
	UserID                 string                                `json:"user_id"`
	RegulatoryRequirements []SubNumberOrderRegulatoryRequirement `json:"regulatory_requirements"`
	RecordType             string                                `json:"record_type"`
	PhoneNumbersCount      int                                   `json:"phone_numbers_count"`
	CreatedAt              time.Time                             `json:"created_at"`
	UpdatedAt              time.Time                             `json:"updated_at"`
	RequirementsMet        bool                                  `json:"requirements_met"`
	Status                 string                                `json:"status"`
	CustomerReference      string                                `json:"customer_reference"`
	IsBlockSubNumberOrder  bool                                  `json:"is_block_sub_number_order"`
}

// type PhoneNumberOrderRequest struct {
// 	PhoneNumbers       []string `json:"phone_numbers"`
// 	ConnectionID       string   `json:"connection_id"`
// 	MessagingProfileID string   `json:"messaging_profile_id"`
// 	BillingGroupID     string   `json:"billing_group_id"`
// 	CustomerReference  string   `json:"customer_reference"`
// 	SubNumberOrderIDs  []string `json:"sub_number_order_ids,omitempty"`
// }

type NumberOrderRegulatoryRequirement struct {
	RequirementID string `json:"requirement_id"`
	FieldValue    string `json:"field_value"`
	FieldType     string `json:"field_type"`
}

type UpdateNumberOrderRequest struct {
	CustomerReference      string                             `json:"customer_reference"`
	RegulatoryRequirements []NumberOrderRegulatoryRequirement `json:"regulatory_requirements"`
}

type PhoneNumberOrderResponse struct {
	ID                 string                      `json:"id"`
	RecordType         string                      `json:"record_type"`
	PhoneNumbersCount  int                         `json:"phone_numbers_count"`
	ConnectionID       string                      `json:"connection_id"`
	MessagingProfileID string                      `json:"messaging_profile_id"`
	BillingGroupID     string                      `json:"billing_group_id"`
	PhoneNumbers       []OrderResponsePhoneNumbers `json:"phone_numbers"`
	SubNumberOrderIDs  []string                    `json:"sub_number_orders_ids,omitempty"`
	Status             string                      `json:"status"`
	CustomerReference  string                      `json:"customer_reference"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
	RequirementsMet    bool                        `json:"requirements_met"`
}

type OrderResponsePhoneNumbers struct {
	ID                     string                             `json:"id"`
	RecordType             string                             `json:"record_type"`
	PhoneNumber            string                             `json:"phone_number"`
	BundleID               string                             `json:"bundle_id"`
	PhoneNumberType        string                             `json:"phone_number_type"`
	CountryCode            string                             `json:"country_code"`
	RequirementsMet        bool                               `json:"requirements_met"`
	Status                 string                             `json:"status"`
	RegulatoryRequirements []NumberOrderRegulatoryRequirement `json:"regulatory_requirements"`
}

// PhoneNumberResponse represents the response for retrieving a phone number.
type PhoneNumberResponse struct {
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
	ChannelLimit                *int     `json:"channel_limit,omitempty"`
	GenerateRingbackTone        *bool    `json:"generate_ringback_tone,omitempty"`
	ISUPHeadersEnabled          *bool    `json:"isup_headers_enabled,omitempty"`
	PRACKEnabled                *bool    `json:"prack_enabled,omitempty"`
	PrivacyZoneEnabled          *bool    `json:"privacy_zone_enabled,omitempty"`
	SIPCompactHeadersEnabled    *bool    `json:"sip_compact_headers_enabled,omitempty"`
	SIPRegion                   string   `json:"sip_region"`
	SIPSubdomain                string   `json:"sip_subdomain"`
	SIPSubdomainReceiveSettings string   `json:"sip_subdomain_receive_settings"`
	Timeout1xxSecs              *int     `json:"timeout_1xx_secs,omitempty"`
	Timeout2xxSecs              *int     `json:"timeout_2xx_secs,omitempty"`
	ShakenSTIREnabled           *bool    `json:"shaken_stir_enabled,omitempty"`
}

type OutboundSettings struct {
	ANIOverride            string  `json:"ani_override"`
	ANIOverrideType        string  `json:"ani_override_type"`
	CallParkingEnabled     *bool   `json:"call_parking_enabled,omitempty"`
	ChannelLimit           *int    `json:"channel_limit,omitempty"`
	GenerateRingbackTone   *bool   `json:"generate_ringback_tone,omitempty"`
	InstantRingbackEnabled *bool   `json:"instant_ringback_enabled,omitempty"`
	IPAuthenticationMethod string  `json:"ip_authentication_method"`
	IPAuthenticationToken  *string `json:"ip_authentication_token,omitempty"`
	Localization           string  `json:"localization"`
	OutboundVoiceProfileID string  `json:"outbound_voice_profile_id"`
	T38ReinviteSource      string  `json:"t38_reinvite_source"`
}

// New FQDN Struct
type FQDN struct {
	ID            string    `json:"id"`
	ConnectionID  int       `json:"connection_id"`
	FQDN          string    `json:"fqdn"`
	Port          int       `json:"port"`
	DNSRecordType string    `json:"dns_record_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FQDNConnection Struct with Credential Fields
type FQDNConnection struct {
	ID string `json:"id"`
	// RecordType                       string           `json:"record_type"`
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
	MicrosoftTeamsSbc                bool             `json:"microsoft_teams_sbc"`
	WebhookEventURL                  string           `json:"webhook_event_url"`
	WebhookEventFailoverURL          string           `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion                string           `json:"webhook_api_version"`
	WebhookTimeoutSecs               int              `json:"webhook_timeout_secs,omitempty"`
	RTCPSettings                     RTCPSettings     `json:"rtcp_settings"`
	Inbound                          InboundSettings  `json:"inbound"`
	Outbound                         OutboundSettings `json:"outbound"`
	CreatedAt                        time.Time        `json:"created_at"`
	UpdatedAt                        time.Time        `json:"updated_at"`
	Username                         *string          `json:"user_name,omitempty"`
	Password                         *string          `json:"password,omitempty"`
	SipUriCallingPreference          *string          `json:"sip_uri_calling_preference,omitempty"`
}

// CredentialConnection Struct
type CredentialConnection struct {
	ID                               string           `json:"id"`
	RecordType                       string           `json:"record_type"`
	Active                           bool             `json:"active"`
	AnchorsiteOverride               string           `json:"anchorsite_override"`
	ConnectionName                   string           `json:"connection_name"`
	DefaultOnHoldComfortNoiseEnabled bool             `json:"default_on_hold_comfort_noise_enabled"`
	DTMFType                         string           `json:"dtmf_type"`
	EncodeContactHeaderEnabled       bool             `json:"encode_contact_header_enabled"`
	OnnetT38PassthroughEnabled       bool             `json:"onnet_t38_passthrough_enabled"`
	MicrosoftTeamsSbc                bool             `json:"microsoft_teams_sbc"`
	WebhookEventURL                  string           `json:"webhook_event_url"`
	WebhookEventFailoverURL          string           `json:"webhook_event_failover_url"`
	WebhookAPIVersion                string           `json:"webhook_api_version"`
	WebhookTimeoutSecs               int              `json:"webhook_timeout_secs"`
	RTCPSettings                     RTCPSettings     `json:"rtcp_settings"`
	Inbound                          InboundSettings  `json:"inbound"`
	Outbound                         OutboundSettings `json:"outbound"`
	CreatedAt                        time.Time        `json:"created_at"`
	UpdatedAt                        time.Time        `json:"updated_at"`
	Username                         string           `json:"user_name"`
	Password                         string           `json:"password"`
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

// AvailablePhoneNumbersResponse represents the response from the list available phone numbers API.
type AvailablePhoneNumbersResponse struct {
	Data []AvailablePhoneNumber `json:"data"`
	Meta struct {
		TotalResults      int `json:"total_results"`
		BestEffortResults int `json:"best_effort_results"`
	} `json:"meta"`
}

// AvailablePhoneNumber represents a single available phone number in the response.
type AvailablePhoneNumber struct {
	RecordType        string              `json:"record_type"`
	PhoneNumber       string              `json:"phone_number"`
	VanityFormat      string              `json:"vanity_format"`
	BestEffort        bool                `json:"best_effort"`
	Quickship         bool                `json:"quickship"`
	Reservable        bool                `json:"reservable"`
	RegionInformation []RegionInformation `json:"region_information"`
	CostInformation   CostInformation     `json:"cost_information"`
	Features          []Feature           `json:"features"`
}

// RegionInformation represents the region information for an available phone number.
type RegionInformation struct {
	RegionType string `json:"region_type"`
	RegionName string `json:"region_name"`
}

// CostInformation represents the cost information for an available phone number.
type CostInformation struct {
	UpfrontCost string `json:"upfront_cost"`
	MonthlyCost string `json:"monthly_cost"`
	Currency    string `json:"currency"`
}

// Feature represents a feature available for a phone number.
type Feature struct {
	Name string `json:"name"`
}

// AvailablePhoneNumbersRequest encapsulates the filters for listing available phone numbers.
type AvailablePhoneNumbersRequest struct {
	StartsWith              string   `json:"filter[phone_number][starts_with],omitempty"`
	EndsWith                string   `json:"filter[phone_number][ends_with],omitempty"`
	Contains                string   `json:"filter[phone_number][contains],omitempty"`
	Locality                string   `json:"filter[locality],omitempty"`
	AdministrativeArea      string   `json:"filter[administrative_area],omitempty"`
	CountryCode             string   `json:"filter[country_code],omitempty"`
	NationalDestinationCode string   `json:"filter[national_destination_code],omitempty"`
	RateCenter              string   `json:"filter[rate_center],omitempty"`
	PhoneNumberType         string   `json:"filter[phone_number_type],omitempty"`
	Features                []string `json:"filter[features],omitempty"`
	Limit                   int      `json:"filter[limit],omitempty"`
	BestEffort              bool     `json:"filter[best_effort],omitempty"`
	Quickship               bool     `json:"filter[quickship],omitempty"`
	Reservable              bool     `json:"filter[reservable],omitempty"`
	ExcludeHeldNumbers      bool     `json:"filter[exclude_held_numbers],omitempty"`
}

// UpdatePhoneNumberRequest represents the request payload for updating a phone number.
type UpdatePhoneNumberRequest struct {
	CustomerReference  string   `json:"customer_reference"`
	ConnectionID       int      `json:"connection_id"`
	BillingGroupID     string   `json:"billing_group_id"`
	Tags               []string `json:"tags"`
	HDVoiceEnabled     bool     `json:"hd_voice_enabled"`
	ExternalPin        string   `json:"external_pin,omitempty"`
	NumberLevelRouting string   `json:"number_level_routing,omitempty"`
}

// UpdatePhoneNumberResponse represents the response from updating a phone number.
type UpdatePhoneNumberResponse struct {
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

// TeXMLApplication represents a TeXML application in Telnyx.
type TeXMLApplication struct {
	ID                      string           `json:"id"`
	RecordType              string           `json:"record_type"`
	FriendlyName            string           `json:"friendly_name"`
	Active                  bool             `json:"active"`
	AnchorsiteOverride      string           `json:"anchorsite_override"`
	DTMFType                string           `json:"dtmf_type"`
	FirstCommandTimeout     bool             `json:"first_command_timeout"`
	FirstCommandTimeoutSecs int              `json:"first_command_timeout_secs"`
	VoiceURL                string           `json:"voice_url"`
	VoiceFallbackURL        string           `json:"voice_fallback_url"`
	VoiceMethod             string           `json:"voice_method"`
	StatusCallback          string           `json:"status_callback"`
	StatusCallbackMethod    string           `json:"status_callback_method"`
	Inbound                 InboundSettings  `json:"inbound"`
	Outbound                OutboundSettings `json:"outbound"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
}

// TeXMLApplicationRequest represents the request payload for creating or updating a TeXML application.
type TeXMLApplicationRequest struct {
	FriendlyName            string           `json:"friendly_name"`
	Active                  bool             `json:"active"`
	AnchorsiteOverride      string           `json:"anchorsite_override"`
	DTMFType                string           `json:"dtmf_type"`
	FirstCommandTimeout     bool             `json:"first_command_timeout"`
	FirstCommandTimeoutSecs int              `json:"first_command_timeout_secs"`
	VoiceURL                string           `json:"voice_url"`
	VoiceFallbackURL        string           `json:"voice_fallback_url"`
	VoiceMethod             string           `json:"voice_method"`
	StatusCallback          string           `json:"status_callback"`
	StatusCallbackMethod    string           `json:"status_callback_method"`
	Inbound                 InboundSettings  `json:"inbound"`
	Outbound                OutboundSettings `json:"outbound"`
}

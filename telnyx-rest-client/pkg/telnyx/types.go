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
	ConnectionID  int       `json:"connection_id"`
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
	MicrosoftTeamsSbc                bool             `json:"microsoft_teams_sbc`
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
	FQDNOutboundAuthentication       string           `json:"fqdn_outbound_authentication`
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
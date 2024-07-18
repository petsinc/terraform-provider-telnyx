package test_runner

import (
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type TestRunner struct {
	client                 *telnyx.TelnyxClient
	logger                 *zap.Logger
	billingGroupID         string
	outboundVoiceProfileID string
	messagingProfileID     string
	fqdnConnectionID       string
	fqdnConnectionIDInt    int
	fqdnID                 string
	phoneNumberID          string
	numberReservationID    string
	numberOrderID          string
	credentialConnectionID string
}

func NewTestRunner(client *telnyx.TelnyxClient, logger *zap.Logger) *TestRunner {
	return &TestRunner{client: client, logger: logger}
}

func (runner *TestRunner) PerformCreates() {
	runner.logger.Info("Performing create operations")

	// Create a Billing Group
	billingGroup, err := runner.client.CreateBillingGroup("Test Billing Group")
	if err != nil {
		runner.logger.Error("Error creating billing group", zap.Error(err))
		os.Exit(1)
	}
	runner.billingGroupID = billingGroup.ID
	runner.logger.Info("Created Billing Group",
		zap.String("ID", billingGroup.ID),
		zap.String("Name", billingGroup.Name),
		zap.Time("Created At", billingGroup.CreatedAt))

	// Create an Outbound Voice Profile
	outboundVoiceProfile, err := runner.client.CreateOutboundVoiceProfile(telnyx.OutboundVoiceProfile{
		Name:                    "Test Outbound Profile",
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
		CallRecording: telnyx.CallRecording{
			Type:               "all",
			CallerPhoneNumbers: []string{},
			Channels:           "single",
			Format:             "wav",
		},
	})

	if err != nil {
		runner.logger.Error("Error creating outbound voice profile", zap.Error(err))
		os.Exit(1)
	}
	runner.outboundVoiceProfileID = outboundVoiceProfile.ID
	runner.logger.Info("Created Outbound Voice Profile",
		zap.String("ID", outboundVoiceProfile.ID),
		zap.String("Name", outboundVoiceProfile.Name),
		zap.Time("Created At", outboundVoiceProfile.CreatedAt))

	// Create a Messaging Profile
	messagingProfile, err := runner.client.CreateMessagingProfile(telnyx.MessagingProfile{
		Name:                    "Test Profile",
		Enabled:                 true,
		WebhookURL:              "https://www.example.com/hooks",
		WebhookFailoverURL:      "https://backup.example.com/hooks",
		WebhookAPIVersion:       "2",
		WhitelistedDestinations: []string{"US"},
	})
	if err != nil {
		runner.logger.Error("Error creating messaging profile", zap.Error(err))
		os.Exit(1)
	}
	runner.messagingProfileID = messagingProfile.ID
	runner.logger.Info("Created Messaging Profile",
		zap.String("ID", messagingProfile.ID),
		zap.String("Name", messagingProfile.Name),
		zap.Bool("Enabled", messagingProfile.Enabled),
		zap.String("Webhook URL", messagingProfile.WebhookURL),
		zap.String("Webhook Failover URL", messagingProfile.WebhookFailoverURL),
		zap.String("Webhook API Version", messagingProfile.WebhookAPIVersion),
		zap.Strings("Whitelisted Destinations", messagingProfile.WhitelistedDestinations),
		zap.Time("Created At", messagingProfile.CreatedAt),
		zap.Time("Updated At", messagingProfile.UpdatedAt))

	// Create a Credential Connection
	credentialConnection, err := runner.client.CreateCredentialConnection(telnyx.FQDNConnection{
		Username:                         telnyx.StringPtr("hellopatienttest12345"),
		Password:                         telnyx.StringPtr("54321testpatienthello"),
		Active:                           true,
		AnchorsiteOverride:               "Latency",
		ConnectionName:                   "Test Credential Connection",
		TransportProtocol:                "UDP",
		DefaultOnHoldComfortNoiseEnabled: true,
		DTMFType:                         "RFC 2833",
		EncodeContactHeaderEnabled:       false,
		EncryptedMedia:                   nil,
		OnnetT38PassthroughEnabled:       false,
		MicrosoftTeamsSbc:                false,
		WebhookEventURL:                  "https://www.example.com/hooks",
		WebhookEventFailoverURL:          "https://failover.example.com/hooks",
		WebhookAPIVersion:                "1",
		WebhookTimeoutSecs:               25,
		RTCPSettings: telnyx.RTCPSettings{
			Port:                "rtp+1",
			CaptureEnabled:      false,
			ReportFrequencySecs: 5,
		},
		Inbound: telnyx.InboundSettings{
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
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            "+12345678901",
			ANIOverrideType:        "always",
			CallParkingEnabled:     true,
			ChannelLimit:           10,
			GenerateRingbackTone:   true,
			InstantRingbackEnabled: false,
			IPAuthenticationMethod: "token",
			IPAuthenticationToken:  "aBcD1234aBcD1234",
			Localization:           "US",
			OutboundVoiceProfileID: runner.outboundVoiceProfileID,
			T38ReinviteSource:      "customer",
			EncryptedMedia:         "SRTP",
			Timeout1xxSecs:         3,
			Timeout2xxSecs:         90,
		},
	})
	if err != nil {
		runner.logger.Error("Error creating credential connection", zap.Error(err))
		os.Exit(1)
	}
	runner.credentialConnectionID = credentialConnection.ID
	runner.logger.Info("Created Credential Connection",
		zap.String("ID", credentialConnection.ID),
		zap.String("Name", credentialConnection.ConnectionName),
		zap.Time("Created At", credentialConnection.CreatedAt))

	// Create an FQDN Connection
	fqdnConnection, err := runner.client.CreateFQDNConnection(telnyx.FQDNConnection{
		Username:                         telnyx.StringPtr("hellopatienttest123456"),
		Password:                         telnyx.StringPtr("54321testpatienthello"),
		Active:                           true,
		AnchorsiteOverride:               "Latency",
		ConnectionName:                   "Test FQDN Connection",
		TransportProtocol:                "UDP",
		DefaultOnHoldComfortNoiseEnabled: true,
		DTMFType:                         "RFC 2833",
		EncodeContactHeaderEnabled:       false,
		EncryptedMedia:                   nil,
		OnnetT38PassthroughEnabled:       false,
		MicrosoftTeamsSbc:                false,
		WebhookEventURL:                  "https://www.example.com/hooks",
		WebhookEventFailoverURL:          "https://failover.example.com/hooks",
		WebhookAPIVersion:                "1",
		WebhookTimeoutSecs:               25,
		RTCPSettings: telnyx.RTCPSettings{
			Port:                "rtp+1",
			CaptureEnabled:      false,
			ReportFrequencySecs: 5,
		},
		Inbound: telnyx.InboundSettings{
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
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            "+12345678901",
			ANIOverrideType:        "always",
			CallParkingEnabled:     true,
			ChannelLimit:           10,
			GenerateRingbackTone:   true,
			InstantRingbackEnabled: false,
			IPAuthenticationMethod: "token",
			IPAuthenticationToken:  "aBcD1234aBcD1234",
			Localization:           "US",
			OutboundVoiceProfileID: runner.outboundVoiceProfileID,
			T38ReinviteSource:      "customer",
			EncryptedMedia:         "SRTP",
			Timeout1xxSecs:         3,
			Timeout2xxSecs:         90,
		},
	})
	if err != nil {
		runner.logger.Error("Error creating FQDN connection", zap.Error(err))
		os.Exit(1)
	}
	runner.fqdnConnectionID = fqdnConnection.ID
	asInt, err := strconv.Atoi(runner.fqdnConnectionID)
	runner.fqdnConnectionIDInt = asInt
	runner.logger.Info("Created FQDN Connection",
		zap.String("ID", fqdnConnection.ID),
		zap.String("Name", fqdnConnection.ConnectionName),
		zap.Time("Created At", fqdnConnection.CreatedAt))

	// Create an FQDN and bind it to the connection
	fqdn, err := runner.client.CreateFQDN(telnyx.FQDN{
		ConnectionID:  runner.fqdnConnectionIDInt,
		FQDN:          "test.sip.livekit.cloud",
		DNSRecordType: "a",
		Port:          5060,
	})
	if err != nil {
		runner.logger.Error("Error creating FQDN", zap.Error(err))
		os.Exit(1)
	}
	runner.fqdnID = fqdn.ID
	runner.logger.Info("Created FQDN",
		zap.String("ID", fqdn.ID),
		zap.String("FQDN", fqdn.FQDN),
		zap.Time("Created At", fqdn.CreatedAt))

	// List all available phone numbers
	filters := telnyx.AvailablePhoneNumbersRequest{
		StartsWith:  "312",
		CountryCode: "US",
		Limit:       10,
	}
	response, err := runner.client.ListAvailablePhoneNumbers(filters)
	if err != nil {
		runner.logger.Error("Error retrieving available phone numbers", zap.Error(err))
		os.Exit(1)
	}

	if len(response.Data) == 0 {
		runner.logger.Info("No available phone numbers found")
		return
	}

	firstPhoneNumber := response.Data[0].PhoneNumber
	runner.logger.Info("First Available Phone Number",
		zap.String("PhoneNumber", firstPhoneNumber),
		zap.Bool("Reservable", response.Data[0].Reservable),
		zap.String("UpfrontCost", response.Data[0].CostInformation.UpfrontCost),
		zap.String("MonthlyCost", response.Data[0].CostInformation.MonthlyCost))

	// Create a number order for the first available phone number
	numberOrderRequest := telnyx.CreateNumberOrderRequest{
		PhoneNumbers:       []telnyx.PhoneNumberRequest{{PhoneNumber: firstPhoneNumber}},
		ConnectionID:       runner.fqdnConnectionID,
		MessagingProfileID: runner.messagingProfileID,
		BillingGroupID:     runner.billingGroupID,
		CustomerReference:  "Test Order",
	}

	numberOrder, err := runner.client.CreateNumberOrder(numberOrderRequest)
	if err != nil {
		runner.logger.Error("Error creating number order", zap.Error(err))
		os.Exit(1)
	}

	runner.numberOrderID = numberOrder.ID
	runner.phoneNumberID = numberOrder.PhoneNumbers[0].ID
	runner.logger.Info("Created Number Order",
		zap.String("ID", numberOrder.ID),
		zap.Int("PhoneNumbersCount", numberOrder.PhoneNumbersCount),
		zap.String("Status", numberOrder.Status),
		zap.String("CustomerReference", numberOrder.CustomerReference),
		zap.Time("CreatedAt", numberOrder.CreatedAt),
		zap.Time("UpdatedAt", numberOrder.UpdatedAt))
}

func (runner *TestRunner) PerformUpdates() {
	runner.logger.Info("Performing update operations")

	// Update the Messaging Profile
	updatedMessagingProfile, err := runner.client.UpdateMessagingProfile(runner.messagingProfileID, telnyx.MessagingProfile{
		Name:                    "Updated Profile for Messages",
		Enabled:                 true,
		WebhookURL:              "https://www.example.com/hooks",
		WebhookFailoverURL:      "https://backup.example.com/hooks",
		WebhookAPIVersion:       "2",
		WhitelistedDestinations: []string{"US"},
	})
	if err != nil {
		runner.logger.Error("Error updating messaging profile", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Updated Messaging Profile",
		zap.String("ID", updatedMessagingProfile.ID),
		zap.String("Name", updatedMessagingProfile.Name),
		zap.Bool("Enabled", updatedMessagingProfile.Enabled),
		zap.String("Webhook URL", updatedMessagingProfile.WebhookURL),
		zap.String("Webhook Failover URL", updatedMessagingProfile.WebhookFailoverURL),
		zap.String("Webhook API Version", updatedMessagingProfile.WebhookAPIVersion),
		zap.Strings("Whitelisted Destinations", updatedMessagingProfile.WhitelistedDestinations),
		zap.Time("Created At", updatedMessagingProfile.CreatedAt),
		zap.Time("Updated At", updatedMessagingProfile.UpdatedAt))

	// Update the Outbound Voice Profile
	updatedOutboundVoiceProfile, err := runner.client.UpdateOutboundVoiceProfile(runner.outboundVoiceProfileID, telnyx.OutboundVoiceProfile{
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
		CallRecording: telnyx.CallRecording{
			Type:               "all",
			CallerPhoneNumbers: []string{},
			Channels:           "single",
			Format:             "wav",
		},
	})

	if err != nil {
		runner.logger.Error("Error updating outbound voice profile", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Updated Outbound Voice Profile",
		zap.String("ID", updatedOutboundVoiceProfile.ID),
		zap.String("Name", updatedOutboundVoiceProfile.Name),
		zap.Time("Updated At", updatedOutboundVoiceProfile.UpdatedAt))

	// Update the Credential Connection
	updatedCredentialConnection, err := runner.client.UpdateCredentialConnection(runner.credentialConnectionID, telnyx.FQDNConnection{
		Username:                         telnyx.StringPtr("updatedtest12345"),
		Password:                         telnyx.StringPtr("updatedpassword54321"),
		Active:                           true,
		AnchorsiteOverride:               "Latency",
		ConnectionName:                   "Updated Credential Connection",
		TransportProtocol:                "UDP",
		DefaultOnHoldComfortNoiseEnabled: true,
		DTMFType:                         "RFC 2833",
		EncodeContactHeaderEnabled:       false,
		EncryptedMedia:                   nil,
		OnnetT38PassthroughEnabled:       false,
		MicrosoftTeamsSbc:                false,
		WebhookEventURL:                  "https://www.example.com/hooks",
		WebhookEventFailoverURL:          "https://failover.example.com/hooks",
		WebhookAPIVersion:                "1",
		WebhookTimeoutSecs:               25,
		RTCPSettings: telnyx.RTCPSettings{
			Port:                "rtp+1",
			CaptureEnabled:      false,
			ReportFrequencySecs: 5,
		},
		Inbound: telnyx.InboundSettings{
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
			SIPSubdomain:                "updatedexample.sip.telnyx.com",
			SIPSubdomainReceiveSettings: "only_my_connections",
			Timeout1xxSecs:              3,
			Timeout2xxSecs:              90,
			ShakenSTIREnabled:           true,
		},
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            "+12345678901",
			ANIOverrideType:        "always",
			CallParkingEnabled:     true,
			ChannelLimit:           10,
			GenerateRingbackTone:   true,
			InstantRingbackEnabled: false,
			IPAuthenticationMethod: "token",
			IPAuthenticationToken:  "updatedtoken1234",
			Localization:           "US",
			OutboundVoiceProfileID: runner.outboundVoiceProfileID,
			T38ReinviteSource:      "customer",
			EncryptedMedia:         "SRTP",
			Timeout1xxSecs:         3,
			Timeout2xxSecs:         90,
		},
	})
	if err != nil {
		runner.logger.Error("Error updating credential connection", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Updated Credential Connection",
		zap.String("ID", updatedCredentialConnection.ID),
		zap.String("Name", updatedCredentialConnection.ConnectionName),
		zap.Time("Updated At", updatedCredentialConnection.UpdatedAt))

	// Update the FQDN Connection
	updatedFQDNConnection, err := runner.client.UpdateFQDNConnection(runner.fqdnConnectionID, telnyx.FQDNConnection{
		Username:                         telnyx.StringPtr("updatedtest123456"),
		Password:                         telnyx.StringPtr("updatedpasswordhello"),
		Active:                           true,
		AnchorsiteOverride:               "Latency",
		ConnectionName:                   "Updated FQDN Connection",
		TransportProtocol:                "UDP",
		DefaultOnHoldComfortNoiseEnabled: true,
		DTMFType:                         "RFC 2833",
		EncodeContactHeaderEnabled:       false,
		EncryptedMedia:                   nil,
		OnnetT38PassthroughEnabled:       false,
		MicrosoftTeamsSbc:                false,
		WebhookEventURL:                  "https://www.example.com/hooks",
		WebhookEventFailoverURL:          "https://failover.example.com/hooks",
		WebhookAPIVersion:                "1",
		WebhookTimeoutSecs:               25,
		RTCPSettings: telnyx.RTCPSettings{
			Port:                "rtp+1",
			CaptureEnabled:      false,
			ReportFrequencySecs: 5,
		},
		Inbound: telnyx.InboundSettings{
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
			SIPSubdomain:                "updatedexample.sip.telnyx.com",
			SIPSubdomainReceiveSettings: "only_my_connections",
			Timeout1xxSecs:              3,
			Timeout2xxSecs:              90,
			ShakenSTIREnabled:           true,
		},
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            "+12345678901",
			ANIOverrideType:        "always",
			CallParkingEnabled:     true,
			ChannelLimit:           10,
			GenerateRingbackTone:   true,
			InstantRingbackEnabled: false,
			IPAuthenticationMethod: "token",
			IPAuthenticationToken:  "updatedtoken1234",
			Localization:           "US",
			OutboundVoiceProfileID: runner.outboundVoiceProfileID,
			T38ReinviteSource:      "customer",
			EncryptedMedia:         "SRTP",
			Timeout1xxSecs:         3,
			Timeout2xxSecs:         90,
		},
	})
	if err != nil {
		runner.logger.Error("Error updating FQDN connection", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Updated FQDN Connection",
		zap.String("ID", updatedFQDNConnection.ID),
		zap.String("Name", updatedFQDNConnection.ConnectionName),
		zap.Time("Updated At", updatedFQDNConnection.UpdatedAt))

	// Update the FQDN
	updatedFQDN, err := runner.client.UpdateFQDN(runner.fqdnID, telnyx.FQDN{
		ConnectionID:  runner.fqdnConnectionIDInt,
		FQDN:          "updated.test.sip.livekit.cloud",
		DNSRecordType: "a",
		Port:          5060,
	})
	if err != nil {
		runner.logger.Error("Error updating FQDN", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Updated FQDN",
		zap.String("ID", updatedFQDN.ID),
		zap.String("FQDN", updatedFQDN.FQDN),
		zap.Time("Updated At", updatedFQDN.UpdatedAt))

	// Update the Number Order
	updateRequest := telnyx.UpdateNumberOrderRequest{
		CustomerReference: "Updated Test Order",
		RegulatoryRequirements: []telnyx.NumberOrderRegulatoryRequirement{
			{
				RequirementID: "requirement_id_example",
				FieldValue:    "field_value_example",
			},
		},
	}

	updatedNumberOrder, err := runner.client.UpdateNumberOrder(runner.numberOrderID, updateRequest)
	if err != nil {
		runner.logger.Error("Error updating number order", zap.Error(err))
		os.Exit(1)
	}

	runner.logger.Info("Updated Number Order",
		zap.String("ID", updatedNumberOrder.ID),
		zap.Int("PhoneNumbersCount", updatedNumberOrder.PhoneNumbersCount),
		zap.String("Status", updatedNumberOrder.Status),
		zap.String("CustomerReference", updatedNumberOrder.CustomerReference),
		zap.Time("CreatedAt", updatedNumberOrder.CreatedAt),
		zap.Time("UpdatedAt", updatedNumberOrder.UpdatedAt))

	// Update the Phone Number
	phoneNumberUpdateRequest := telnyx.UpdatePhoneNumberRequest{
		CustomerReference:  "Updated Test Number",
		ConnectionID:       runner.fqdnConnectionIDInt,
		BillingGroupID:     runner.billingGroupID,
		Tags:               []string{"test", "updated"},
		HDVoiceEnabled:     true,
		ExternalPin:        "1234",
		NumberLevelRouting: "ENABLED",
	}

	updatedPhoneNumber, err := runner.client.UpdatePhoneNumber(runner.phoneNumberID, phoneNumberUpdateRequest)
	if err != nil {
		runner.logger.Error("Error updating phone number", zap.Error(err))
		os.Exit(1)
	}

	runner.logger.Info("Updated Phone Number",
		zap.String("ID", updatedPhoneNumber.ID),
		zap.String("PhoneNumber", updatedPhoneNumber.PhoneNumber),
		zap.String("CustomerReference", updatedPhoneNumber.CustomerReference),
		zap.Time("UpdatedAt", updatedPhoneNumber.UpdatedAt))
}

func (runner *TestRunner) PerformCascadingDeletes() {
	runner.logger.Info("Performing cascading delete operations")

	// Delete the phone number
	err := runner.client.DeletePhoneNumber(runner.phoneNumberID)
	if err != nil {
		runner.logger.Error("Error deleting phone number", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted Phone Number")

	// Delete the FQDN
	err = runner.client.DeleteFQDN(runner.fqdnID)
	if err != nil {
		runner.logger.Error("Error deleting FQDN", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted FQDN")

	// Delete the Credential Connection
	err = runner.client.DeleteCredentialConnection(runner.credentialConnectionID)
	if err != nil {
		runner.logger.Error("Error deleting credential connection", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted Credential Connection")

	// Delete the FQDN Connection
	err = runner.client.DeleteFQDNConnection(runner.fqdnConnectionID)
	if err != nil {
		runner.logger.Error("Error deleting FQDN connection", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted FQDN Connection")

	// Delete the Messaging Profile
	err = runner.client.DeleteMessagingProfile(runner.messagingProfileID)
	if err != nil {
		runner.logger.Error("Error deleting messaging profile", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted Messaging Profile")

	// Delete the Outbound Voice Profile
	err = runner.client.DeleteOutboundVoiceProfile(runner.outboundVoiceProfileID)
	if err != nil {
		runner.logger.Error("Error deleting outbound voice profile", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted Outbound Voice Profile")

	// Delete the Billing Group
	err = runner.client.DeleteBillingGroup(runner.billingGroupID)
	if err != nil {
		runner.logger.Error("Error deleting billing group", zap.Error(err))
		os.Exit(1)
	}
	runner.logger.Info("Deleted Billing Group")
}

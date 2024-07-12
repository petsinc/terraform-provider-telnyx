package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	providerConfig = `
provider "telnyx" {
  endpoint = "https://api.telnyx.com/v2"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"telnyx": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestAccTelnyxResources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "telnyx_billing_group" "test" {
  name = "Test Billing Group Terraform"
}
resource "telnyx_outbound_voice_profile" "test" {
  name                    = "Test Outbound Profile Terraform"
  billing_group_id        = telnyx_billing_group.test.id
  traffic_type            = "conversational"
  service_plan            = "global"
  concurrent_call_limit   = 10
  enabled                 = true
  tags                    = ["test-profile"]
  usage_payment_method    = "rate-deck"
  whitelisted_destinations = ["US"]
  max_destination_rate    = 10.0
  daily_spend_limit       = "100.00"
  daily_spend_limit_enabled = true
  call_recording = {
    type = "all"
    caller_phone_numbers = []
    channels = "single"
    format = "wav"
  }
}
resource "telnyx_messaging_profile" "test" {
  name                      = "Test Messaging Profile Terraform"
  enabled                   = true
  webhook_url               = "https://example.com/webhook"
  webhook_failover_url      = "https://example.com/failover"
  webhook_api_version       = "2"
  whitelisted_destinations  = ["US"]
}
resource "telnyx_credential_connection" "test" {
  connection_name                    = "Test Credential Connection Terraform"
  username                           = "hellopatienttest12345terraform"
  password                           = "54321testpatienthelloterraform"
  active                             = true
  anchorsite_override                = "Latency"
  default_on_hold_comfort_noise_enabled = true
  dtmf_type                          = "RFC 2833"
  encode_contact_header_enabled      = false
  onnet_t38_passthrough_enabled      = false
  microsoft_teams_sbc                = false
  webhook_event_url                  = "https://www.example.com/hooks"
  webhook_event_failover_url         = "https://failover.example.com/hooks"
  webhook_api_version                = "1"
  webhook_timeout_secs               = 25
  rtcp_settings = {
    port = "rtp+1"
    capture_enabled = false
    report_frequency_secs = 5
  }
  inbound = {
    ani_number_format = "E.164-national"
    dnis_number_format = "e164"
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    default_routing_method = "sequential"
    channel_limit = 10
    generate_ringback_tone = true
    isup_headers_enabled = true
    prack_enabled = true
    privacy_zone_enabled = true
    sip_compact_headers_enabled = true
    sip_region = "US"
    sip_subdomain = "uniqueexample.sip.telnyx.com"
    sip_subdomain_receive_settings = "only_my_connections"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
    shaken_stir_enabled = true
  }
  outbound = {
    ani_override = "+12345678901"
    ani_override_type = "always"
    call_parking_enabled = true
    channel_limit = 10
    generate_ringback_tone = true
    instant_ringback_enabled = false
    ip_authentication_method = "token"
    ip_authentication_token = "aBcD1234aBcD1234"
    localization = "US"
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
    t38_reinvite_source = "customer"
    encrypted_media = "SRTP"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
  }
}

resource "telnyx_fqdn_connection" "test" {
  connection_name                    = "Test FQDN Connection Terraform"
  username                           = "fqdnhellopatientest"
  password                           = "fqdnhellopatientestlmao"
  active                             = true
  transport_protocol                 = "UDP"
  encrypted_media                    = null
  anchorsite_override                = "Latency"
  default_on_hold_comfort_noise_enabled = true
  dtmf_type                          = "RFC 2833"
  encode_contact_header_enabled      = false
  onnet_t38_passthrough_enabled      = false
  microsoft_teams_sbc                = false
  webhook_event_url                  = "https://www.example.com/hooks"
  webhook_event_failover_url         = "https://failover.example.com/hooks"
  webhook_api_version                = "1"
  webhook_timeout_secs               = 25
  rtcp_settings = {
    port = "rtp+1"
    capture_enabled = false
    report_frequency_secs = 5
  }
  inbound = {
    ani_number_format = "E.164-national"
    dnis_number_format = "e164"
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    default_routing_method = "sequential"
    channel_limit = 10
    generate_ringback_tone = true
    isup_headers_enabled = true
    prack_enabled = true
    privacy_zone_enabled = true
    sip_compact_headers_enabled = true
    sip_region = "US"
    sip_subdomain = "test.fqdn.connection.uniqueexample.sip.telnyx.com"
    sip_subdomain_receive_settings = "only_my_connections"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
    shaken_stir_enabled = true
  }
  outbound = {
    ani_override = "+12345678901"
    ani_override_type = "always"
    call_parking_enabled = true
    channel_limit = 10
    generate_ringback_tone = true
    instant_ringback_enabled = false
    ip_authentication_method = "token"
    ip_authentication_token = "BBcD1234aBcD1234"
    localization = "US"
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
    t38_reinvite_source = "customer"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
  }
}

resource "telnyx_fqdn" "test" {
  connection_id  = telnyx_fqdn_connection.test.id
  fqdn           = "terraform.test.sip.livekit.cloud"
  dns_record_type = "a"
  port           = 5060
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Test Billing Group Terraform"),
					resource.TestCheckResourceAttr("telnyx_outbound_voice_profile.test", "name", "Test Outbound Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_messaging_profile.test", "name", "Test Messaging Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "connection_name", "Test Credential Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "username", "hellopatienttest12345terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "webhook_event_url", "https://www.example.com/hooks"),
					resource.TestCheckResourceAttr("telnyx_fqdn_connection.test", "connection_name", "Test FQDN Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "fqdn", "terraform.test.sip.livekit.cloud"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "dns_record_type", "a"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "port", "5060"),
				),
			},
			{
				ResourceName:      "telnyx_billing_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
resource "telnyx_billing_group" "test" {
  name = "Updated Billing Group Terraform"
}
resource "telnyx_outbound_voice_profile" "test" {
  name                    = "Updated Test Outbound Profile Terraform"
  billing_group_id        = telnyx_billing_group.test.id
  traffic_type            = "conversational"
  service_plan            = "global"
  concurrent_call_limit   = 10
  enabled                 = true
  tags                    = ["test-profile"]
  usage_payment_method    = "rate-deck"
  whitelisted_destinations = ["US"]
  max_destination_rate    = 10.0
  daily_spend_limit       = "100.00"
  daily_spend_limit_enabled = true
  call_recording = {
    type = "all"
    caller_phone_numbers = []
    channels = "single"
    format = "wav"
  }
}
resource "telnyx_messaging_profile" "test" {
  name                      = "Updated Test Messaging Profile Terraform"
  enabled                   = true
  webhook_url               = "https://example.com/webhook"
  webhook_failover_url      = "https://example.com/failover"
  webhook_api_version       = "2"
  whitelisted_destinations  = ["US"]
}
resource "telnyx_credential_connection" "test" {
  connection_name                    = "Updated Test Credential Connection Terraform"
  username                           = "updatedtest12345terraform"
  password                           = "updatedpassword54321terraform"
  active                             = true
  anchorsite_override                = "Latency"
  default_on_hold_comfort_noise_enabled = true
  dtmf_type                          = "RFC 2833"
  encode_contact_header_enabled      = false
  onnet_t38_passthrough_enabled      = false
  microsoft_teams_sbc                = false
  webhook_event_url                  = "https://www.example.com/hooks"
  webhook_event_failover_url         = "https://failover.example.com/failover"
  webhook_api_version                = "1"
  webhook_timeout_secs               = 25
  rtcp_settings = {
    port = "rtp+1"
    capture_enabled = false
    report_frequency_secs = 5
  }
  inbound = {
    ani_number_format = "E.164-national"
    dnis_number_format = "e164"
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    default_routing_method = "sequential"
    channel_limit = 10
    generate_ringback_tone = true
    isup_headers_enabled = true
    prack_enabled = true
    privacy_zone_enabled = true
    sip_compact_headers_enabled = true
    sip_region = "US"
    sip_subdomain = "uniqueexample.sip.telnyx.com"
    sip_subdomain_receive_settings = "only_my_connections"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
    shaken_stir_enabled = true
  }
  outbound = {
    ani_override = "+12345678901"
    ani_override_type = "always"
    call_parking_enabled = true
    channel_limit = 10
    generate_ringback_tone = true
    instant_ringback_enabled = false
    ip_authentication_method = "token"
    ip_authentication_token = "aBcD1234aBcD1234"
    localization = "US"
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
    t38_reinvite_source = "customer"
    encrypted_media = "SRTP"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
  }
}

resource "telnyx_fqdn_connection" "test" {
  connection_name                    = "Updated Test FQDN Connection Terraform"
  username                           = "fqdnhellopatientest"
  password                           = "fqdnhellopatientestlmao"
  active                             = true
  transport_protocol                 = "UDP"
  encrypted_media                    = null
  anchorsite_override                = "Latency"
  default_on_hold_comfort_noise_enabled = true
  dtmf_type                          = "RFC 2833"
  encode_contact_header_enabled      = false
  onnet_t38_passthrough_enabled      = false
  microsoft_teams_sbc                = false
  webhook_event_url                  = "https://www.example.com/hooks"
  webhook_event_failover_url         = "https://failover.example.com/hooks"
  webhook_api_version                = "1"
  webhook_timeout_secs               = 25
  rtcp_settings = {
    port = "rtp+1"
    capture_enabled = false
    report_frequency_secs = 5
  }
  inbound = {
    ani_number_format = "E.164-national"
    dnis_number_format = "e164"
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    default_routing_method = "sequential"
    channel_limit = 10
    generate_ringback_tone = true
    isup_headers_enabled = true
    prack_enabled = true
    privacy_zone_enabled = true
    sip_compact_headers_enabled = true
    sip_region = "US"
    sip_subdomain = "test.fqdn.connection.uniqueexample.sip.telnyx.com"
    sip_subdomain_receive_settings = "only_my_connections"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
    shaken_stir_enabled = true
  }
  outbound = {
    ani_override = "+12345678901"
    ani_override_type = "always"
    call_parking_enabled = true
    channel_limit = 10
    generate_ringback_tone = true
    instant_ringback_enabled = false
    ip_authentication_method = "token"
    ip_authentication_token = "BBcD1234aBcD1234"
    localization = "US"
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
    t38_reinvite_source = "customer"
    timeout_1xx_secs = 3
    timeout_2xx_secs = 90
  }
}

resource "telnyx_fqdn" "test" {
  connection_id  = telnyx_fqdn_connection.test.id
  fqdn           = "updated.terraform.test.sip.livekit.cloud"
  dns_record_type = "a"
  port           = 5060
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Updated Billing Group Terraform"),
					resource.TestCheckResourceAttr("telnyx_outbound_voice_profile.test", "name", "Updated Test Outbound Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_messaging_profile.test", "name", "Updated Test Messaging Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "connection_name", "Updated Test Credential Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "username", "updatedtest12345terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "webhook_event_url", "https://www.example.com/hooks"),
					resource.TestCheckResourceAttr("telnyx_fqdn_connection.test", "connection_name", "Updated Test FQDN Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "fqdn", "updated.terraform.test.sip.livekit.cloud"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "dns_record_type", "a"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "port", "5060"),
				),
			},
			{
				ResourceName:      "telnyx_billing_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

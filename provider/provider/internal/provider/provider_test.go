package provider

import (
	"flag"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	providerConfig = `
provider "telnyx" {
}
`
)

var (
	includeNumberOrder bool

	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"telnyx": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func init() {
	// Define flags
	flag.BoolVar(&includeNumberOrder, "include-number-order", false, "Include number order test")
}

func TestMain(m *testing.M) {
	// Parse the flags for testing
	flag.Parse()
	// Run the tests
	m.Run()
}

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
  name             = "Test Outbound Voice Profile Terraform"
  billing_group_id = telnyx_billing_group.test.id
  tags             = ["test-profile"]
}

resource "telnyx_messaging_profile" "test" {
  name                = "Test Messaging Profile Terraform"
  enabled             = true
  webhook_url         = "https://example.com/webhook"
  webhook_api_version = "2"
}

resource "telnyx_credential_connection" "test" {
  connection_name            = "Test Credential Connection Terraform"
  username                   = "test12345terraform"
  password                   = "test12345terraform"
  webhook_event_url          = "https://www.example.com/hooks"
  webhook_event_failover_url = "https://failover.example.com/hooks"
  webhook_api_version        = "2"
  inbound = {
    codecs                     = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    sip_subdomain_receive_settings = "from_anyone"
  }
  outbound = {
    ani_override              = "+12345678901"
    ani_override_type         = "always"
    call_parking_enabled      = true
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
  }
}

resource "telnyx_fqdn_connection" "test" {
  connection_name = "Test FQDN Connection Terraform"
  username       = "test12345terraformlmao"
  password        = "test12345terraformlmao"
  inbound = {
    sip_subdomain              = "terraform.test.fqdn.connection.uniqueexample.sip.telnyx.com"
    codecs                     = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    sip_subdomain_receive_settings = "from_anyone"
  }
  outbound = {
    ani_override              = "+12345678901"
    ani_override_type         = "always"
    ip_authentication_token   = "BBcD1234aBcD12345"
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
  }
}

resource "telnyx_fqdn" "test" {
  connection_id   = telnyx_fqdn_connection.test.id
  fqdn            = "terraform.test.sip.livekit.cloud"
  dns_record_type = "a"
  port            = 5060
}

resource "telnyx_texml_application" "test" {
  friendly_name    = "Test TeXML Application Terraform"
  voice_url        = "https://example.com/voice"
  voice_fallback_url = "https://example.com/failover"
  voice_method     = "post"
  inbound = {
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    sip_subdomain                  = "lmao.terraform.test.provider.lol"
    sip_subdomain_receive_settings = "from_anyone"
  }
  outbound = {
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
  }
}

resource "telnyx_phone_number_lookup" "test" {
  starts_with  = "312"
  country_code = "US"
  limit        = 1
  features     = ["sms", "voice"]
}
` + getOptionalNumberOrderConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Test Billing Group Terraform"),
					resource.TestCheckResourceAttr("telnyx_outbound_voice_profile.test", "name", "Test Outbound Voice Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_messaging_profile.test", "name", "Test Messaging Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "connection_name", "Test Credential Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_fqdn_connection.test", "connection_name", "Test FQDN Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "fqdn", "terraform.test.sip.livekit.cloud"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "dns_record_type", "a"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "port", "5060"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "friendly_name", "Test TeXML Application Terraform"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "voice_url", "https://example.com/voice"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "voice_fallback_url", "https://example.com/failover"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "voice_method", "post"),
				),
			},
			{
				ResourceName:      "telnyx_billing_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
locals {
  api_number = telnyx_phone_number_lookup.test.phone_numbers[0].phone_number
}

resource "telnyx_billing_group" "test" {
  name = "Updated Billing Group Terraform"
}

resource "telnyx_outbound_voice_profile" "test" {
  name             = "Updated Test Outbound Voice Profile Terraform"
  billing_group_id = telnyx_billing_group.test.id
  tags             = ["test-profile"]
}

resource "telnyx_messaging_profile" "test" {
  name                      = "Updated Test Messaging Profile Terraform"
  enabled                   = true
  webhook_url               = "https://example.com/webhook"
  webhook_api_version       = "2"
}

resource "telnyx_credential_connection" "test" {
  connection_name                 = "Updated Test Credential Connection Terraform"
  username                        = "test12345terraform"
  password                        = "test12345terraform"
  webhook_event_url               = "https://www.example.com/hooks"
  webhook_event_failover_url      = "https://failover.example.com/hooks"
  webhook_api_version             = "2"
  inbound = {
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    sip_subdomain_receive_settings = "from_anyone"
  }
  outbound = {
    ani_override              = "+12345678901"
    ani_override_type         = "always"
    call_parking_enabled      = true
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
  }
}

resource "telnyx_fqdn_connection" "test" {
  connection_name = "Updated Test FQDN Connection Terraform"
  username       = "test12345terraformlmao"
  password        = "test12345terraformlmao"
  inbound = {
    sip_subdomain              = "terraform.test.fqdn.connection.uniqueexample.sip.telnyx.com"
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    sip_subdomain_receive_settings = "from_anyone"
  }
  outbound = {
    ip_authentication_token   = "BBcD1234aBcD12345"
    ani_override              = "+12345678901"
    ani_override_type         = "always"
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
  }
}

resource "telnyx_fqdn" "test" {
  connection_id  = telnyx_fqdn_connection.test.id
  fqdn           = "updated.terraform.test.sip.livekit.cloud"
  dns_record_type = "a"
  port           = 5060
}

resource "telnyx_texml_application" "test" {
  friendly_name            = "Updated Test TeXML Application Terraform"
  voice_url                = "https://example.com/voice"
  voice_fallback_url       = "https://example.com/failover"
  voice_method             = "post"
  inbound = {
    codecs = ["G722", "G711U", "G711A", "G729", "OPUS", "H.264"]
    sip_subdomain                  = "lmao.terraform.test.provider.lol"
    sip_subdomain_receive_settings = "from_anyone"
  }
  outbound = {
    outbound_voice_profile_id = telnyx_outbound_voice_profile.test.id
  }
}

resource "telnyx_phone_number_lookup" "test" {
  starts_with  = "312"
  country_code = "US"
  limit        = 1
  features     = ["sms", "voice"]
}
` + getOptionalNumberOrderConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Updated Billing Group Terraform"),
					resource.TestCheckResourceAttr("telnyx_outbound_voice_profile.test", "name", "Updated Test Outbound Voice Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_messaging_profile.test", "name", "Updated Test Messaging Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_credential_connection.test", "connection_name", "Updated Test Credential Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_fqdn_connection.test", "connection_name", "Updated Test FQDN Connection Terraform"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "fqdn", "updated.terraform.test.sip.livekit.cloud"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "dns_record_type", "a"),
					resource.TestCheckResourceAttr("telnyx_fqdn.test", "port", "5060"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "friendly_name", "Updated Test TeXML Application Terraform"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "voice_url", "https://example.com/voice"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "voice_fallback_url", "https://example.com/failover"),
					resource.TestCheckResourceAttr("telnyx_texml_application.test", "voice_method", "post"),
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

func getOptionalNumberOrderConfig() string {
	if includeNumberOrder {
		return `
resource "telnyx_number_order" "this" {
  connection_id      = telnyx_texml_application.test.id
  billing_group_id   = telnyx_billing_group.test.id
  customer_reference = "terraform-test-api-number"
  phone_numbers = [
    {
      phone_number = telnyx_phone_number_lookup.test.phone_numbers[0].phone_number
    }
  ]
}
`
	}
	return ""
}

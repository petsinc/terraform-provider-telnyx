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

func TestAccTelnyxBillingGroupResource(t *testing.T) {
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Test Billing Group Terraform"),
					resource.TestCheckResourceAttr("telnyx_outbound_voice_profile.test", "name", "Test Outbound Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_messaging_profile.test", "name", "Test Messaging Profile Terraform"),
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Updated Billing Group Terraform"),
					resource.TestCheckResourceAttr("telnyx_outbound_voice_profile.test", "name", "Updated Test Outbound Profile Terraform"),
					resource.TestCheckResourceAttr("telnyx_messaging_profile.test", "name", "Updated Test Messaging Profile Terraform"),
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

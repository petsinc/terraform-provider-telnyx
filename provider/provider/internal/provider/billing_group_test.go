package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Telnyx client is properly configured.
	providerConfig = `
provider "telnyx" {
  endpoint = "https://api.telnyx.com/v2"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Test Billing Group Terraform"),
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("telnyx_billing_group.test", "name", "Updated Billing Group Terraform"),
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

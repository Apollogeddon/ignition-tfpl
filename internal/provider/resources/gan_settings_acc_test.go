package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGanSettingsResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGanSettingsResourceConfig("ApprovedOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_gan_settings.test", "name", "gateway-network-settings"),
					resource.TestCheckResourceAttr("ignition_gan_settings.test", "security_policy", "ApprovedOnly"),
				),
			},
			// Update and Read testing
			{
				Config: testAccGanSettingsResourceConfig("Unrestricted"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_gan_settings.test", "security_policy", "Unrestricted"),
				),
			},
		},
	})
}

func testAccGanSettingsResourceConfig(securityPolicy string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_gan_settings" "test" {
  security_policy = %[1]q
}
`, securityPolicy)
}

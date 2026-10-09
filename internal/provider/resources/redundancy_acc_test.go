package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpr/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRedundancyResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRedundancyResourceConfig("Independent", "Automatic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_redundancy.test", "name", "gateway-redundancy"),
					resource.TestCheckResourceAttr("ignition_redundancy.test", "role", "Independent"),
					resource.TestCheckResourceAttr("ignition_redundancy.test", "recovery_mode", "Automatic"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRedundancyResourceConfig("Independent", "Manual"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_redundancy.test", "recovery_mode", "Manual"),
				),
			},
		},
	})
}

func testAccRedundancyResourceConfig(role, recoveryMode string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_redundancy" "test" {
  role          = %[1]q
  recovery_mode = %[2]q
}
`, role, recoveryMode)
}

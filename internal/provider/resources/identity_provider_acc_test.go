package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityProviderResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIdentityProviderResourceConfig(rName, 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_identity_provider.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_identity_provider.test", "type", "internal"),
					resource.TestCheckResourceAttr("ignition_identity_provider.test", "user_source", "default"),
					resource.TestCheckResourceAttr("ignition_identity_provider.test", "session_inactivity_timeout", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_identity_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIdentityProviderResourceConfig(rName, 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_identity_provider.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_identity_provider.test", "session_inactivity_timeout", "60"),
				),
			},
			{
				ResourceName:      "ignition_identity_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityProviderResourceConfig(name string, sessionInactivityTimeout int) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_identity_provider" "test" {
  name                        = %[1]q
  type                        = "internal"
  user_source                 = "default"
  session_inactivity_timeout  = %[2]d
}
`, name, sessionInactivityTimeout)
}

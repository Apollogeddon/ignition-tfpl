package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGanOutgoingResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGanOutgoingResourceConfig(rName, "test-gan-host.invalid", 8060),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_gan_outgoing.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_gan_outgoing.test", "host", "test-gan-host.invalid"),
					resource.TestCheckResourceAttr("ignition_gan_outgoing.test", "port", "8060"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_gan_outgoing.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccGanOutgoingResourceConfig(rName, "test-gan-host.invalid", 8070),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_gan_outgoing.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_gan_outgoing.test", "port", "8070"),
				),
			},
			{
				ResourceName:      "ignition_gan_outgoing.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGanOutgoingResourceConfig(name, host string, port int) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_gan_outgoing" "test" {
  name = %[1]q
  host = %[2]q
  port = %[3]d
}
`, name, host, port)
}

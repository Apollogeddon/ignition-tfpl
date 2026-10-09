package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tofu/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpcUaConnectionResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOpcUaConnectionResourceConfig(rName, "opc.tcp://test-opcua-host.invalid:4840/discovery"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.test", "security_policy", "None"),
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.test", "security_mode", "None"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_opc_ua_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOpcUaConnectionResourceConfig(rName, "opc.tcp://updated-opcua-host.invalid:4840/discovery"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.test", "discovery_url", "opc.tcp://updated-opcua-host.invalid:4840/discovery"),
				),
			},
			{
				ResourceName:      "ignition_opc_ua_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOpcUaConnectionResourceConfig(name, discoveryURL string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_opc_ua_connection" "test" {
  name            = %[1]q
  discovery_url   = %[2]q
  endpoint_url    = "opc.tcp://test-opcua-host.invalid:4840"
  security_policy = "None"
  security_mode   = "None"
}
`, name, discoveryURL)
}

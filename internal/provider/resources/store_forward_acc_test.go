package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStoreForwardResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStoreForwardResourceConfig(rName, "ALL", 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_store_forward.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_store_forward.test", "forwarding_policy", "ALL"),
					resource.TestCheckResourceAttr("ignition_store_forward.test", "batch_size", "100"),
				),
			},
			{
				ResourceName:      "ignition_store_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStoreForwardResourceConfig(name, policy string, batchSize int) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_store_forward" "test" {
  name              = %[1]q
  forwarding_policy = %[2]q
  batch_size        = %[3]d
  time_threshold_ms = 1000
}
`, name, policy, batchSize)
}

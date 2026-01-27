package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserSourceResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserSourceResourceConfig(rName, "INTERNAL", "Test User Source"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_user_source.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_user_source.test", "type", "INTERNAL"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "description", "Test User Source"),
				),
			},
			{
				Config: testAccUserSourceResourceConfig(rName, "INTERNAL", "Updated User Source"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_user_source.test", "description", "Updated User Source"),
				),
			},
			{
				ResourceName:      "ignition_user_source.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccUserSourceResourceConfig(name, typeVal, desc string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_user_source" "test" {
  name          = %[1]q
  type          = %[2]q
  description   = %[3]q
}
`, name, typeVal, desc)
}

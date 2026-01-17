package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagProviderResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagProviderResourceConfig(rName, "STANDARD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "type", "STANDARD"),
				),
			},
			{
				ResourceName:      "ignition_tag_provider.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTagProviderResourceConfig(name, typeVal string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_tag_provider" "test" {
  name = %[1]q
  type = %[2]q
}
`, name, typeVal)
}

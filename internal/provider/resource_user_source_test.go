package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserSourceResource(t *testing.T) {
	// Only run acceptance tests if IGNITION_HOST and IGNITION_TOKEN are set
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserSourceResourceConfig("test_users_4", "INTERNAL", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_user_source.test", "name", "test_users_4"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "type", "INTERNAL"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "schedule_restricted", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_user_source.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccUserSourceResourceConfig("test_users_4", "INTERNAL", "Updated Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_user_source.test", "name", "test_users_4"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "type", "INTERNAL"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "description", "Updated Description"),
				),
			},
		},
	})
}

func testAccUserSourceResourceConfig(name, typeVal, desc string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_user_source" "test" {
  name        = %[1]q
  type        = %[2]q
  description = %[3]q
}
`, name, typeVal, desc)
}
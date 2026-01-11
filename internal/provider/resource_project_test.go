package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	// Only run acceptance tests if IGNITION_HOST and IGNITION_TOKEN are set
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig("test_project_4", "A test project", "Test Project"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_project.test", "name", "test_project_4"),
					resource.TestCheckResourceAttr("ignition_project.test", "description", "A test project"),
					resource.TestCheckResourceAttr("ignition_project.test", "title", "Test Project"),
					resource.TestCheckResourceAttr("ignition_project.test", "enabled", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig("test_project_4", "Updated description", "Updated Title"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_project.test", "name", "test_project_4"),
					resource.TestCheckResourceAttr("ignition_project.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("ignition_project.test", "title", "Updated Title"),
				),
			},
		},
	})
}

func testAccProjectResourceConfig(name, desc, title string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_project" "test" {
  name        = %[1]q
  description = %[2]q
  title       = %[3]q
  enabled     = true
}
`, name, desc, title)
}
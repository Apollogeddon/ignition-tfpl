package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAuditProfileResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	dbName := "db_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAuditProfileResourceConfig(rName, "database", dbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "type", "database"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "name", dbName),
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "retention_days", "90"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_audit_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAuditProfileResourceConfig(rName, "database", dbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "name", rName),
				),
			},
			{
				ResourceName:      "ignition_audit_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAuditProfileResourceConfig(name, profileType, dbName string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_database_connection" "test" {
  name        = %[3]q
  type        = "MariaDB"
  translator  = "MYSQL"
  connect_url = "jdbc:mariadb://localhost:3306/test"
}

resource "ignition_audit_profile" "test" {
  name        = %[1]q
  type        = %[2]q
  database    = ignition_database_connection.test.name
  retention_days = 90
}
`, name, profileType, dbName)
}

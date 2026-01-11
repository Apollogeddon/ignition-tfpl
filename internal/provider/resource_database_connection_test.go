package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseConnectionResource(t *testing.T) {
	// Only run acceptance tests if IGNITION_HOST and IGNITION_TOKEN are set
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseConnectionResourceConfig("test_db_4", "MariaDB", "MYSQL", "jdbc:mariadb://localhost:3306/testdb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_database_connection.test", "name", "test_db_4"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "type", "MariaDB"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "translator", "MYSQL"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "connect_url", "jdbc:mariadb://localhost:3306/testdb"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_database_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDatabaseConnectionResourceConfig("test_db_4", "PostgreSQL", "POSTGRESQL", "jdbc:postgresql://localhost:5432/testdb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_database_connection.test", "name", "test_db_4"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "type", "PostgreSQL"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "translator", "POSTGRESQL"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "connect_url", "jdbc:postgresql://localhost:5432/testdb"),
				),
			},
		},
	})
}

func testAccDatabaseConnectionResourceConfig(name, dbType, translator, url string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_database_connection" "test" {
  name        = %[1]q
  type        = %[2]q
  translator  = %[3]q
  connect_url = %[4]q
}
`, name, dbType, translator, url)
}
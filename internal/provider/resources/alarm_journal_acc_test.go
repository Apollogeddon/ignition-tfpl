package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tofu/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlarmJournalResource(t *testing.T) {
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
				Config: testAccAlarmJournalResourceConfig(rName, dbName, "Low"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_journal.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_alarm_journal.test", "type", "DATASOURCE"),
					resource.TestCheckResourceAttr("ignition_alarm_journal.test", "datasource", dbName),
					resource.TestCheckResourceAttr("ignition_alarm_journal.test", "min_priority", "Low"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_alarm_journal.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAlarmJournalResourceConfig(rName, dbName, "Medium"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_journal.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_alarm_journal.test", "min_priority", "Medium"),
				),
			},
			{
				ResourceName:      "ignition_alarm_journal.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlarmJournalResourceConfig(name, dbName, minPriority string) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_database_connection" "test" {
  name        = %[2]q
  type        = "MariaDB"
  translator  = "MYSQL"
  connect_url = "jdbc:mariadb://localhost:3306/test"
}

resource "ignition_alarm_journal" "test" {
  name         = %[1]q
  type         = "DATASOURCE"
  datasource   = ignition_database_connection.test.name
  min_priority = %[3]q
}
`, name, dbName, minPriority)
}

package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlarmNotificationProfileResource(t *testing.T) {
	// Only run acceptance tests if IGNITION_HOST and IGNITION_TOKEN are set
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAlarmNotificationProfileResourceConfig(rName, "EmailNotificationProfileType", "smtp.example.com", 25),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "type", "EmailNotificationProfileType"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "email_config.hostname", "smtp.example.com"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "email_config.port", "25"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ignition_alarm_notification_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Password/secrets are not returned by API, so we skip verifying them
				ImportStateVerifyIgnore: []string{"email_config.password"},
			},
			// Update and Read testing
			{
				Config: testAccAlarmNotificationProfileResourceConfig(rName, "EmailNotificationProfileType", "smtp.updated.com", 587),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "type", "EmailNotificationProfileType"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "email_config.hostname", "smtp.updated.com"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.test", "email_config.port", "587"),
				),
			},
		},
	})
}

func testAccAlarmNotificationProfileResourceConfig(name, typeVal, hostname string, port int) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_alarm_notification_profile" "test" {
  name = %[1]q
  type = %[2]q

  email_config {
    use_smtp_profile = false
    hostname         = %[3]q
    port             = %[4]d
    ssl_enabled      = false
    username         = "user"
    password         = "pass"
  }
}
`, name, typeVal, hostname, port)
}

package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSMTPProfileResource(t *testing.T) {
	if os.Getenv("IGNITION_HOST") == "" || os.Getenv("IGNITION_TOKEN") == "" {
		t.Skip("Skipping acceptance test: IGNITION_HOST and/or IGNITION_TOKEN not set")
	}

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSMTPProfileResourceConfig(rName, "smtp.example.com", 25),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_smtp_profile.test", "name", rName),
					resource.TestCheckResourceAttr("ignition_smtp_profile.test", "hostname", "smtp.example.com"),
					resource.TestCheckResourceAttr("ignition_smtp_profile.test", "port", "25"),
				),
			},
			{
				ResourceName:      "ignition_smtp_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccSMTPProfileResourceConfig(name, hostname string, port int) string {
	return fmt.Sprintf(`
provider "ignition" {}

resource "ignition_smtp_profile" "test" {
  name     = %[1]q
  hostname = %[2]q
  port     = %[3]d
}
`, name, hostname, port)
}

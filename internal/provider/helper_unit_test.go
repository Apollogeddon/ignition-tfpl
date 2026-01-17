package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitHelper_ErrorPaths(t *testing.T) {
	mockClient := &client.MockClient{
		CreateSMTPProfileFunc: func(ctx context.Context, item client.ResourceResponse[client.SMTPProfileConfig]) (*client.ResourceResponse[client.SMTPProfileConfig], error) {
			return nil, fmt.Errorf("simulated creation failure")
		},
	}

	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"ignition": providerserver.NewProtocol6WithError(&IgnitionProvider{
			version: "test",
			client:  mockClient,
		}),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_smtp_profile" "fail" {
						name     = "fail"
						hostname = "fail"
					}
				`,
				ExpectError: regexp.MustCompile(`simulated creation failure`),
			},
		},
	})
}

package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitProjectDataSource(t *testing.T) {
	mockClient := &client.MockClient{
		GetProjectFunc: func(ctx context.Context, name string) (*client.Project, error) {
			return &client.Project{
				Name:        "test-project",
				Description: "A test project",
				Enabled:     true,
				Title:       "Test Project",
			}, nil
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
					data "ignition_project" "unit" {
						name = "test-project"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ignition_project.unit", "name", "test-project"),
					resource.TestCheckResourceAttr("data.ignition_project.unit", "description", "A test project"),
					resource.TestCheckResourceAttr("data.ignition_project.unit", "enabled", "true"),
				),
			},
		},
	})
}

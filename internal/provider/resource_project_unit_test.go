package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitProjectResource(t *testing.T) {
	mockProject := &client.Project{
		Name:        "test-project",
		Description: "Initial Description",
		Enabled:     true,
	}

	mockClient := &client.MockClient{
		CreateProjectFunc: func(ctx context.Context, p client.Project) (*client.Project, error) {
			mockProject = &p
			return mockProject, nil
		},
		GetProjectFunc: func(ctx context.Context, name string) (*client.Project, error) {
			return mockProject, nil
		},
		UpdateProjectFunc: func(ctx context.Context, p client.Project) (*client.Project, error) {
			mockProject = &p
			return mockProject, nil
		},
		DeleteProjectFunc: func(ctx context.Context, name string) error {
			return nil
		},
	}

	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"ignition": providerserver.NewProtocol6WithError(&IgnitionProvider{
			version: "test",
			client:  mockClient,
		}),
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_project" "test" {
						name        = "test-project"
						description = "Initial Description"
						enabled     = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_project.test", "name", "test-project"),
					resource.TestCheckResourceAttr("ignition_project.test", "description", "Initial Description"),
					resource.TestCheckResourceAttr("ignition_project.test", "enabled", "true"),
				),
			},
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_project" "test" {
						name        = "test-project"
						description = "Updated Description"
						enabled     = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_project.test", "description", "Updated Description"),
					resource.TestCheckResourceAttr("ignition_project.test", "enabled", "false"),
				),
			},
		},
	})
}

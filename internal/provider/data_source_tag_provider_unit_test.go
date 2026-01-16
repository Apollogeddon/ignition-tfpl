package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitTagProviderDataSource(t *testing.T) {
	mockClient := &client.MockClient{
		GetTagProviderFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.TagProviderConfig], error) {
			return &client.ResourceResponse[client.TagProviderConfig]{
				Name: name,
				Config: client.TagProviderConfig{
					Type:        "standard",
					Description: "A mock tag provider",
				},
			}, nil
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
					data "ignition_tag_provider" "test" {
						name = "mock-tags"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ignition_tag_provider.test", "name", "mock-tags"),
					resource.TestCheckResourceAttr("data.ignition_tag_provider.test", "type", "standard"),
					resource.TestCheckResourceAttr("data.ignition_tag_provider.test", "description", "A mock tag provider"),
				),
			},
		},
	})
}

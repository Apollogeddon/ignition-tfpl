package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitUserSourceResource(t *testing.T) {
	mockClient := &client.MockClient{
		CreateUserSourceFunc: func(ctx context.Context, us client.ResourceResponse[client.UserSourceConfig]) (*client.ResourceResponse[client.UserSourceConfig], error) {
			us.Signature = "sig-123"
			return &us, nil
		},
		GetUserSourceFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.UserSourceConfig], error) {
			return &client.ResourceResponse[client.UserSourceConfig]{
				Name:      name,
				Enabled:   boolPtr(true),
				Signature: "sig-123",
				Config: client.UserSourceConfig{
					Profile: client.UserSourceProfile{
						Type:         "INTERNAL",
						FailoverMode: "Soft",
					},
				},
			}, nil
		},
		UpdateUserSourceFunc: func(ctx context.Context, us client.ResourceResponse[client.UserSourceConfig]) (*client.ResourceResponse[client.UserSourceConfig], error) {
			us.Signature = "sig-456"
			return &us, nil
		},
		DeleteUserSourceFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_user_source" "test" {
						name          = "TestUserSource"
						type          = "INTERNAL"
						failover_mode = "Soft"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_user_source.test", "name", "TestUserSource"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "type", "INTERNAL"),
					resource.TestCheckResourceAttr("ignition_user_source.test", "failover_mode", "Soft"),
				),
			},
		},
	})
}

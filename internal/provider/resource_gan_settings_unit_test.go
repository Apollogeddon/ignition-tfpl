package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitGanSettingsResource(t *testing.T) {
	mockClient := &client.MockClient{
		GetGanGeneralSettingsFunc: func(ctx context.Context) (*client.ResourceResponse[client.GanGeneralSettingsConfig], error) {
			return &client.ResourceResponse[client.GanGeneralSettingsConfig]{
				Name:    "gateway-network-settings",
				Enabled: boolPtr(true),
				Config: client.GanGeneralSettingsConfig{
					RequireSSL:                  true,
					RequireTwoWayAuth:           true,
					AllowIncoming:               true,
					SecurityPolicy:              "ApprovedOnly",
					Whitelist:                   "",
					WebsocketSessionIdleTimeout: 30000,
					TempFilesMaxAgeHours:        24,
				},
			}, nil
		},
		UpdateGanGeneralSettingsFunc: func(ctx context.Context, item client.ResourceResponse[client.GanGeneralSettingsConfig]) (*client.ResourceResponse[client.GanGeneralSettingsConfig], error) {
			return &item, nil
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
					resource "ignition_gan_settings" "unit" {
						require_ssl = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_gan_settings.unit", "require_ssl", "true"),
					resource.TestCheckResourceAttr("ignition_gan_settings.unit", "security_policy", "ApprovedOnly"),
				),
			},
		},
	})
}

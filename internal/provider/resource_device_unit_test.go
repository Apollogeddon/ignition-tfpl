package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitDeviceResource(t *testing.T) {
	currentRate := 1000

	mockClient := &client.MockClient{
		CreateDeviceFunc: func(ctx context.Context, item client.ResourceResponse[client.DeviceConfig]) (*client.ResourceResponse[client.DeviceConfig], error) {
			item.Signature = "sig-123"
			return &item, nil
		},
		GetDeviceFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.DeviceConfig], error) {
			return &client.ResourceResponse[client.DeviceConfig]{
				Name:      name,
				Type:      "ProgrammableSimulatorDevice",
				Enabled:   boolPtr(true),
				Signature: "sig-123",
				Config:    client.DeviceConfig{"baseRate": currentRate},
			}, nil
		},
		UpdateDeviceFunc: func(ctx context.Context, item client.ResourceResponse[client.DeviceConfig]) (*client.ResourceResponse[client.DeviceConfig], error) {
			item.Signature = "sig-456"
			if val, ok := item.Config["baseRate"]; ok {
				if f, ok := val.(float64); ok {
					currentRate = int(f)
				} else if i, ok := val.(int); ok {
					currentRate = i
				}
			}
			return &item, nil
		},
		DeleteDeviceFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_device" "test" {
						name = "SimDevice"
						type = "ProgrammableSimulatorDevice"
						parameters = "{\"baseRate\": 1000}"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_device.test", "name", "SimDevice"),
					resource.TestCheckResourceAttr("ignition_device.test", "type", "ProgrammableSimulatorDevice"),
					resource.TestCheckResourceAttr("ignition_device.test", "parameters", `{"baseRate": 1000}`),
				),
			},
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_device" "test" {
						name = "SimDevice"
						type = "ProgrammableSimulatorDevice"
						parameters = "{\"baseRate\": 2000}"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_device.test", "parameters", `{"baseRate": 2000}`),
				),
			},
		},
	})
}

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitOpcUaConnectionResource_Create(t *testing.T) {
	mockClient := &client.MockClient{
		CreateOpcUaConnectionFunc: func(ctx context.Context, item client.ResourceResponse[client.OpcUaConnectionConfig]) (*client.ResourceResponse[client.OpcUaConnectionConfig], error) {
			if item.Name != "unit-test-opc" {
				return nil, fmt.Errorf("expected name 'unit-test-opc', got '%s'", item.Name)
			}
			if item.Config.Settings.Endpoint.EndpointURL != "opc.tcp://localhost:4840" {
				return nil, fmt.Errorf("unexpected endpoint URL: %s", item.Config.Settings.Endpoint.EndpointURL)
			}
			
			// Simulate successful creation
			item.Signature = "mock-signature-opc"
			item.Config.Profile.Type = "com.inductiveautomation.OpcUaServerType"
			return &item, nil
		},
		GetOpcUaConnectionFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.OpcUaConnectionConfig], error) {
			if name != "unit-test-opc" {
				return nil, fmt.Errorf("not found")
			}
			return &client.ResourceResponse[client.OpcUaConnectionConfig]{
				Name:      "unit-test-opc",
				Enabled:   boolPtr(true),
				Signature: "mock-signature-opc",
				Config: client.OpcUaConnectionConfig{
					Profile: client.OpcUaConnectionProfile{Type: "com.inductiveautomation.OpcUaServerType"},
					Settings: client.OpcUaConnectionSettings{
						Endpoint: client.OpcUaConnectionEndpoint{
							DiscoveryURL:   "opc.tcp://localhost:4840",
							EndpointURL:    "opc.tcp://localhost:4840",
							SecurityPolicy: "None",
							SecurityMode:   "None",
						},
					},
				},
			}, nil
		},
		DeleteOpcUaConnectionFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_opc_ua_connection" "unit" {
						name = "unit-test-opc"
						discovery_url = "opc.tcp://localhost:4840"
						endpoint_url  = "opc.tcp://localhost:4840"
						security_policy = "None"
						security_mode   = "None"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.unit", "name", "unit-test-opc"),
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.unit", "signature", "mock-signature-opc"),
				),
			},
		},
	})
}

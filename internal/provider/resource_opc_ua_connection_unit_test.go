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
			if item.Name != "test-opc" {
				return nil, fmt.Errorf("expected name 'test-opc', got '%s'", item.Name)
			}
			if item.Config.Settings.Endpoint.EndpointURL != "opc.tcp://localhost:53530" {
				return nil, fmt.Errorf("unexpected endpoint URL")
			}
			
			// Simulate successful creation
			item.Signature = "mock-signature-opc"
			item.Config.Profile.Type = "com.inductiveautomation.OpcUaServerType"
			return &item, nil
		},
		GetOpcUaConnectionFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.OpcUaConnectionConfig], error) {
			if name != "test-opc" {
				return nil, fmt.Errorf("not found")
			}
			return &client.ResourceResponse[client.OpcUaConnectionConfig]{
				Name:      "test-opc",
				Enabled:   true,
				Signature: "mock-signature-opc",
				Config: client.OpcUaConnectionConfig{
					Profile: client.OpcUaConnectionProfile{Type: "com.inductiveautomation.OpcUaServerType"},
					Settings: client.OpcUaConnectionSettings{
						Endpoint: client.OpcUaConnectionEndpoint{
							DiscoveryURL:   "opc.tcp://localhost:53530/discovery",
							EndpointURL:    "opc.tcp://localhost:53530",
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_opc_ua_connection" "unit" {
						name            = "test-opc"
						discovery_url   = "opc.tcp://localhost:53530/discovery"
						endpoint_url    = "opc.tcp://localhost:53530"
						security_policy = "None"
						security_mode   = "None"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.unit", "name", "test-opc"),
					resource.TestCheckResourceAttr("ignition_opc_ua_connection.unit", "signature", "mock-signature-opc"),
				),
			},
		},
	})
}

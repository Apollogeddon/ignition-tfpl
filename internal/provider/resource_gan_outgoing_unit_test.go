package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitGanOutgoingResource_Create(t *testing.T) {
	mockClient := &client.MockClient{
		CreateGanOutgoingFunc: func(ctx context.Context, item client.ResourceResponse[client.GanOutgoingConfig]) (*client.ResourceResponse[client.GanOutgoingConfig], error) {
			if item.Config.Host != "192.168.1.100" {
				return nil, fmt.Errorf("expected host '192.168.1.100', got '%s'", item.Config.Host)
			}
			item.Signature = "mock-signature-gan"
			return &item, nil
		},
		GetGanOutgoingFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.GanOutgoingConfig], error) {
			return &client.ResourceResponse[client.GanOutgoingConfig]{
				Name:      name,
				Enabled:   boolPtr(true),
				Signature: "mock-signature-gan",
				Config: client.GanOutgoingConfig{
					Host:                     "192.168.1.100",
					Port:                     8060,
					PingRateMillis:           2000,
					PingTimeoutMillis:        60000,
					PingMaxMissed:            3,
					WsTimeoutMillis:          10000,
					HttpConnectTimeoutMillis: 10000,
					HttpReadTimeoutMillis:    30000,
					SendThreads:              1,
					ReceiveThreads:           1,
				},
			}, nil
		},
		DeleteGanOutgoingFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_gan_outgoing" "unit" {
						name = "test-gan"
						host = "192.168.1.100"
						port = 8060
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_gan_outgoing.unit", "host", "192.168.1.100"),
					resource.TestCheckResourceAttr("ignition_gan_outgoing.unit", "signature", "mock-signature-gan"),
				),
			},
		},
	})
}

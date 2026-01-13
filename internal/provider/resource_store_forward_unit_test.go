package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitStoreForwardResource(t *testing.T) {
	mockClient := &client.MockClient{
		CreateStoreAndForwardFunc: func(ctx context.Context, item client.ResourceResponse[client.StoreAndForwardConfig]) (*client.ResourceResponse[client.StoreAndForwardConfig], error) {
			item.Signature = "sig"
			return &item, nil
		},
		GetStoreAndForwardFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.StoreAndForwardConfig], error) {
			return &client.ResourceResponse[client.StoreAndForwardConfig]{
				Name:      name,
				Signature: "sig",
				Enabled:   true,
				Config: client.StoreAndForwardConfig{
					TimeThresholdMs:    1000,
					ForwardRateMs:      1000,
					ForwardingPolicy:   "ALL",
					ForwardingSchedule: "",
					IsThirdParty:       false,
					DataThreshold:      100,
					BatchSize:          100,
					ScanRateMs:         1000,
					PrimaryMaintenancePolicy: &client.StoreAndForwardMaintenancePolicy{
						Action: "EVICT_OLDEST_DATA",
						Limit: struct {
							LimitType string `json:"limitType"`
							Value     int    `json:"value"`
						}{
							LimitType: "COUNT",
							Value:     1000,
						},
					},
				},
			}, nil
		},
		DeleteStoreAndForwardFunc: func(ctx context.Context, name, sig string) error { return nil },
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
					resource "ignition_store_forward" "unit" {
						name = "sf"
						primary_policy = {
							action = "EVICT_OLDEST_DATA"
							limit_type = "COUNT"
							value = 1000
						}
						# Match defaults to ensure clean plan
						time_threshold_ms = 1000
						forward_rate_ms = 1000
						forwarding_policy = "ALL"
						data_threshold = 100
						batch_size = 100
						scan_rate_ms = 1000
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_store_forward.unit", "name", "sf"),
					resource.TestCheckResourceAttr("ignition_store_forward.unit", "primary_policy.action", "EVICT_OLDEST_DATA"),
				),
			},
		},
	})
}

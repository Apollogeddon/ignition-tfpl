package resources

import (
	"context"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/apollogeddon/ignition-tfpl/internal/provider/base"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitRedundancyResource(t *testing.T) {
	mockClient := &client.MockClient{
		GetRedundancyConfigFunc: func(ctx context.Context) (*client.RedundancyConfig, error) {
			return &client.RedundancyConfig{
				Role:               "Master",
				ActiveHistoryLevel: "Full",
				JoinWaitTime:       10000,
				RecoveryMode:       "Automatic",
			}, nil
		},
		UpdateRedundancyConfigFunc: func(ctx context.Context, config client.RedundancyConfig) error {
			return nil
		},
	}

	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
			ResourceFactory: NewRedundancyResource,
			Client:          mockClient,
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
					resource "ignition_redundancy" "unit" {
						role = "Master"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_redundancy.unit", "role", "Master"),
					resource.TestCheckResourceAttr("ignition_redundancy.unit", "active_history_level", "Full"),
				),
			},
		},
	})
}

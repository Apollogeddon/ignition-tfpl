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

func TestUnitTagProviderResource(t *testing.T) {
	currentDescription := "Test Description"
	currentSignature := "sig-123"

	mockClient := &client.MockClient{
		CreateTagProviderFunc: func(ctx context.Context, tp client.ResourceResponse[client.TagProviderConfig]) (*client.ResourceResponse[client.TagProviderConfig], error) {
			tp.Signature = "sig-123"
			return &tp, nil
		},
		GetTagProviderFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.TagProviderConfig], error) {
			if name != "test-tags" {
				return nil, fmt.Errorf("not found")
			}
			return &client.ResourceResponse[client.TagProviderConfig]{
				Name:      "test-tags",
				Signature: currentSignature,
				Config: client.TagProviderConfig{
					Type:        "standard",
					Description: currentDescription,
				},
			}, nil
		},
		UpdateTagProviderFunc: func(ctx context.Context, tp client.ResourceResponse[client.TagProviderConfig]) (*client.ResourceResponse[client.TagProviderConfig], error) {
			currentDescription = tp.Config.Description
			currentSignature = "sig-456"
			tp.Signature = currentSignature
			return &tp, nil
		},
		DeleteTagProviderFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_tag_provider" "test" {
						name        = "test-tags"
						type        = "standard"
						description = "Test Description"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "name", "test-tags"),
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "type", "standard"),
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "description", "Test Description"),
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "signature", "sig-123"),
				),
			},
			// Update
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_tag_provider" "test" {
						name        = "test-tags"
						type        = "standard"
						description = "Updated Description"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "description", "Updated Description"),
					resource.TestCheckResourceAttr("ignition_tag_provider.test", "signature", "sig-456"),
				),
			},
		},
	})
}

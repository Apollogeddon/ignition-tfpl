package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/apollogeddon/ignition-tfpl/internal/provider/base"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitIdentityProviderResource_Create(t *testing.T) {
	mockClient := &client.MockClient{
		CreateIdentityProviderFunc: func(ctx context.Context, item client.ResourceResponse[client.IdentityProviderConfig]) (*client.ResourceResponse[client.IdentityProviderConfig], error) {
			if item.Name != "unit-test-idp" {
				return nil, fmt.Errorf("expected name 'unit-test-idp', got '%s'", item.Name)
			}

			// Simulate successful creation
			item.Signature = "mock-signature-idp"
			return &item, nil
		},
		GetIdentityProviderFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.IdentityProviderConfig], error) {
			if name != "unit-test-idp" {
				return nil, fmt.Errorf("not found")
			}

			// Return an internal IdP
			internalConfig := client.IdentityProviderInternalConfig{
				UserSource:               "default",
				SessionInactivityTimeout: 30,
				SessionExp:               0,
				RememberMeExp:            0,
			}

			return &client.ResourceResponse[client.IdentityProviderConfig]{
				Name:      "unit-test-idp",
				Enabled:   base.BoolPtr(true),
				Signature: "mock-signature-idp",
				Config: client.IdentityProviderConfig{
					Type:   "internal",
					Config: internalConfig,
				},
			}, nil
		},
		DeleteIdentityProviderFunc: func(ctx context.Context, name, signature string) error {
			return nil
		},
	}

	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
			ResourceFactory: NewIdentityProviderResource,
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
					resource "ignition_identity_provider" "unit" {
						name        = "unit-test-idp"
						type        = "internal"
						user_source = "default"
						session_inactivity_timeout = 30
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_identity_provider.unit", "name", "unit-test-idp"),
					resource.TestCheckResourceAttr("ignition_identity_provider.unit", "type", "internal"),
					resource.TestCheckResourceAttr("ignition_identity_provider.unit", "signature", "mock-signature-idp"),
				),
			},
		},
	})
}

func TestUnitIdentityProviderResource_OIDC(t *testing.T) {
	mockClient := &client.MockClient{
		CreateIdentityProviderFunc: func(ctx context.Context, item client.ResourceResponse[client.IdentityProviderConfig]) (*client.ResourceResponse[client.IdentityProviderConfig], error) {
			if item.Config.Type != "oidc" {
				return nil, fmt.Errorf("expected type 'oidc', got '%s'", item.Config.Type)
			}
			item.Signature = "mock-signature-oidc"
			return &item, nil
		},
		GetIdentityProviderFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.IdentityProviderConfig], error) {
			oidcConfig := client.IdentityProviderOidcConfig{
				ClientId:                   "test-client",
				ProviderId:                 "https://auth.com",
				JsonWebKeysEndpointEnabled: true,
			}
			return &client.ResourceResponse[client.IdentityProviderConfig]{
				Name:      name,
				Enabled:   base.BoolPtr(true),
				Signature: "mock-signature-oidc",
				Config: client.IdentityProviderConfig{
					Type:   "oidc",
					Config: oidcConfig,
				},
			}, nil
		},
		DeleteIdentityProviderFunc: func(ctx context.Context, name, signature string) error {
			return nil
		},
	}

	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
			ResourceFactory: NewIdentityProviderResource,
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
					resource "ignition_identity_provider" "unit_oidc" {
						name        = "unit-test-oidc"
						type        = "oidc"
						client_id   = "test-client"
						provider_id = "https://auth.com"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_identity_provider.unit_oidc", "type", "oidc"),
					resource.TestCheckResourceAttr("ignition_identity_provider.unit_oidc", "client_id", "test-client"),
				),
			},
		},
	})
}

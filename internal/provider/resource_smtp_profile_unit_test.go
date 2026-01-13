package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitSMTPProfileResource_Create(t *testing.T) {
	mockClient := &client.MockClient{
		CreateSMTPProfileFunc: func(ctx context.Context, item client.ResourceResponse[client.SMTPProfileConfig]) (*client.ResourceResponse[client.SMTPProfileConfig], error) {
			item.Signature = "mock-signature-smtp"
			return &item, nil
		},
		GetSMTPProfileFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.SMTPProfileConfig], error) {
			return &client.ResourceResponse[client.SMTPProfileConfig]{
				Name:      "unit-test-smtp",
				Enabled:   true,
				Signature: "mock-signature-smtp",
				Config: client.SMTPProfileConfig{
					Profile: client.SMTPProfileProfile{Type: "smtp.classic"},
					Settings: client.SMTPProfileSettings{
						Settings: &client.SMTPProfileSettingsClassic{
							Hostname: "smtp.test.com",
							Port:     25,
						},
					},
				},
			}, nil
		},
		DeleteSMTPProfileFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_smtp_profile" "unit" {
						name     = "unit-test-smtp"
						hostname = "smtp.test.com"
						port     = 25
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_smtp_profile.unit", "name", "unit-test-smtp"),
					resource.TestCheckResourceAttr("ignition_smtp_profile.unit", "hostname", "smtp.test.com"),
					resource.TestCheckResourceAttr("ignition_smtp_profile.unit", "signature", "mock-signature-smtp"),
				),
			},
		},
	})
}
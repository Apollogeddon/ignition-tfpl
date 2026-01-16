package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitAuditProfileResource(t *testing.T) {
	mockClient := &client.MockClient{
		CreateAuditProfileFunc: func(ctx context.Context, ap client.ResourceResponse[client.AuditProfileConfig]) (*client.ResourceResponse[client.AuditProfileConfig], error) {
			ap.Signature = "sig-123"
			return &ap, nil
		},
		GetAuditProfileFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.AuditProfileConfig], error) {
			return &client.ResourceResponse[client.AuditProfileConfig]{
				Name:      name,
				Enabled:   boolPtr(true),
				Signature: "sig-123",
				Config: client.AuditProfileConfig{
					Profile: client.AuditProfileProfile{
						Type:          "database",
						RetentionDays: 90,
					},
					Settings: client.AuditProfileSettings{
						DatabaseName: "IgnitionDB",
						TableName:    "audit_events",
						AutoCreate:   true,
						PruneEnabled: false,
					},
				},
			}, nil
		},
		UpdateAuditProfileFunc: func(ctx context.Context, ap client.ResourceResponse[client.AuditProfileConfig]) (*client.ResourceResponse[client.AuditProfileConfig], error) {
			ap.Signature = "sig-456"
			return &ap, nil
		},
		DeleteAuditProfileFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_audit_profile" "test" {
						name        = "TestAuditProfile"
						type        = "database"
						database    = "IgnitionDB"
						table_name  = "audit_events"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "name", "TestAuditProfile"),
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "type", "database"),
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "database", "IgnitionDB"),
					resource.TestCheckResourceAttr("ignition_audit_profile.test", "table_name", "audit_events"),
				),
			},
		},
	})
}

package provider

import (
	"context"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitDatabaseConnectionResource(t *testing.T) {
	mockClient := &client.MockClient{
		CreateDatabaseConnectionFunc: func(ctx context.Context, db client.ResourceResponse[client.DatabaseConfig]) (*client.ResourceResponse[client.DatabaseConfig], error) {
			db.Signature = "sig-123"
			return &db, nil
		},
		GetDatabaseConnectionFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.DatabaseConfig], error) {
			return &client.ResourceResponse[client.DatabaseConfig]{
				Name:      name,
				Enabled:   boolPtr(true),
				Signature: "sig-123",
				Config: client.DatabaseConfig{
					Driver:     "PostgreSQL",
					Translator: "POSTGRESQL",
					ConnectURL: "jdbc:postgresql://localhost:5432/test",
					Username:   "dbuser",
				},
			}, nil
		},
		UpdateDatabaseConnectionFunc: func(ctx context.Context, db client.ResourceResponse[client.DatabaseConfig]) (*client.ResourceResponse[client.DatabaseConfig], error) {
			db.Signature = "sig-456"
			return &db, nil
		},
		DeleteDatabaseConnectionFunc: func(ctx context.Context, name, signature string) error {
			return nil
		},
		EncryptSecretFunc: func(ctx context.Context, plaintext string) (*client.IgnitionSecret, error) {
			return &client.IgnitionSecret{Type: "Embedded", Data: map[string]interface{}{"value": plaintext}}, nil
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
					resource "ignition_database_connection" "test" {
						name        = "TestDB"
						type        = "PostgreSQL"
						translator  = "POSTGRESQL"
						connect_url = "jdbc:postgresql://localhost:5432/test"
						username    = "dbuser"
						password    = "secret"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_database_connection.test", "name", "TestDB"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "type", "PostgreSQL"),
					resource.TestCheckResourceAttr("ignition_database_connection.test", "username", "dbuser"),
				),
			},
		},
	})
}

package datasources

import (
	"context"
	"testing"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/apollogeddon/ignition-tfpl/internal/provider/base"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitDataSources(t *testing.T) {
	mockClient := &client.MockClient{
		GetDatabaseConnectionFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.DatabaseConfig], error) {
			return &client.ResourceResponse[client.DatabaseConfig]{
				Name: name,
				Config: client.DatabaseConfig{
					Driver:     "test-driver",
					ConnectURL: "jdbc:test://localhost",
					Username:   "admin",
				},
			}, nil
		},
		GetSMTPProfileFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.SMTPProfileConfig], error) {
			return &client.ResourceResponse[client.SMTPProfileConfig]{
				Name: name,
				Config: client.SMTPProfileConfig{
					Settings: client.SMTPProfileSettings{
						Settings: &client.SMTPProfileSettingsClassic{
							Hostname: "smtp.test.com",
							Port:     25,
						},
					},
				},
			}, nil
		},
		GetStoreAndForwardFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.StoreAndForwardConfig], error) {
			return &client.ResourceResponse[client.StoreAndForwardConfig]{
				Name: name,
				Config: client.StoreAndForwardConfig{
					ForwardingPolicy: "ALL",
					BatchSize:        100,
				},
			}, nil
		},
		GetTagProviderFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.TagProviderConfig], error) {
			return &client.ResourceResponse[client.TagProviderConfig]{
				Name: name,
				Config: client.TagProviderConfig{
					Type: "standard",
				},
			}, nil
		},
		GetUserSourceFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.UserSourceConfig], error) {
			return &client.ResourceResponse[client.UserSourceConfig]{
				Name: name,
				Config: client.UserSourceConfig{
					Profile: client.UserSourceProfile{
						Type: "internal",
					},
				},
			}, nil
		},
	}

	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
			Client: mockClient,
			DataSourceFactories: []func() datasource.DataSource{
				NewDatabaseConnectionDataSource,
				NewSMTPProfileDataSource,
				NewStoreAndForwardDataSource,
				NewTagProviderDataSource,
				NewUserSourceDataSource,
			},
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
					data "ignition_database_connection" "test" { name = "db" }
					data "ignition_smtp_profile" "test" { name = "smtp" }
					data "ignition_store_forward" "test" { name = "sf" }
					data "ignition_tag_provider" "test" { name = "tags" }
					data "ignition_user_source" "test" { name = "users" }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ignition_database_connection.test", "type", "test-driver"),
					resource.TestCheckResourceAttr("data.ignition_smtp_profile.test", "hostname", "smtp.test.com"),
					resource.TestCheckResourceAttr("data.ignition_store_forward.test", "forwarding_policy", "ALL"),
					resource.TestCheckResourceAttr("data.ignition_tag_provider.test", "type", "standard"),
					resource.TestCheckResourceAttr("data.ignition_user_source.test", "type", "internal"),
				),
			},
		},
	})
}
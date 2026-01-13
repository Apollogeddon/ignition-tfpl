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

func TestUnitAlarmJournalResource_Create(t *testing.T) {
	mockClient := &client.MockClient{
		CreateAlarmJournalFunc: func(ctx context.Context, item client.ResourceResponse[client.AlarmJournalConfig]) (*client.ResourceResponse[client.AlarmJournalConfig], error) {
			if item.Name != "test-journal" {
				return nil, fmt.Errorf("expected name 'test-journal', got '%s'", item.Name)
			}
			if item.Config.Profile.Type != "DATASOURCE" {
				return nil, fmt.Errorf("expected type 'DATASOURCE', got '%s'", item.Config.Profile.Type)
			}
			if item.Config.Settings.Datasource != "db_connection" {
				return nil, fmt.Errorf("expected datasource 'db_connection', got '%s'", item.Config.Settings.Datasource)
			}
			
			// Simulate successful creation
			item.Signature = "mock-signature-journal"
			return &item, nil
		},
		GetAlarmJournalFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.AlarmJournalConfig], error) {
			if name != "test-journal" {
				return nil, fmt.Errorf("not found")
			}
			return &client.ResourceResponse[client.AlarmJournalConfig]{
				Name:      "test-journal",
				Enabled:   boolPtr(true),
				Signature: "mock-signature-journal",
				Config: client.AlarmJournalConfig{
					Profile: client.AlarmJournalProfile{Type: "DATASOURCE"},
					Settings: client.AlarmJournalSettings{
						Datasource: "db_connection",
						Advanced: &struct {
							TableName          string `json:"tableName,omitempty"`
							DataTableName      string `json:"dataTableName,omitempty"`
							UseStoreAndForward bool   `json:"useStoreAndForward,omitempty"`
						}{
							TableName: "alarm_events",
						},
						Events: &struct {
							MinPriority            string `json:"minPriority,omitempty"`
							StoreShelvedEvents     bool   `json:"storeShelvedEvents,omitempty"`
							StoreFromEnabledChange bool   `json:"storeFromEnabledChange,omitempty"`
						}{
							MinPriority: "Low",
						},
					},
				},
			}, nil
		},
		DeleteAlarmJournalFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_alarm_journal" "unit" {
						name         = "test-journal"
						type         = "DATASOURCE"
						datasource   = "db_connection"
						table_name   = "alarm_events"
						min_priority = "Low"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_journal.unit", "name", "test-journal"),
					resource.TestCheckResourceAttr("ignition_alarm_journal.unit", "signature", "mock-signature-journal"),
				),
			},
		},
	})
}

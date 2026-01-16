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

func TestUnitAlarmNotificationProfileResource_Create(t *testing.T) {
	currentHostname := "mock.smtp.com"
	currentSignature := "mock-signature-123"

	mockClient := &client.MockClient{
		CreateAlarmNotificationProfileFunc: func(ctx context.Context, anp client.ResourceResponse[client.AlarmNotificationProfileConfig]) (*client.ResourceResponse[client.AlarmNotificationProfileConfig], error) {
			if anp.Name != "unit-test-profile" {
				return nil, fmt.Errorf("expected name 'unit-test-profile', got '%s'", anp.Name)
			}
			if anp.Config.Profile.Type != "EmailNotificationProfileType" {
				return nil, fmt.Errorf("expected type 'EmailNotificationProfileType', got '%s'", anp.Config.Profile.Type)
			}
			
			// Simulate successful creation
			anp.Signature = "mock-signature-123"
			return &anp, nil
		},
		GetAlarmNotificationProfileFunc: func(ctx context.Context, name string) (*client.ResourceResponse[client.AlarmNotificationProfileConfig], error) {
			if name != "unit-test-profile" {
				return nil, fmt.Errorf("not found")
			}
			return &client.ResourceResponse[client.AlarmNotificationProfileConfig]{
				Name:      "unit-test-profile",
				Enabled:   boolPtr(true),
				Signature: currentSignature,
				Config: client.AlarmNotificationProfileConfig{
					Profile: client.AlarmNotificationProfileProfile{Type: "EmailNotificationProfileType"},
					Settings: map[string]any{
						"useSmtpProfile": false,
						"hostname":       currentHostname,
						"port":           25,
						"sslEnabled":     false,
					},
				},
			}, nil
		},
		UpdateAlarmNotificationProfileFunc: func(ctx context.Context, anp client.ResourceResponse[client.AlarmNotificationProfileConfig]) (*client.ResourceResponse[client.AlarmNotificationProfileConfig], error) {
			if anp.Name != "unit-test-profile" {
				return nil, fmt.Errorf("expected name 'unit-test-profile', got '%s'", anp.Name)
			}
			
			// In Ignition, settings for Email type are often nested
			settings := anp.Config.Settings
			if s, ok := settings["settings"].(map[string]any); ok {
				settings = s
			}
			
			if v, ok := settings["hostname"].(string); ok {
				currentHostname = v
			}
			currentSignature = "mock-signature-456" // Update current signature
			anp.Signature = currentSignature
			return &anp, nil
		},
		DeleteAlarmNotificationProfileFunc: func(ctx context.Context, name, signature string) error {
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
					resource "ignition_alarm_notification_profile" "unit" {
						name = "unit-test-profile"
						type = "EmailNotificationProfileType"
						email_config {
							use_smtp_profile = false
							hostname         = "mock.smtp.com"
							port             = 25
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.unit", "name", "unit-test-profile"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.unit", "signature", "mock-signature-123"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.unit", "email_config.hostname", "mock.smtp.com"),
				),
			},
			// Update Step
			{
				Config: `
					provider "ignition" {
						host  = "http://mock-host"
						token = "mock-token"
					}
					resource "ignition_alarm_notification_profile" "unit" {
						name = "unit-test-profile"
						type = "EmailNotificationProfileType"
						email_config {
							use_smtp_profile = false
							hostname         = "updated.smtp.com"
							port             = 25
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.unit", "email_config.hostname", "updated.smtp.com"),
					resource.TestCheckResourceAttr("ignition_alarm_notification_profile.unit", "signature", "mock-signature-456"),
				),
			},
		},
	})
}

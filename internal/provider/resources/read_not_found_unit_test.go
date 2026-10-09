package resources

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/apollogeddon/ignition-tofu/internal/client"
	"github.com/apollogeddon/ignition-tofu/internal/provider/base"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// errGatewayNotFound is what the client returns for a 404.
var errGatewayNotFound = fmt.Errorf("status: 404, body: : %w", client.ErrNotFound)

const smtpDeletedConfig = `
	provider "ignition" {
		host  = "http://mock-host"
		token = "mock-token"
	}
	resource "ignition_smtp_profile" "test" {
		name     = "smtp"
		hostname = "smtp.test.com"
		port     = 25
	}
`

// smtpMock is a gateway holding at most one SMTP profile, which the test can
// delete behind OpenTofu's back. getErr, when set, replaces the lookup result.
func smtpMock(exists *bool, creates *int, getErr *error) *client.MockClient {
	return &client.MockClient{
		CreateSMTPProfileFunc: func(_ context.Context, item client.ResourceResponse[client.SMTPProfileConfig]) (*client.ResourceResponse[client.SMTPProfileConfig], error) {
			*exists = true
			*creates++
			item.Signature = "sig"
			return &item, nil
		},
		GetSMTPProfileFunc: func(_ context.Context, name string) (*client.ResourceResponse[client.SMTPProfileConfig], error) {
			if *getErr != nil {
				return nil, *getErr
			}
			if !*exists {
				return nil, errGatewayNotFound
			}
			return &client.ResourceResponse[client.SMTPProfileConfig]{
				Name:      name,
				Enabled:   base.BoolPtr(true),
				Signature: "sig",
				Config: client.SMTPProfileConfig{
					Profile: client.SMTPProfileProfile{Type: "smtp.classic"},
					Settings: client.SMTPProfileSettings{
						Settings: &client.SMTPProfileSettingsClassic{Hostname: "smtp.test.com", Port: 25},
					},
				},
			}, nil
		},
		DeleteSMTPProfileFunc: func(context.Context, string, string) error {
			*exists = false
			return nil
		},
	}
}

func TestUnitRead_DeletedOutsideOpenTofuIsRecreated(t *testing.T) {
	exists, creates := false, 0
	var getErr error
	mockClient := smtpMock(&exists, &creates, &getErr)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
				ResourceFactory: NewSMTPProfileResource,
				Client:          mockClient,
			}),
		},
		Steps: []resource.TestStep{
			{Config: smtpDeletedConfig},
			{
				PreConfig: func() { exists = false }, // deleted on the gateway
				Config:    smtpDeletedConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ignition_smtp_profile.test", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})

	if creates != 2 {
		t.Errorf("expected the profile to be created twice (initial + recreate), got %d", creates)
	}
}

func TestUnitRead_OtherErrorsStillFail(t *testing.T) {
	exists, creates := false, 0
	var getErr error
	mockClient := smtpMock(&exists, &creates, &getErr)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
				ResourceFactory: NewSMTPProfileResource,
				Client:          mockClient,
			}),
		},
		Steps: []resource.TestStep{
			{Config: smtpDeletedConfig},
			{
				PreConfig:   func() { getErr = fmt.Errorf("status: 500, body: boom") },
				Config:      smtpDeletedConfig,
				ExpectError: regexp.MustCompile(`Error reading resource`),
			},
			{
				// let the post-test destroy succeed
				PreConfig: func() { getErr = nil },
				Config:    smtpDeletedConfig,
			},
		},
	})
}

func TestUnitDeviceRead_DeletedOutsideOpenTofuIsRecreated(t *testing.T) {
	exists, creates := false, 0
	mockClient := &client.MockClient{
		CreateDeviceFunc: func(_ context.Context, item client.ResourceResponse[client.DeviceConfig]) (*client.ResourceResponse[client.DeviceConfig], error) {
			exists = true
			creates++
			item.Signature = "sig"
			return &item, nil
		},
		GetDeviceFunc: func(_ context.Context, name string) (*client.ResourceResponse[client.DeviceConfig], error) {
			if !exists {
				return nil, errGatewayNotFound
			}
			return &client.ResourceResponse[client.DeviceConfig]{
				Name:      name,
				Type:      "ProgrammableSimulatorDevice",
				Enabled:   base.BoolPtr(true),
				Signature: "sig",
				Config:    client.DeviceConfig{"baseRate": 1000},
			}, nil
		},
		DeleteDeviceFunc: func(context.Context, string, string) error {
			exists = false
			return nil
		},
	}
	config := `
		provider "ignition" {
			host  = "http://mock-host"
			token = "mock-token"
		}
		resource "ignition_device" "test" {
			name       = "SimDevice"
			type       = "ProgrammableSimulatorDevice"
			parameters = "{\"baseRate\": 1000}"
		}
	`

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"ignition": providerserver.NewProtocol6WithError(&base.TestProvider{
				ResourceFactory: NewDeviceResource,
				Client:          mockClient,
			}),
		},
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() { exists = false },
				Config:    config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ignition_device.test", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})

	if creates != 2 {
		t.Errorf("expected the device to be created twice (initial + recreate), got %d", creates)
	}
}

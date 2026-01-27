package provider

import (
	"context"
	"os"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/apollogeddon/ignition-tfpl/internal/provider/datasources"
	"github.com/apollogeddon/ignition-tfpl/internal/provider/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure IgnitionProvider satisfies various provider interfaces.
var _ provider.Provider = &IgnitionProvider{}

// IgnitionProvider defines the provider implementation.
type IgnitionProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	client  client.IgnitionClient
}

// IgnitionProviderModel describes the provider data model.
type IgnitionProviderModel struct {
	Host             types.String `tfsdk:"host"`
	Token            types.String `tfsdk:"token"`
	AllowInsecureTLS types.Bool   `tfsdk:"allow_insecure_tls"`
}

func (p *IgnitionProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ignition"
	resp.Version = p.version
}

func (p *IgnitionProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Ignition provider allows you to manage Inductive Automation Ignition resources.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The base URL of the Ignition Gateway (e.g., http://localhost:8088).",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "The API token for authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"allow_insecure_tls": schema.BoolAttribute{
				Description: "Whether to allow insecure TLS connections (e.g., self-signed certs).",
				Optional:    true,
			},
		},
	}
}

func (p *IgnitionProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data IgnitionProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Host",
			"The provider host cannot be unknown. Please verify your Terraform configuration.",
		)
	}

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Token",
			"The provider token cannot be unknown. Please verify your Terraform configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := data.Host.ValueString()
	token := data.Token.ValueString()
	allowInsecure := data.AllowInsecureTLS.ValueBool()

	// Default to environment variables if not configured
	if host == "" {
		host = os.Getenv("IGNITION_HOST")
	}
	if token == "" {
		token = os.Getenv("IGNITION_TOKEN")
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Host",
			"The provider host must be configured via the 'host' attribute or IGNITION_HOST environment variable.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Token",
			"The provider token must be configured via the 'token' attribute or IGNITION_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var apiClient client.IgnitionClient
	var err error

	if p.client != nil {
		apiClient = p.client
	} else {
		apiClient, err = client.NewClient(host, token, allowInsecure)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create Ignition API Client",
				"An unexpected error occurred when creating the Ignition API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Ignition Client Error: "+err.Error(),
			)
			return
		}
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *IgnitionProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewDatabaseConnectionResource,
		resources.NewTagProviderResource,
		resources.NewUserSourceResource,
		resources.NewProjectResource,
		resources.NewAuditProfileResource,
		resources.NewAlarmNotificationProfileResource,
		resources.NewOpcUaConnectionResource,
		resources.NewAlarmJournalResource,
		resources.NewSMTPProfileResource,
		resources.NewStoreAndForwardResource,
		resources.NewIdentityProviderResource,
		resources.NewGanOutgoingResource,
		resources.NewRedundancyResource,
		resources.NewGanGeneralSettingsResource,
		resources.NewDeviceResource,
	}
}

func (p *IgnitionProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewDatabaseConnectionDataSource,
		datasources.NewProjectDataSource,
		datasources.NewUserSourceDataSource,
		datasources.NewTagProviderDataSource,
		datasources.NewSMTPProfileDataSource,
		datasources.NewStoreAndForwardDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IgnitionProvider{
			version: version,
		}
	}
}

package base

import (
	"context"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestProvider is a minimal implementation of provider.Provider for unit testing.
type TestProvider struct {
	ResourceFactory     func() resource.Resource
	ResourceFactories   []func() resource.Resource
	DataSourceFactory   func() datasource.DataSource
	DataSourceFactories []func() datasource.DataSource
	Client              client.IgnitionClient
}

func (p *TestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ignition"
}

func (p *TestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host":  schema.StringAttribute{Optional: true},
			"token": schema.StringAttribute{Optional: true},
		},
	}
}

func (p *TestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.DataSourceData = p.Client
	resp.ResourceData = p.Client
}

func (p *TestProvider) Resources(_ context.Context) []func() resource.Resource {
	var factories []func() resource.Resource
	if p.ResourceFactory != nil {
		factories = append(factories, p.ResourceFactory)
	}
	factories = append(factories, p.ResourceFactories...)
	return factories
}

func (p *TestProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	var factories []func() datasource.DataSource
	if p.DataSourceFactory != nil {
		factories = append(factories, p.DataSourceFactory)
	}
	factories = append(factories, p.DataSourceFactories...)
	return factories
}

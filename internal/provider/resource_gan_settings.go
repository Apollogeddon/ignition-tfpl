package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GanGeneralSettingsResource{}

func NewGanGeneralSettingsResource() resource.Resource {
	return &GanGeneralSettingsResource{}
}

// GanGeneralSettingsResource defines the resource implementation.
type GanGeneralSettingsResource struct {
	client  client.IgnitionClient
	generic GenericIgnitionResource[client.GanGeneralSettingsConfig, GanGeneralSettingsResourceModel]
}

// GanGeneralSettingsResourceModel describes the resource data model.
type GanGeneralSettingsResourceModel struct {
	BaseResourceModel
	RequireSSL                  types.Bool    `tfsdk:"require_ssl"`
	RequireTwoWayAuth           types.Bool    `tfsdk:"require_two_way_auth"`
	AllowIncoming               types.Bool    `tfsdk:"allow_incoming"`
	SecurityPolicy              types.String  `tfsdk:"security_policy"`
	Whitelist                   types.String  `tfsdk:"whitelist"`
	AllowedProxyHops            types.Float64 `tfsdk:"allowed_proxy_hops"`
	WebsocketSessionIdleTimeout types.Float64 `tfsdk:"websocket_session_idle_timeout"`
	TempFilesMaxAgeHours        types.Float64 `tfsdk:"temp_files_max_age_hours"`
}

func (r *GanGeneralSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gan_settings"
}

func (r *GanGeneralSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages General Gateway Network Settings. This is a singleton resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "Internal name for the resource (fixed to 'gateway-network-settings').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("gateway-network-settings"),
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"require_ssl": schema.BoolAttribute{
				Description: "If true, only connections that use SSL will be allowed.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"require_two_way_auth": schema.BoolAttribute{
				Description: "Enforces two-way SSL authentication.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_incoming": schema.BoolAttribute{
				Description: "If false, only outward connections will be allowed.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"security_policy": schema.StringAttribute{
				Description: "Dictates what connections are allowed (Unrestricted, ApprovedOnly, SpecifiedList).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ApprovedOnly"),
				Validators: []validator.String{
					stringvalidator.OneOf("Unrestricted", "ApprovedOnly", "SpecifiedList"),
				},
			},
			"whitelist": schema.StringAttribute{
				Description: "Comma-separated list of Gateway names allowed to connect if security_policy is SpecifiedList.",
				Optional:    true,
			},
			"allowed_proxy_hops": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(0),
			},
			"websocket_session_idle_timeout": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(30000),
			},
			"temp_files_max_age_hours": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(24),
			},
			"signature": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *GanGeneralSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(client.IgnitionClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.IgnitionClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = c
	r.generic = GenericIgnitionResource[client.GanGeneralSettingsConfig, GanGeneralSettingsResourceModel]{
		Client:       c,
		Handler:      r,
		Module:       "ignition",
		ResourceType: "gateway-network-settings",
		CreateFunc:   c.UpdateGanGeneralSettings,
		GetFunc: func(ctx context.Context, _ string) (*client.ResourceResponse[client.GanGeneralSettingsConfig], error) {
			return c.GetGanGeneralSettings(ctx)
		},
		UpdateFunc: c.UpdateGanGeneralSettings,
		DeleteFunc: func(ctx context.Context, name, signature string) error {
			return nil
		},
	}
}

func (r *GanGeneralSettingsResource) MapPlanToClient(ctx context.Context, model *GanGeneralSettingsResourceModel) (client.GanGeneralSettingsConfig, error) {
	return client.GanGeneralSettingsConfig{
		RequireSSL:                  model.RequireSSL.ValueBool(),
		RequireTwoWayAuth:           model.RequireTwoWayAuth.ValueBool(),
		AllowIncoming:               model.AllowIncoming.ValueBool(),
		SecurityPolicy:              model.SecurityPolicy.ValueString(),
		Whitelist:                   model.Whitelist.ValueString(),
		AllowedProxyHops:            model.AllowedProxyHops.ValueFloat64(),
		WebsocketSessionIdleTimeout: model.WebsocketSessionIdleTimeout.ValueFloat64(),
		TempFilesMaxAgeHours:        model.TempFilesMaxAgeHours.ValueFloat64(),
	}, nil
}

func (r *GanGeneralSettingsResource) MapClientToState(ctx context.Context, name string, config *client.GanGeneralSettingsConfig, model *GanGeneralSettingsResourceModel) error {
	model.Name = types.StringValue(name)
	model.RequireSSL = types.BoolValue(config.RequireSSL)
	model.RequireTwoWayAuth = types.BoolValue(config.RequireTwoWayAuth)
	model.AllowIncoming = types.BoolValue(config.AllowIncoming)
	model.SecurityPolicy = types.StringValue(config.SecurityPolicy)
	model.Whitelist = stringToNullableString(config.Whitelist)
	model.AllowedProxyHops = types.Float64Value(config.AllowedProxyHops)
	model.WebsocketSessionIdleTimeout = types.Float64Value(config.WebsocketSessionIdleTimeout)
	model.TempFilesMaxAgeHours = types.Float64Value(config.TempFilesMaxAgeHours)
	return nil
}

func (r *GanGeneralSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GanGeneralSettingsResourceModel
	// Ensure name is fixed
	data.Name = types.StringValue("gateway-network-settings")
	r.generic.Create(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *GanGeneralSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GanGeneralSettingsResourceModel
	r.generic.Read(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *GanGeneralSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GanGeneralSettingsResourceModel
	r.generic.Update(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *GanGeneralSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GanGeneralSettingsResourceModel
	r.generic.Delete(ctx, req, resp, &data, &data.BaseResourceModel)
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IdentityProviderResource{}
var _ resource.ResourceWithImportState = &IdentityProviderResource{}

func NewIdentityProviderResource() resource.Resource {
	return &IdentityProviderResource{}
}

// IdentityProviderResource defines the resource implementation.
type IdentityProviderResource struct {
	client client.IgnitionClient
}

// IdentityProviderResourceModel describes the resource data model.
type IdentityProviderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	Type                     types.String `tfsdk:"type"`
	UserSource               types.String `tfsdk:"user_source"`
	SessionInactivityTimeout types.Float64 `tfsdk:"session_inactivity_timeout"`
	SessionExp               types.Float64 `tfsdk:"session_expiration"`
	RememberMeExp            types.Float64 `tfsdk:"remember_me_expiration"`
	Signature                types.String `tfsdk:"signature"`
}

func (r *IdentityProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider"
}

func (r *IdentityProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Identity Provider in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the identity provider.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the identity provider.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the identity provider is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "The type of the identity provider (currently only 'internal' is supported by this resource).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("internal"),
				},
			},
			"user_source": schema.StringAttribute{
				Description: "The name of the User Source Profile used to authenticate users (for 'internal' type).",
				Optional:    true,
			},
			"session_inactivity_timeout": schema.Float64Attribute{
				Description: "Minutes before expiring a session due to user inactivity.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"session_expiration": schema.Float64Attribute{
				Description: "Maximum minutes a session may exist before it is expired.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"remember_me_expiration": schema.Float64Attribute{
				Description: "Maximum hours a user will be remembered.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource.",
				Computed:    true,
			},
		},
	}
}

func (r *IdentityProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.IgnitionClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.IgnitionClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *IdentityProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IdentityProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internalConfig := client.IdentityProviderInternalConfig{
		UserSource:               data.UserSource.ValueString(),
		SessionInactivityTimeout: data.SessionInactivityTimeout.ValueFloat64(),
		SessionExp:               data.SessionExp.ValueFloat64(),
		RememberMeExp:            data.RememberMeExp.ValueFloat64(),
		AuthMethods: []client.IdentityProviderAuthMethod{
			{
				Type: "basic",
				Config: map[string]any{},
			},
		},
	}

	config := client.IdentityProviderConfig{
		Type:   data.Type.ValueString(),
		Config: internalConfig,
	}

	res := client.ResourceResponse[client.IdentityProviderConfig]{
		Name:    data.Name.ValueString(),
		Enabled: data.Enabled.ValueBool(),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateIdentityProvider(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Identity Provider", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IdentityProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetIdentityProvider(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Identity Provider", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Enabled = types.BoolValue(res.Enabled)
	data.Description = types.StringValue(res.Description)
	data.Type = types.StringValue(res.Config.Type)

	// Unmarshal internal config if applicable
	if res.Config.Type == "internal" {
		var internalConfig client.IdentityProviderInternalConfig
		configBytes, _ := json.Marshal(res.Config.Config)
		if err := json.Unmarshal(configBytes, &internalConfig); err == nil {
			data.UserSource = types.StringValue(internalConfig.UserSource)
			data.SessionInactivityTimeout = types.Float64Value(internalConfig.SessionInactivityTimeout)
			data.SessionExp = types.Float64Value(internalConfig.SessionExp)
			data.RememberMeExp = types.Float64Value(internalConfig.RememberMeExp)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IdentityProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internalConfig := client.IdentityProviderInternalConfig{
		UserSource:               data.UserSource.ValueString(),
		SessionInactivityTimeout: data.SessionInactivityTimeout.ValueFloat64(),
		SessionExp:               data.SessionExp.ValueFloat64(),
		RememberMeExp:            data.RememberMeExp.ValueFloat64(),
		AuthMethods: []client.IdentityProviderAuthMethod{
			{
				Type: "basic",
				Config: map[string]any{},
			},
		},
	}

	config := client.IdentityProviderConfig{
		Type:   data.Type.ValueString(),
		Config: internalConfig,
	}

	res := client.ResourceResponse[client.IdentityProviderConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   data.Enabled.ValueBool(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateIdentityProvider(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Identity Provider", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IdentityProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIdentityProvider(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Identity Provider", err.Error())
		return
	}
}

func (r *IdentityProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

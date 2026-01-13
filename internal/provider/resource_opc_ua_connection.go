package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OpcUaConnectionResource{}
var _ resource.ResourceWithImportState = &OpcUaConnectionResource{}

func NewOpcUaConnectionResource() resource.Resource {
	return &OpcUaConnectionResource{}
}

// OpcUaConnectionResource defines the resource implementation.
type OpcUaConnectionResource struct {
	client client.IgnitionClient
}

// OpcUaConnectionResourceModel describes the resource data model.
type OpcUaConnectionResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Type           types.String `tfsdk:"type"`
	DiscoveryURL   types.String `tfsdk:"discovery_url"`
	EndpointURL    types.String `tfsdk:"endpoint_url"`
	SecurityPolicy types.String `tfsdk:"security_policy"`
	SecurityMode   types.String `tfsdk:"security_mode"`
	Signature      types.String `tfsdk:"signature"`
}

func (r *OpcUaConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_opc_ua_connection"
}

func (r *OpcUaConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OPC UA Connection in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the OPC UA connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the OPC UA connection.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the OPC UA connection is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "The type of the OPC UA connection (e.g., com.inductiveautomation.OpcUaServerType).",
				Optional:    true,
				Computed:    true,
				// Default to standard OPC UA Server type if not provided
			},
			"discovery_url": schema.StringAttribute{
				Description: "The discovery URL of the OPC UA server.",
				Required:    true,
			},
			"endpoint_url": schema.StringAttribute{
				Description: "The endpoint URL of the OPC UA server.",
				Required:    true,
			},
			"security_policy": schema.StringAttribute{
				Description: "The security policy to use (e.g., None, Basic256Sha256).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"None",
						"Basic128Rsa15",
						"Basic256",
						"Basic256Sha256",
						"Aes128_Sha256_RsaOaep",
						"Aes256_Sha256_RsaPss",
					),
				},
			},
			"security_mode": schema.StringAttribute{
				Description: "The security mode to use (e.g., None, SignAndEncrypt).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"None",
						"Sign",
						"SignAndEncrypt",
					),
				},
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *OpcUaConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OpcUaConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpcUaConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default type if not set
	profileType := "com.inductiveautomation.OpcUaServerType"
	if !data.Type.IsNull() {
		profileType = data.Type.ValueString()
	}

	config := client.OpcUaConnectionConfig{
		Profile: client.OpcUaConnectionProfile{
			Type: profileType,
		},
		Settings: client.OpcUaConnectionSettings{
			Endpoint: client.OpcUaConnectionEndpoint{
				DiscoveryURL:   data.DiscoveryURL.ValueString(),
				EndpointURL:    data.EndpointURL.ValueString(),
				SecurityPolicy: data.SecurityPolicy.ValueString(),
				SecurityMode:   data.SecurityMode.ValueString(),
			},
		},
	}

	res := client.ResourceResponse[client.OpcUaConnectionConfig]{
		Name:    data.Name.ValueString(),
		Enabled: boolPtr(data.Enabled.ValueBool()),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateOpcUaConnection(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OPC UA connection", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	data.Name = types.StringValue(created.Name)
	if created.Enabled != nil {
		data.Enabled = types.BoolValue(*created.Enabled)
	} else {
		data.Enabled = types.BoolValue(true)
	}
	
	if created.Description != "" {
		data.Description = types.StringValue(created.Description)
	} else {
		data.Description = types.StringNull()
	}

	data.Type = types.StringValue(created.Config.Profile.Type)
	data.DiscoveryURL = types.StringValue(created.Config.Settings.Endpoint.DiscoveryURL)
	data.EndpointURL = types.StringValue(created.Config.Settings.Endpoint.EndpointURL)
	data.SecurityPolicy = types.StringValue(created.Config.Settings.Endpoint.SecurityPolicy)
	data.SecurityMode = types.StringValue(created.Config.Settings.Endpoint.SecurityMode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpcUaConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpcUaConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetOpcUaConnection(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading OPC UA connection", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Name = types.StringValue(res.Name)
	if res.Enabled != nil {
		data.Enabled = types.BoolValue(*res.Enabled)
	} else {
		data.Enabled = types.BoolValue(true)
	}
	
	if res.Description != "" {
		data.Description = types.StringValue(res.Description)
	} else {
		data.Description = types.StringNull()
	}

	data.Type = types.StringValue(res.Config.Profile.Type)
	data.DiscoveryURL = types.StringValue(res.Config.Settings.Endpoint.DiscoveryURL)
	data.EndpointURL = types.StringValue(res.Config.Settings.Endpoint.EndpointURL)
	data.SecurityPolicy = types.StringValue(res.Config.Settings.Endpoint.SecurityPolicy)
	data.SecurityMode = types.StringValue(res.Config.Settings.Endpoint.SecurityMode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpcUaConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OpcUaConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default type if not set
	profileType := "com.inductiveautomation.OpcUaServerType"
	if !data.Type.IsNull() {
		profileType = data.Type.ValueString()
	}

	config := client.OpcUaConnectionConfig{
		Profile: client.OpcUaConnectionProfile{
			Type: profileType,
		},
		Settings: client.OpcUaConnectionSettings{
			Endpoint: client.OpcUaConnectionEndpoint{
				DiscoveryURL:   data.DiscoveryURL.ValueString(),
				EndpointURL:    data.EndpointURL.ValueString(),
				SecurityPolicy: data.SecurityPolicy.ValueString(),
				SecurityMode:   data.SecurityMode.ValueString(),
			},
		},
	}

	res := client.ResourceResponse[client.OpcUaConnectionConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   boolPtr(data.Enabled.ValueBool()),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateOpcUaConnection(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating OPC UA connection", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	data.Type = types.StringValue(updated.Config.Profile.Type)
	data.DiscoveryURL = types.StringValue(updated.Config.Settings.Endpoint.DiscoveryURL)
	data.EndpointURL = types.StringValue(updated.Config.Settings.Endpoint.EndpointURL)
	data.SecurityPolicy = types.StringValue(updated.Config.Settings.Endpoint.SecurityPolicy)
	data.SecurityMode = types.StringValue(updated.Config.Settings.Endpoint.SecurityMode)
	
	if updated.Description != "" {
		data.Description = types.StringValue(updated.Description)
	}

	if updated.Enabled != nil {
		data.Enabled = types.BoolValue(*updated.Enabled)
	} else {
		data.Enabled = types.BoolValue(true)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpcUaConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpcUaConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOpcUaConnection(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OPC UA connection", err.Error())
		return
	}
}

func (r *OpcUaConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
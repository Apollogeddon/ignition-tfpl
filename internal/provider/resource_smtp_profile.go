package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SMTPProfileResource{}
var _ resource.ResourceWithImportState = &SMTPProfileResource{}

func NewSMTPProfileResource() resource.Resource {
	return &SMTPProfileResource{}
}

// SMTPProfileResource defines the resource implementation.
type SMTPProfileResource struct {
	client client.IgnitionClient
}

// SMTPProfileResourceModel describes the resource data model.
type SMTPProfileResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Hostname        types.String `tfsdk:"hostname"`
	Port            types.Int64  `tfsdk:"port"`
	UseSslPort      types.Bool   `tfsdk:"use_ssl_port"`
	StartTlsEnabled types.Bool   `tfsdk:"start_tls_enabled"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	Signature       types.String `tfsdk:"signature"`
}

func (r *SMTPProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smtp_profile"
}

func (r *SMTPProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an SMTP Email Profile in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the SMTP profile.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the SMTP profile.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the SMTP profile is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname of the SMTP server.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Port SMTP service is running on.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(25),
			},
			"use_ssl_port": schema.BoolAttribute{
				Description: "Connect using dedicated SSL/TLS port.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"start_tls_enabled": schema.BoolAttribute{
				Description: "Connect using STARTTLS.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"username": schema.StringAttribute{
				Description: "The username for logging into the email server.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for logging into the email server.",
				Optional:    true,
				Sensitive:   true,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *SMTPProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SMTPProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SMTPProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.SMTPProfileConfig{
		Profile: client.SMTPProfileProfile{
			Type: "smtp.classic",
		},
		Settings: client.SMTPProfileSettings{
			Settings: &client.SMTPProfileSettingsClassic{
				Hostname:        data.Hostname.ValueString(),
				Port:            int(data.Port.ValueInt64()),
				UseSslPort:      data.UseSslPort.ValueBool(),
				StartTlsEnabled: data.StartTlsEnabled.ValueBool(),
			},
		},
	}

	if !data.Username.IsNull() {
		config.Settings.Settings.Username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		config.Settings.Settings.Password = data.Password.ValueString()
	}

	res := client.ResourceResponse[client.SMTPProfileConfig]{
		Name:    data.Name.ValueString(),
		Enabled: data.Enabled.ValueBool(),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateSMTPProfile(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SMTP profile", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	
	if created.Description != "" {
		data.Description = types.StringValue(created.Description)
	} else {
		data.Description = types.StringNull()
	}

	if created.Config.Settings.Settings != nil {
		data.Hostname = types.StringValue(created.Config.Settings.Settings.Hostname)
		data.Port = types.Int64Value(int64(created.Config.Settings.Settings.Port))
		data.UseSslPort = types.BoolValue(created.Config.Settings.Settings.UseSslPort)
		data.StartTlsEnabled = types.BoolValue(created.Config.Settings.Settings.StartTlsEnabled)
		if created.Config.Settings.Settings.Username != "" {
			data.Username = types.StringValue(created.Config.Settings.Settings.Username)
		} else {
			data.Username = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SMTPProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SMTPProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetSMTPProfile(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SMTP profile", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Enabled = types.BoolValue(res.Enabled)
	
	if res.Description != "" {
		data.Description = types.StringValue(res.Description)
	} else {
		data.Description = types.StringNull()
	}

	if res.Config.Settings.Settings != nil {
		data.Hostname = types.StringValue(res.Config.Settings.Settings.Hostname)
		data.Port = types.Int64Value(int64(res.Config.Settings.Settings.Port))
		data.UseSslPort = types.BoolValue(res.Config.Settings.Settings.UseSslPort)
		data.StartTlsEnabled = types.BoolValue(res.Config.Settings.Settings.StartTlsEnabled)
		if res.Config.Settings.Settings.Username != "" {
			data.Username = types.StringValue(res.Config.Settings.Settings.Username)
		} else {
			data.Username = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SMTPProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SMTPProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.SMTPProfileConfig{
		Profile: client.SMTPProfileProfile{
			Type: "smtp.classic",
		},
		Settings: client.SMTPProfileSettings{
			Settings: &client.SMTPProfileSettingsClassic{
				Hostname:        data.Hostname.ValueString(),
				Port:            int(data.Port.ValueInt64()),
				UseSslPort:      data.UseSslPort.ValueBool(),
				StartTlsEnabled: data.StartTlsEnabled.ValueBool(),
			},
		},
	}

	if !data.Username.IsNull() {
		config.Settings.Settings.Username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		config.Settings.Settings.Password = data.Password.ValueString()
	}

	res := client.ResourceResponse[client.SMTPProfileConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   data.Enabled.ValueBool(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateSMTPProfile(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating SMTP profile", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	
	if updated.Description != "" {
		data.Description = types.StringValue(updated.Description)
	}

	if updated.Config.Settings.Settings != nil {
		data.Hostname = types.StringValue(updated.Config.Settings.Settings.Hostname)
		data.Port = types.Int64Value(int64(updated.Config.Settings.Settings.Port))
		data.UseSslPort = types.BoolValue(updated.Config.Settings.Settings.UseSslPort)
		data.StartTlsEnabled = types.BoolValue(updated.Config.Settings.Settings.StartTlsEnabled)
		if updated.Config.Settings.Settings.Username != "" {
			data.Username = types.StringValue(updated.Config.Settings.Settings.Username)
		} else {
			data.Username = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SMTPProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SMTPProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSMTPProfile(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting SMTP profile", err.Error())
		return
	}
}

func (r *SMTPProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

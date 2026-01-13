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
var _ resource.Resource = &AlarmNotificationProfileResource{}
var _ resource.ResourceWithImportState = &AlarmNotificationProfileResource{}

func NewAlarmNotificationProfileResource() resource.Resource {
	return &AlarmNotificationProfileResource{}
}

// AlarmNotificationProfileResource defines the resource implementation.
type AlarmNotificationProfileResource struct {
	client client.IgnitionClient
}

// AlarmNotificationProfileResourceModel describes the resource data model.
type AlarmNotificationProfileResourceModel struct {
	Id          types.String                      `tfsdk:"id"`
	Name        types.String                      `tfsdk:"name"`
	Description types.String                      `tfsdk:"description"`
	Enabled     types.Bool                        `tfsdk:"enabled"`
	Type        types.String                      `tfsdk:"type"`
	EmailConfig *AlarmNotificationProfileEmailModel `tfsdk:"email_config"`
	Signature   types.String                      `tfsdk:"signature"`
}

type AlarmNotificationProfileEmailModel struct {
	UseSMTPProfile types.Bool   `tfsdk:"use_smtp_profile"`
	EmailProfile   types.String `tfsdk:"email_profile"`
	Hostname       types.String `tfsdk:"hostname"`
	Port           types.Int64  `tfsdk:"port"`
	SSLEnabled     types.Bool   `tfsdk:"ssl_enabled"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
}

func (r *AlarmNotificationProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alarm_notification_profile"
}

func (r *AlarmNotificationProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Alarm Notification Profile in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the alarm notification profile.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the alarm notification profile.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the alarm notification profile is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "The type of the alarm notification profile (e.g., EmailNotificationProfileType).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"EmailNotificationProfileType",
						// Add others as they are implemented
					),
				},
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"email_config": schema.SingleNestedBlock{
				Description: "Configuration for Email Notification Profiles.",
				Attributes: map[string]schema.Attribute{
					"use_smtp_profile": schema.BoolAttribute{
						Description: "Whether to use an existing SMTP profile.",
						Required:    true,
					},
					"email_profile": schema.StringAttribute{
						Description: "The name of the SMTP profile to use (if use_smtp_profile is true).",
						Optional:    true,
					},
					"hostname": schema.StringAttribute{
						Description: "SMTP Server Hostname (if use_smtp_profile is false).",
						Optional:    true,
						Computed:    true,
					},
					"port": schema.Int64Attribute{
						Description: "SMTP Server Port (if use_smtp_profile is false).",
						Optional:    true,
						Computed:    true,
					},
					"ssl_enabled": schema.BoolAttribute{
						Description: "Enable SSL/TLS (if use_smtp_profile is false).",
						Optional:    true,
						Computed:    true,
					},
					"username": schema.StringAttribute{
						Description: "SMTP Username (if use_smtp_profile is false).",
						Optional:    true,
					},
					"password": schema.StringAttribute{
						Description: "SMTP Password (if use_smtp_profile is false).",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},
	}
}

func (r *AlarmNotificationProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlarmNotificationProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AlarmNotificationProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the client configuration
	config := client.AlarmNotificationProfileConfig{
		Profile: client.AlarmNotificationProfileProfile{
			Type: data.Type.ValueString(),
		},
		Settings: make(map[string]any),
	}

	if data.Type.ValueString() == "EmailNotificationProfileType" && data.EmailConfig != nil {
		emailSettings := make(map[string]any)
		emailSettings["useSmtpProfile"] = data.EmailConfig.UseSMTPProfile.ValueBool()
		
		if !data.EmailConfig.EmailProfile.IsNull() {
			emailSettings["emailProfile"] = data.EmailConfig.EmailProfile.ValueString()
		}
		if !data.EmailConfig.Hostname.IsNull() {
			emailSettings["hostname"] = data.EmailConfig.Hostname.ValueString()
		}
		if !data.EmailConfig.Port.IsNull() {
			emailSettings["port"] = int(data.EmailConfig.Port.ValueInt64())
		}
		if !data.EmailConfig.SSLEnabled.IsNull() {
			emailSettings["sslEnabled"] = data.EmailConfig.SSLEnabled.ValueBool()
		}
		if !data.EmailConfig.Username.IsNull() {
			emailSettings["username"] = data.EmailConfig.Username.ValueString()
		}
		if !data.EmailConfig.Password.IsNull() {
			emailSettings["password"] = data.EmailConfig.Password.ValueString()
		}
		config.Settings = emailSettings
	}

	res := client.ResourceResponse[client.AlarmNotificationProfileConfig]{
		Name:    data.Name.ValueString(),
		Enabled: data.Enabled.ValueBool(),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateAlarmNotificationProfile(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating alarm notification profile", err.Error())
		return
	}

	// Refresh from API to get server-side defaults
	refreshed, err := r.client.GetAlarmNotificationProfile(ctx, created.Name)
	if err == nil {
		created = refreshed
	}

	// Map response back to model
	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	data.Name = types.StringValue(created.Name)
	data.Enabled = types.BoolValue(created.Enabled)
	
	if created.Description != "" {
		data.Description = types.StringValue(created.Description)
	} else {
		data.Description = types.StringNull()
	}
	
	// Only map back settings if it was an Email type
	if created.Config.Profile.Type == "EmailNotificationProfileType" && data.EmailConfig != nil {
		settings := created.Config.Settings
		if settings == nil {
			settings = make(map[string]any)
		}
		
		if data.EmailConfig == nil {
			data.EmailConfig = &AlarmNotificationProfileEmailModel{}
		}
		
		if v, ok := settings["useSmtpProfile"].(bool); ok {
			data.EmailConfig.UseSMTPProfile = types.BoolValue(v)
		}
		
		if v, ok := settings["emailProfile"].(string); ok && v != "" {
			data.EmailConfig.EmailProfile = types.StringValue(v)
		} else {
			data.EmailConfig.EmailProfile = types.StringNull()
		}
		
		if v, ok := settings["hostname"].(string); ok && v != "" {
			data.EmailConfig.Hostname = types.StringValue(v)
		} else {
			data.EmailConfig.Hostname = types.StringNull()
		}
		
		// Handle port being float64 (from JSON unmarshal) or int
		if v, ok := settings["port"].(float64); ok && v != 0 {
			data.EmailConfig.Port = types.Int64Value(int64(v))
		} else if v, ok := settings["port"].(int); ok && v != 0 {
			data.EmailConfig.Port = types.Int64Value(int64(v))
		} else {
			data.EmailConfig.Port = types.Int64Null()
		}
		
		if v, ok := settings["sslEnabled"].(bool); ok {
			data.EmailConfig.SSLEnabled = types.BoolValue(v)
		} else {
			data.EmailConfig.SSLEnabled = types.BoolNull()
		}
		
		if v, ok := settings["username"].(string); ok && v != "" {
			data.EmailConfig.Username = types.StringValue(v)
		} else {
			data.EmailConfig.Username = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlarmNotificationProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AlarmNotificationProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetAlarmNotificationProfile(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading alarm notification profile", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Name = types.StringValue(res.Name)
	data.Enabled = types.BoolValue(res.Enabled)
	
	if res.Description != "" {
		data.Description = types.StringValue(res.Description)
	} else {
		data.Description = types.StringNull()
	}

	if res.Config.Profile.Type == "EmailNotificationProfileType" && res.Config.Settings != nil {
		data.Type = types.StringValue(res.Config.Profile.Type)
		settings := res.Config.Settings
		
		if data.EmailConfig == nil {
			data.EmailConfig = &AlarmNotificationProfileEmailModel{}
		}
		
		if v, ok := settings["useSmtpProfile"].(bool); ok {
			data.EmailConfig.UseSMTPProfile = types.BoolValue(v)
		}
		
		if v, ok := settings["emailProfile"].(string); ok && v != "" {
			data.EmailConfig.EmailProfile = types.StringValue(v)
		} else {
			data.EmailConfig.EmailProfile = types.StringNull()
		}
		
		if v, ok := settings["hostname"].(string); ok && v != "" {
			data.EmailConfig.Hostname = types.StringValue(v)
		} else {
			data.EmailConfig.Hostname = types.StringNull()
		}
		
		if v, ok := settings["port"].(float64); ok && v != 0 {
			data.EmailConfig.Port = types.Int64Value(int64(v))
		} else if v, ok := settings["port"].(int); ok && v != 0 {
			data.EmailConfig.Port = types.Int64Value(int64(v))
		} else {
			data.EmailConfig.Port = types.Int64Null()
		}
		
		if v, ok := settings["sslEnabled"].(bool); ok {
			data.EmailConfig.SSLEnabled = types.BoolValue(v)
		}
		
		if v, ok := settings["username"].(string); ok && v != "" {
			data.EmailConfig.Username = types.StringValue(v)
		} else {
			data.EmailConfig.Username = types.StringNull()
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlarmNotificationProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AlarmNotificationProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.AlarmNotificationProfileConfig{
		Profile: client.AlarmNotificationProfileProfile{
			Type: data.Type.ValueString(),
		},
		Settings: make(map[string]any),
	}

	if data.Type.ValueString() == "EmailNotificationProfileType" && data.EmailConfig != nil {
		emailSettings := make(map[string]any)
		emailSettings["useSmtpProfile"] = data.EmailConfig.UseSMTPProfile.ValueBool()
		
		if !data.EmailConfig.EmailProfile.IsNull() {
			emailSettings["emailProfile"] = data.EmailConfig.EmailProfile.ValueString()
		}
		if !data.EmailConfig.Hostname.IsNull() {
			emailSettings["hostname"] = data.EmailConfig.Hostname.ValueString()
		}
		if !data.EmailConfig.Port.IsNull() {
			emailSettings["port"] = int(data.EmailConfig.Port.ValueInt64())
		}
		if !data.EmailConfig.SSLEnabled.IsNull() {
			emailSettings["sslEnabled"] = data.EmailConfig.SSLEnabled.ValueBool()
		}
		if !data.EmailConfig.Username.IsNull() {
			emailSettings["username"] = data.EmailConfig.Username.ValueString()
		}
		if !data.EmailConfig.Password.IsNull() {
			emailSettings["password"] = data.EmailConfig.Password.ValueString()
		}
		config.Settings = emailSettings
	}

	res := client.ResourceResponse[client.AlarmNotificationProfileConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   data.Enabled.ValueBool(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}
	
	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateAlarmNotificationProfile(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating alarm notification profile", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	
	// Re-sync state with response
	if updated.Config.Profile.Type == "EmailNotificationProfileType" && updated.Config.Settings != nil {
		settings := updated.Config.Settings
		
		if v, ok := settings["useSmtpProfile"].(bool); ok {
			data.EmailConfig.UseSMTPProfile = types.BoolValue(v)
		}
		// We trust the rest from plan or check similarly
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlarmNotificationProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AlarmNotificationProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlarmNotificationProfile(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting alarm notification profile", err.Error())
		return
	}
}

func (r *AlarmNotificationProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
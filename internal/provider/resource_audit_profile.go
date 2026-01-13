package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AuditProfileResource{}
var _ resource.ResourceWithImportState = &AuditProfileResource{}

func NewAuditProfileResource() resource.Resource {
	return &AuditProfileResource{}
}

// AuditProfileResource defines the resource implementation.
type AuditProfileResource struct {
	client client.IgnitionClient
}

// AuditProfileResourceModel describes the resource data model.
type AuditProfileResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Type                  types.String `tfsdk:"type"`
	RetentionDays         types.Int64  `tfsdk:"retention_days"`
	Database              types.String `tfsdk:"database"`
	PruneEnabled          types.Bool   `tfsdk:"prune_enabled"`
	AutoCreate            types.Bool   `tfsdk:"auto_create"`
	TableName             types.String `tfsdk:"table_name"`
	RemoteServer          types.String `tfsdk:"remote_server"`
	RemoteProfile         types.String `tfsdk:"remote_profile"`
	EnableStoreAndForward types.Bool   `tfsdk:"enable_store_and_forward"`
	Signature             types.String `tfsdk:"signature"`
}

func (r *AuditProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_profile"
}

func (r *AuditProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Audit Profile in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the audit profile.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the audit profile.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the audit profile (database, remote, edge, local).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("database", "remote", "edge", "local"),
				},
			},
			"retention_days": schema.Int64Attribute{
				Description: "How long (in days) should audit records be kept? If set to 0, pruning will be disabled. Default: 90.",
				Optional:    true,
				Computed:    true,
			},
			"database": schema.StringAttribute{
				Description: "The database connection to use to store audit events (for 'database' type).",
				Optional:    true,
			},
			"prune_enabled": schema.BoolAttribute{
				Description: "If false, this audit profile will never prune records (for 'database' type). Default: false.",
				Optional:    true,
				Computed:    true,
			},
			"auto_create": schema.BoolAttribute{
				Description: "If true, the table schema will be automatically created (for 'database' type). Default: true.",
				Optional:    true,
				Computed:    true,
			},
			"table_name": schema.StringAttribute{
				Description: "The name of the table to store audit events (for 'database' type). Default: 'audit_events'.",
				Optional:    true,
				Computed:    true,
			},
			"remote_server": schema.StringAttribute{
				Description: "The remote system to send audit events to (for 'remote' type).",
				Optional:    true,
			},
			"remote_profile": schema.StringAttribute{
				Description: "The audit profile on the remote system to log events into (for 'remote' type).",
				Optional:    true,
			},
			"enable_store_and_forward": schema.BoolAttribute{
				Description: "If enabled, audit events will be stored through the store and forward system (for 'remote' type).",
				Optional:    true,
				Computed:    true,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *AuditProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AuditProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AuditProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := client.AuditProfileProfile{
		Type: data.Type.ValueString(),
	}

	if !data.RetentionDays.IsNull() {
		profile.RetentionDays = int(data.RetentionDays.ValueInt64())
	}

	settings := client.AuditProfileSettings{}
	if !data.Database.IsNull() {
		settings.DatabaseName = data.Database.ValueString()
	}
	if !data.PruneEnabled.IsNull() {
		settings.PruneEnabled = data.PruneEnabled.ValueBool()
	}
	if !data.AutoCreate.IsNull() {
		settings.AutoCreate = data.AutoCreate.ValueBool()
	}
	if !data.TableName.IsNull() {
		settings.TableName = data.TableName.ValueString()
	}
	if !data.RemoteServer.IsNull() {
		settings.RemoteServer = data.RemoteServer.ValueString()
	}
	if !data.RemoteProfile.IsNull() {
		settings.RemoteProfile = data.RemoteProfile.ValueString()
	}
	if !data.EnableStoreAndForward.IsNull() {
		settings.EnableStoreAndForward = data.EnableStoreAndForward.ValueBool()
	}

	config := client.AuditProfileConfig{
		Profile:  profile,
		Settings: settings,
	}

	res := client.ResourceResponse[client.AuditProfileConfig]{
		Name:    data.Name.ValueString(),
		Enabled: true,
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateAuditProfile(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating audit profile", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(data.Name.ValueString())
	if created.Config.Profile.Type != "" {
		data.Type = types.StringValue(created.Config.Profile.Type)
	}
	if !data.RetentionDays.IsNull() {
		profile.RetentionDays = int(data.RetentionDays.ValueInt64())
	}
	data.Description = types.StringValue(created.Description)

	if created.Config.Settings.DatabaseName != "" {
		data.Database = types.StringValue(created.Config.Settings.DatabaseName)
	}
	data.PruneEnabled = types.BoolValue(created.Config.Settings.PruneEnabled)
	data.AutoCreate = types.BoolValue(created.Config.Settings.AutoCreate)
	if created.Config.Settings.TableName != "" {
		data.TableName = types.StringValue(created.Config.Settings.TableName)
	} else if !data.TableName.IsNull() && !data.TableName.IsUnknown() {
		// Keep the plan value if API returns empty but we sent something known
		data.TableName = data.TableName
	} else {
		// Default fallback if not set in plan and not returned
		data.TableName = types.StringValue("audit_events")
	}

	if created.Config.Settings.RemoteServer != "" {
		data.RemoteServer = types.StringValue(created.Config.Settings.RemoteServer)
	}
	if created.Config.Settings.RemoteProfile != "" {
		data.RemoteProfile = types.StringValue(created.Config.Settings.RemoteProfile)
	}
	data.EnableStoreAndForward = types.BoolValue(created.Config.Settings.EnableStoreAndForward)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuditProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AuditProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetAuditProfile(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading audit profile", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(data.Name.ValueString())
	if res.Config.Profile.Type != "" {
		data.Type = types.StringValue(res.Config.Profile.Type)
	}
	data.RetentionDays = types.Int64Value(int64(res.Config.Profile.RetentionDays))
	data.Description = types.StringValue(res.Description)

	if res.Config.Settings.DatabaseName != "" {
		data.Database = types.StringValue(res.Config.Settings.DatabaseName)
	} else {
		data.Database = types.StringNull()
	}
	data.PruneEnabled = types.BoolValue(res.Config.Settings.PruneEnabled)
	data.AutoCreate = types.BoolValue(res.Config.Settings.AutoCreate)
	if res.Config.Settings.TableName != "" {
		data.TableName = types.StringValue(res.Config.Settings.TableName)
	} else {
		// If API doesn't return it, it might be the default or not applicable
		// We shouldn't force it to null if it was set, but for Read we reflect state.
		// If it's a database type, it likely has a default.
		if data.Type.ValueString() == "database" {
			data.TableName = types.StringValue("audit_events")
		} else {
			data.TableName = types.StringNull()
		}
	}
	
	if res.Config.Settings.RemoteServer != "" {
		data.RemoteServer = types.StringValue(res.Config.Settings.RemoteServer)
	} else {
		data.RemoteServer = types.StringNull()
	}
	if res.Config.Settings.RemoteProfile != "" {
		data.RemoteProfile = types.StringValue(res.Config.Settings.RemoteProfile)
	} else {
		data.RemoteProfile = types.StringNull()
	}
	data.EnableStoreAndForward = types.BoolValue(res.Config.Settings.EnableStoreAndForward)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuditProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AuditProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := client.AuditProfileProfile{
		Type: data.Type.ValueString(),
	}
	if !data.RetentionDays.IsNull() {
		profile.RetentionDays = int(data.RetentionDays.ValueInt64())
	}

	settings := client.AuditProfileSettings{}
	if !data.Database.IsNull() {
		settings.DatabaseName = data.Database.ValueString()
	}
	if !data.PruneEnabled.IsNull() {
		settings.PruneEnabled = data.PruneEnabled.ValueBool()
	}
	if !data.AutoCreate.IsNull() {
		settings.AutoCreate = data.AutoCreate.ValueBool()
	}
	if !data.TableName.IsNull() {
		settings.TableName = data.TableName.ValueString()
	}
	if !data.RemoteServer.IsNull() {
		settings.RemoteServer = data.RemoteServer.ValueString()
	}
	if !data.RemoteProfile.IsNull() {
		settings.RemoteProfile = data.RemoteProfile.ValueString()
	}
	if !data.EnableStoreAndForward.IsNull() {
		settings.EnableStoreAndForward = data.EnableStoreAndForward.ValueBool()
	}

	config := client.AuditProfileConfig{
		Profile:  profile,
		Settings: settings,
	}

	res := client.ResourceResponse[client.AuditProfileConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   true,
		Signature: data.Signature.ValueString(),
		Config:    config,
	}
	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateAuditProfile(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating audit profile", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	if updated.Config.Profile.Type != "" {
		data.Type = types.StringValue(updated.Config.Profile.Type)
	}
	data.RetentionDays = types.Int64Value(int64(updated.Config.Profile.RetentionDays))
	data.Description = types.StringValue(updated.Description)

	if updated.Config.Settings.DatabaseName != "" {
		data.Database = types.StringValue(updated.Config.Settings.DatabaseName)
	}
	data.PruneEnabled = types.BoolValue(updated.Config.Settings.PruneEnabled)
	data.AutoCreate = types.BoolValue(updated.Config.Settings.AutoCreate)
	if updated.Config.Settings.TableName != "" {
		data.TableName = types.StringValue(updated.Config.Settings.TableName)
	} else if !data.TableName.IsNull() {
		data.TableName = data.TableName
	} else {
		data.TableName = types.StringValue("audit_events")
	}

	if updated.Config.Settings.RemoteServer != "" {
		data.RemoteServer = types.StringValue(updated.Config.Settings.RemoteServer)
	}
	if updated.Config.Settings.RemoteProfile != "" {
		data.RemoteProfile = types.StringValue(updated.Config.Settings.RemoteProfile)
	}
	data.EnableStoreAndForward = types.BoolValue(updated.Config.Settings.EnableStoreAndForward)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuditProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AuditProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAuditProfile(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting audit profile", err.Error())
		return
	}
}

func (r *AuditProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

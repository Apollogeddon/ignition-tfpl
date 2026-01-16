package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	GenericIgnitionResource[client.AuditProfileConfig, AuditProfileResourceModel]
}

// AuditProfileResourceModel describes the resource data model.
type AuditProfileResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Enabled               types.Bool   `tfsdk:"enabled"`
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
			"enabled": schema.BoolAttribute{
				Description: "Whether the audit profile is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
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

	apiClient, ok := req.ProviderData.(client.IgnitionClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.IgnitionClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.Client = apiClient
	r.Handler = r
	r.Module = "ignition"
	r.ResourceType = "audit-profile"
	r.CreateFunc = apiClient.CreateAuditProfile
	r.GetFunc = apiClient.GetAuditProfile
	r.UpdateFunc = apiClient.UpdateAuditProfile
	r.DeleteFunc = apiClient.DeleteAuditProfile
	
	r.PopulateBase = func(m *AuditProfileResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Enabled = m.Enabled
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.PopulateModel = func(b *BaseResourceModel, m *AuditProfileResourceModel) {
		m.Name = b.Name
		m.Enabled = b.Enabled
		m.Description = b.Description
		m.Signature = b.Signature
		m.Id = b.Id
	}
}

func (r *AuditProfileResource) MapPlanToClient(ctx context.Context, model *AuditProfileResourceModel) (client.AuditProfileConfig, error) {
	profile := client.AuditProfileProfile{
		Type: model.Type.ValueString(),
	}
	if !model.RetentionDays.IsNull() {
		profile.RetentionDays = int(model.RetentionDays.ValueInt64())
	}

	settings := client.AuditProfileSettings{}
	if !model.Database.IsNull() {
		settings.DatabaseName = model.Database.ValueString()
	}
	if !model.PruneEnabled.IsNull() {
		settings.PruneEnabled = model.PruneEnabled.ValueBool()
	}
	if !model.AutoCreate.IsNull() {
		settings.AutoCreate = model.AutoCreate.ValueBool()
	}
	if !model.TableName.IsNull() {
		settings.TableName = model.TableName.ValueString()
	}
	if !model.RemoteServer.IsNull() {
		settings.RemoteServer = model.RemoteServer.ValueString()
	}
	if !model.RemoteProfile.IsNull() {
		settings.RemoteProfile = model.RemoteProfile.ValueString()
	}
	if !model.EnableStoreAndForward.IsNull() {
		settings.EnableStoreAndForward = model.EnableStoreAndForward.ValueBool()
	}

	return client.AuditProfileConfig{
		Profile:  profile,
		Settings: settings,
	}, nil
}

func (r *AuditProfileResource) MapClientToState(ctx context.Context, name string, config *client.AuditProfileConfig, model *AuditProfileResourceModel) error {
	model.Name = types.StringValue(name)
	
	if config.Profile.Type != "" {
		model.Type = types.StringValue(config.Profile.Type)
	}
	
	if config.Profile.RetentionDays != 0 {
		model.RetentionDays = types.Int64Value(int64(config.Profile.RetentionDays))
	}
	
	if config.Settings.DatabaseName != "" {
		model.Database = types.StringValue(config.Settings.DatabaseName)
	}

	model.PruneEnabled = types.BoolValue(config.Settings.PruneEnabled)
	model.AutoCreate = types.BoolValue(config.Settings.AutoCreate)
	
	if config.Settings.TableName != "" {
		model.TableName = types.StringValue(config.Settings.TableName)
	} else if model.TableName.IsNull() || model.TableName.IsUnknown() {
		if model.Type.ValueString() == "database" {
			model.TableName = types.StringValue("audit_events")
		}
	}

	if config.Settings.RemoteServer != "" {
		model.RemoteServer = types.StringValue(config.Settings.RemoteServer)
	}

	if config.Settings.RemoteProfile != "" {
		model.RemoteProfile = types.StringValue(config.Settings.RemoteProfile)
	}

	model.EnableStoreAndForward = types.BoolValue(config.Settings.EnableStoreAndForward)
	
	return nil
}

func (r *AuditProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AuditProfileResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Create(ctx, req, resp, &data, &base)
}

func (r *AuditProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AuditProfileResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Read(ctx, req, resp, &data, &base)
}

func (r *AuditProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AuditProfileResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Update(ctx, req, resp, &data, &base)
}

func (r *AuditProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AuditProfileResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Delete(ctx, req, resp, &data, &base)
}

func (r *AuditProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &AuditProfileResourceModel{
		Id:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	})...)
}
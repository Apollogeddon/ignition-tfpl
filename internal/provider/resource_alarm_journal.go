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
var _ resource.Resource = &AlarmJournalResource{}
var _ resource.ResourceWithImportState = &AlarmJournalResource{}

func NewAlarmJournalResource() resource.Resource {
	return &AlarmJournalResource{}
}

// AlarmJournalResource defines the resource implementation.
type AlarmJournalResource struct {
	client client.IgnitionClient
}

// AlarmJournalResourceModel describes the resource data model.
type AlarmJournalResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Type          types.String `tfsdk:"type"`
	Datasource    types.String `tfsdk:"datasource"`
	TableName     types.String `tfsdk:"table_name"`
	MinPriority   types.String `tfsdk:"min_priority"`
	TargetServer  types.String `tfsdk:"target_server"`
	TargetJournal types.String `tfsdk:"target_journal"`
	Signature     types.String `tfsdk:"signature"`
}

func (r *AlarmJournalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alarm_journal"
}

func (r *AlarmJournalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Alarm Journal in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the alarm journal.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the alarm journal.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the alarm journal is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "The type of the alarm journal (DATASOURCE, REMOTE, LOCAL).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("DATASOURCE", "REMOTE", "LOCAL"),
				},
			},
			"datasource": schema.StringAttribute{
				Description: "The database connection to use (for DATASOURCE type).",
				Optional:    true,
			},
			"table_name": schema.StringAttribute{
				Description: "The table name to store events (for DATASOURCE type).",
				Optional:    true,
				Computed:    true,
			},
			"min_priority": schema.StringAttribute{
				Description: "The minimum priority of events to store (Diagnostic, Low, Medium, High, Critical).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Diagnostic", "Low", "Medium", "High", "Critical"),
				},
			},
			"target_server": schema.StringAttribute{
				Description: "The remote gateway name (for REMOTE type).",
				Optional:    true,
			},
			"target_journal": schema.StringAttribute{
				Description: "The alarm journal on the remote gateway (for REMOTE type).",
				Optional:    true,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *AlarmJournalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AlarmJournalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AlarmJournalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the settings
	settings := client.AlarmJournalSettings{}
	
	if data.Type.ValueString() == "DATASOURCE" {
		if !data.Datasource.IsNull() {
			settings.Datasource = data.Datasource.ValueString()
		}
		
		settings.Advanced = &struct {
			TableName          string `json:"tableName,omitempty"`
			DataTableName      string `json:"dataTableName,omitempty"`
			UseStoreAndForward bool   `json:"useStoreAndForward,omitempty"`
		}{}
		
		if !data.TableName.IsNull() {
			settings.Advanced.TableName = data.TableName.ValueString()
		}
		
		settings.Events = &struct {
			MinPriority            string `json:"minPriority,omitempty"`
			StoreShelvedEvents     bool   `json:"storeShelvedEvents,omitempty"`
			StoreFromEnabledChange bool   `json:"storeFromEnabledChange,omitempty"`
		}{}
		
		if !data.MinPriority.IsNull() {
			settings.Events.MinPriority = data.MinPriority.ValueString()
		}
	} else if data.Type.ValueString() == "REMOTE" {
		settings.RemoteGateway = &struct {
			TargetServer  string `json:"targetServer,omitempty"`
			TargetJournal string `json:"targetJournal,omitempty"`
		}{}
		if !data.TargetServer.IsNull() {
			settings.RemoteGateway.TargetServer = data.TargetServer.ValueString()
		}
		if !data.TargetJournal.IsNull() {
			settings.RemoteGateway.TargetJournal = data.TargetJournal.ValueString()
		}
	}

	config := client.AlarmJournalConfig{
		Profile: client.AlarmJournalProfile{
			Type: data.Type.ValueString(),
		},
		Settings: settings,
	}

	res := client.ResourceResponse[client.AlarmJournalConfig]{
		Name:    data.Name.ValueString(),
		Enabled: data.Enabled.ValueBool(),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateAlarmJournal(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating alarm journal", err.Error())
		return
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
	
	data.Type = types.StringValue(created.Config.Profile.Type)
	
	// Map settings back
	if created.Config.Profile.Type == "DATASOURCE" {
		data.Datasource = types.StringValue(created.Config.Settings.Datasource)
		if created.Config.Settings.Advanced != nil {
			data.TableName = types.StringValue(created.Config.Settings.Advanced.TableName)
		}
		if created.Config.Settings.Events != nil {
			data.MinPriority = types.StringValue(created.Config.Settings.Events.MinPriority)
		}
	} else if created.Config.Profile.Type == "REMOTE" {
		if created.Config.Settings.RemoteGateway != nil {
			data.TargetServer = types.StringValue(created.Config.Settings.RemoteGateway.TargetServer)
			data.TargetJournal = types.StringValue(created.Config.Settings.RemoteGateway.TargetJournal)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlarmJournalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AlarmJournalResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetAlarmJournal(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading alarm journal", err.Error())
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
	
	data.Type = types.StringValue(res.Config.Profile.Type)
	
	// Reset specific fields
	data.Datasource = types.StringNull()
	data.TableName = types.StringNull()
	data.MinPriority = types.StringNull()
	data.TargetServer = types.StringNull()
	data.TargetJournal = types.StringNull()

	if res.Config.Profile.Type == "DATASOURCE" {
		data.Datasource = types.StringValue(res.Config.Settings.Datasource)
		if res.Config.Settings.Advanced != nil {
			data.TableName = types.StringValue(res.Config.Settings.Advanced.TableName)
		}
		if res.Config.Settings.Events != nil {
			data.MinPriority = types.StringValue(res.Config.Settings.Events.MinPriority)
		}
	} else if res.Config.Profile.Type == "REMOTE" {
		if res.Config.Settings.RemoteGateway != nil {
			data.TargetServer = types.StringValue(res.Config.Settings.RemoteGateway.TargetServer)
			data.TargetJournal = types.StringValue(res.Config.Settings.RemoteGateway.TargetJournal)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlarmJournalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AlarmJournalResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the settings
	settings := client.AlarmJournalSettings{}
	
	if data.Type.ValueString() == "DATASOURCE" {
		if !data.Datasource.IsNull() {
			settings.Datasource = data.Datasource.ValueString()
		}
		
		settings.Advanced = &struct {
			TableName          string `json:"tableName,omitempty"`
			DataTableName      string `json:"dataTableName,omitempty"`
			UseStoreAndForward bool   `json:"useStoreAndForward,omitempty"`
		}{}
		
		if !data.TableName.IsNull() {
			settings.Advanced.TableName = data.TableName.ValueString()
		}
		
		settings.Events = &struct {
			MinPriority            string `json:"minPriority,omitempty"`
			StoreShelvedEvents     bool   `json:"storeShelvedEvents,omitempty"`
			StoreFromEnabledChange bool   `json:"storeFromEnabledChange,omitempty"`
		}{}
		
		if !data.MinPriority.IsNull() {
			settings.Events.MinPriority = data.MinPriority.ValueString()
		}
	} else if data.Type.ValueString() == "REMOTE" {
		settings.RemoteGateway = &struct {
			TargetServer  string `json:"targetServer,omitempty"`
			TargetJournal string `json:"targetJournal,omitempty"`
		}{}
		if !data.TargetServer.IsNull() {
			settings.RemoteGateway.TargetServer = data.TargetServer.ValueString()
		}
		if !data.TargetJournal.IsNull() {
			settings.RemoteGateway.TargetJournal = data.TargetJournal.ValueString()
		}
	}

	config := client.AlarmJournalConfig{
		Profile: client.AlarmJournalProfile{
			Type: data.Type.ValueString(),
		},
		Settings: settings,
	}

	res := client.ResourceResponse[client.AlarmJournalConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   data.Enabled.ValueBool(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateAlarmJournal(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating alarm journal", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	data.Type = types.StringValue(updated.Config.Profile.Type)
	
	if updated.Description != "" {
		data.Description = types.StringValue(updated.Description)
	}

	if updated.Config.Profile.Type == "DATASOURCE" {
		data.Datasource = types.StringValue(updated.Config.Settings.Datasource)
		if updated.Config.Settings.Advanced != nil {
			data.TableName = types.StringValue(updated.Config.Settings.Advanced.TableName)
		}
		if updated.Config.Settings.Events != nil {
			data.MinPriority = types.StringValue(updated.Config.Settings.Events.MinPriority)
		}
	} else if updated.Config.Profile.Type == "REMOTE" {
		if updated.Config.Settings.RemoteGateway != nil {
			data.TargetServer = types.StringValue(updated.Config.Settings.RemoteGateway.TargetServer)
			data.TargetJournal = types.StringValue(updated.Config.Settings.RemoteGateway.TargetJournal)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AlarmJournalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AlarmJournalResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlarmJournal(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting alarm journal", err.Error())
		return
	}
}

func (r *AlarmJournalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

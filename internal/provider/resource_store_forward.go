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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StoreAndForwardResource{}
var _ resource.ResourceWithImportState = &StoreAndForwardResource{}

func NewStoreAndForwardResource() resource.Resource {
	return &StoreAndForwardResource{}
}

// StoreAndForwardResource defines the resource implementation.
type StoreAndForwardResource struct {
	client client.IgnitionClient
}

// StoreAndForwardResourceModel describes the resource data model.
type StoreAndForwardResourceModel struct {
	Id                 types.String            `tfsdk:"id"`
	Name               types.String            `tfsdk:"name"`
	Description        types.String            `tfsdk:"description"`
	Enabled            types.Bool              `tfsdk:"enabled"`
	TimeThresholdMs    types.Int64             `tfsdk:"time_threshold_ms"`
	ForwardRateMs      types.Int64             `tfsdk:"forward_rate_ms"`
	ForwardingPolicy   types.String            `tfsdk:"forwarding_policy"`
	ForwardingSchedule types.String            `tfsdk:"forwarding_schedule"`
	IsThirdParty       types.Bool              `tfsdk:"is_third_party"`
	DataThreshold      types.Int64             `tfsdk:"data_threshold"`
	BatchSize          types.Int64             `tfsdk:"batch_size"`
	ScanRateMs         types.Int64             `tfsdk:"scan_rate_ms"`
	PrimaryPolicy      *MaintenancePolicyModel `tfsdk:"primary_policy"`
	SecondaryPolicy    *MaintenancePolicyModel `tfsdk:"secondary_policy"`
	Signature          types.String            `tfsdk:"signature"`
}

type MaintenancePolicyModel struct {
	Action    types.String `tfsdk:"action"`
	LimitType types.String `tfsdk:"limit_type"`
	Value     types.Int64  `tfsdk:"value"`
}

func (r *StoreAndForwardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_store_forward"
}

func (r *StoreAndForwardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	maintenancePolicySchema := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"action": schema.StringAttribute{
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("EVICT_OLDEST_DATA", "PREVENT_NEW_DATA"),
				},
			},
			"limit_type": schema.StringAttribute{
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("TIME_DURATION", "COUNT", "FILE_SIZE"),
				},
			},
			"value": schema.Int64Attribute{
				Required: true,
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages a Store and Forward Engine in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the engine.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the engine.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the engine is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"time_threshold_ms": schema.Int64Attribute{
				Description: "Maximum time data accumulates in memory before writing to disk.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1000),
			},
			"forward_rate_ms": schema.Int64Attribute{
				Description: "Rate at which data is eligible for forwarding.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1000),
			},
			"forwarding_policy": schema.StringAttribute{
				Description: "Policy for forwarding (ALL, PRIMARY_ONLY, SECONDARY_ONLY).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ALL"),
				Validators: []validator.String{
					stringvalidator.OneOf("ALL", "PRIMARY_ONLY", "SECONDARY_ONLY"),
				},
			},
			"forwarding_schedule": schema.StringAttribute{
				Description: "Comma separated list of time ranges (e.g. 9:00-15:00).",
				Optional:    true,
			},
			"is_third_party": schema.BoolAttribute{
				Description: "Whether the engine is managed by a third party.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"data_threshold": schema.Int64Attribute{
				Description: "Maximum data amount in memory before writing to disk.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(100),
			},
			"batch_size": schema.Int64Attribute{
				Description: "Batch size for forwarding.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(100),
			},
			"scan_rate_ms": schema.Int64Attribute{
				Description: "Rate at which data pipelines are scanned.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1000),
			},
			"primary_policy": schema.SingleNestedAttribute{
				Description: "Maintenance policy for primary datastore (memory).",
				Optional:    true,
				Attributes:  maintenancePolicySchema.Attributes,
			},
			"secondary_policy": schema.SingleNestedAttribute{
				Description: "Maintenance policy for secondary datastore (disk).",
				Optional:    true,
				Attributes:  maintenancePolicySchema.Attributes,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource.",
				Computed:    true,
			},
		},
	}
}

func (r *StoreAndForwardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StoreAndForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StoreAndForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.StoreAndForwardConfig{
		TimeThresholdMs:    int(data.TimeThresholdMs.ValueInt64()),
		ForwardRateMs:      int(data.ForwardRateMs.ValueInt64()),
		ForwardingPolicy:   data.ForwardingPolicy.ValueString(),
		ForwardingSchedule: data.ForwardingSchedule.ValueString(),
		IsThirdParty:       data.IsThirdParty.ValueBool(),
		DataThreshold:      int(data.DataThreshold.ValueInt64()),
		BatchSize:          int(data.BatchSize.ValueInt64()),
		ScanRateMs:         int(data.ScanRateMs.ValueInt64()),
	}

	if data.PrimaryPolicy != nil {
		config.PrimaryMaintenancePolicy = &client.StoreAndForwardMaintenancePolicy{
			Action: data.PrimaryPolicy.Action.ValueString(),
		}
		config.PrimaryMaintenancePolicy.Limit.LimitType = data.PrimaryPolicy.LimitType.ValueString()
		config.PrimaryMaintenancePolicy.Limit.Value = int(data.PrimaryPolicy.Value.ValueInt64())
	}

	if data.SecondaryPolicy != nil {
		config.SecondaryMaintenancePolicy = &client.StoreAndForwardMaintenancePolicy{
			Action: data.SecondaryPolicy.Action.ValueString(),
		}
		config.SecondaryMaintenancePolicy.Limit.LimitType = data.SecondaryPolicy.LimitType.ValueString()
		config.SecondaryMaintenancePolicy.Limit.Value = int(data.SecondaryPolicy.Value.ValueInt64())
	}

	res := client.ResourceResponse[client.StoreAndForwardConfig]{
		Name:    data.Name.ValueString(),
		Enabled: data.Enabled.ValueBool(),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateStoreAndForward(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Store and Forward engine", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StoreAndForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StoreAndForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetStoreAndForward(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Store and Forward engine", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Enabled = types.BoolValue(res.Enabled)
	data.Description = types.StringValue(res.Description)
	
	data.TimeThresholdMs = types.Int64Value(int64(res.Config.TimeThresholdMs))
	data.ForwardRateMs = types.Int64Value(int64(res.Config.ForwardRateMs))
	data.ForwardingPolicy = types.StringValue(res.Config.ForwardingPolicy)
	data.ForwardingSchedule = types.StringValue(res.Config.ForwardingSchedule)
	data.IsThirdParty = types.BoolValue(res.Config.IsThirdParty)
	data.DataThreshold = types.Int64Value(int64(res.Config.DataThreshold))
	data.BatchSize = types.Int64Value(int64(res.Config.BatchSize))
	data.ScanRateMs = types.Int64Value(int64(res.Config.ScanRateMs))

	if res.Config.PrimaryMaintenancePolicy != nil {
		data.PrimaryPolicy = &MaintenancePolicyModel{
			Action:    types.StringValue(res.Config.PrimaryMaintenancePolicy.Action),
			LimitType: types.StringValue(res.Config.PrimaryMaintenancePolicy.Limit.LimitType),
			Value:     types.Int64Value(int64(res.Config.PrimaryMaintenancePolicy.Limit.Value)),
		}
	}

	if res.Config.SecondaryMaintenancePolicy != nil {
		data.SecondaryPolicy = &MaintenancePolicyModel{
			Action:    types.StringValue(res.Config.SecondaryMaintenancePolicy.Action),
			LimitType: types.StringValue(res.Config.SecondaryMaintenancePolicy.Limit.LimitType),
			Value:     types.Int64Value(int64(res.Config.SecondaryMaintenancePolicy.Limit.Value)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StoreAndForwardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StoreAndForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.StoreAndForwardConfig{
		TimeThresholdMs:    int(data.TimeThresholdMs.ValueInt64()),
		ForwardRateMs:      int(data.ForwardRateMs.ValueInt64()),
		ForwardingPolicy:   data.ForwardingPolicy.ValueString(),
		ForwardingSchedule: data.ForwardingSchedule.ValueString(),
		IsThirdParty:       data.IsThirdParty.ValueBool(),
		DataThreshold:      int(data.DataThreshold.ValueInt64()),
		BatchSize:          int(data.BatchSize.ValueInt64()),
		ScanRateMs:         int(data.ScanRateMs.ValueInt64()),
	}

	if data.PrimaryPolicy != nil {
		config.PrimaryMaintenancePolicy = &client.StoreAndForwardMaintenancePolicy{
			Action: data.PrimaryPolicy.Action.ValueString(),
		}
		config.PrimaryMaintenancePolicy.Limit.LimitType = data.PrimaryPolicy.LimitType.ValueString()
		config.PrimaryMaintenancePolicy.Limit.Value = int(data.PrimaryPolicy.Value.ValueInt64())
	}

	if data.SecondaryPolicy != nil {
		config.SecondaryMaintenancePolicy = &client.StoreAndForwardMaintenancePolicy{
			Action: data.SecondaryPolicy.Action.ValueString(),
		}
		config.SecondaryMaintenancePolicy.Limit.LimitType = data.SecondaryPolicy.LimitType.ValueString()
		config.SecondaryMaintenancePolicy.Limit.Value = int(data.SecondaryPolicy.Value.ValueInt64())
	}

	res := client.ResourceResponse[client.StoreAndForwardConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   data.Enabled.ValueBool(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateStoreAndForward(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Store and Forward engine", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StoreAndForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StoreAndForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStoreAndForward(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Store and Forward engine", err.Error())
		return
	}
}

func (r *StoreAndForwardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

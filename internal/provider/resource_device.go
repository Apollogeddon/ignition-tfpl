package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DeviceResource{}
var _ resource.ResourceWithImportState = &DeviceResource{}

func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

type DeviceResource struct {
	Res GenericIgnitionResource[client.DeviceConfig, DeviceResourceModel]
}

type DeviceResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Type        types.String `tfsdk:"type"`
	Parameters  types.String `tfsdk:"parameters"`
	Signature   types.String `tfsdk:"signature"`
}

func (r *DeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *DeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Device (Driver) in Ignition (e.g., Modbus, Siemens, Simulator).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the device.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the device.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the device is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "The type of the device (e.g., 'ModbusTcp', 'S71500', 'ProgrammableSimulatorDevice').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parameters": schema.StringAttribute{
				Description: "The JSON configuration parameters for the device. These vary by device type.",
				Required:    true,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *DeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Res.Client = client
	r.Res.Handler = r
	r.Res.Module = "com.inductiveautomation.opcua"
	r.Res.ResourceType = "device"
	r.Res.CreateFunc = client.CreateDevice
	r.Res.GetFunc = client.GetDevice
	r.Res.UpdateFunc = client.UpdateDevice
	r.Res.DeleteFunc = client.DeleteDevice
	r.Res.PopulateBase = func(m *DeviceResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Enabled = m.Enabled
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.Res.PopulateModel = func(b *BaseResourceModel, m *DeviceResourceModel) {
		m.Name = b.Name
		m.Enabled = b.Enabled
		m.Description = b.Description
		m.Signature = b.Signature
		m.Id = b.Id
	}
}

func (r *DeviceResource) MapPlanToClient(ctx context.Context, model *DeviceResourceModel) (client.DeviceConfig, error) {
	config := make(client.DeviceConfig)
	
	// Unmarshal JSON parameters
	if err := json.Unmarshal([]byte(model.Parameters.ValueString()), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters JSON: %w", err)
	}

	return config, nil
}

func (r *DeviceResource) MapClientToState(ctx context.Context, name string, config *client.DeviceConfig, model *DeviceResourceModel) error {
	model.Name = types.StringValue(name)
	
	// Marshal the API response to canonical JSON
	b, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal parameters to JSON: %w", err)
	}
	canonicalAPI := string(b)

	// If the model already has a value (from Plan), check if it's semantically equal
	if !model.Parameters.IsNull() && !model.Parameters.IsUnknown() {
		var modelConfig client.DeviceConfig
		// Try to unmarshal the plan value
		if err := json.Unmarshal([]byte(model.Parameters.ValueString()), &modelConfig); err == nil {
			// Re-marshal to compare canonical forms
			bModel, _ := json.Marshal(modelConfig)
			if string(bModel) == canonicalAPI {
				// They are semantically equal (ignoring whitespace), so keep the Plan value
				// to satisfy Terraform's consistency check.
				return nil
			}
		}
	}

	// If they differ (or model is null/unknown), accept the API's canonical value
	model.Parameters = types.StringValue(canonicalAPI)

	return nil
}

func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceResourceModel
	var base BaseResourceModel
	
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Res.PopulateBase != nil {
		r.Res.PopulateBase(&data, &base)
	}

	config, err := r.MapPlanToClient(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping plan to client", err.Error())
		return
	}

	res := client.ResourceResponse[client.DeviceConfig]{
		Module:       r.Res.Module,
		Type:         data.Type.ValueString(), // Use the driver type from the plan
		Name:         base.Name.ValueString(),
		Enabled:      boolPtr(base.Enabled.ValueBool()),
		Config:       config,
	}

	if !base.Description.IsNull() {
		res.Description = base.Description.ValueString()
	}

	created, err := r.Res.CreateFunc(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}

	base.Signature = types.StringValue(created.Signature)
	base.Id = types.StringValue(created.Name)
	if base.Name.IsNull() || base.Name.IsUnknown() || base.Name.ValueString() == "" {
		base.Name = types.StringValue(created.Name)
	}
	if created.Enabled != nil {
		base.Enabled = types.BoolValue(*created.Enabled)
	} else {
		base.Enabled = types.BoolValue(true)
	}
	if created.Description != "" {
		base.Description = types.StringValue(created.Description)
	} else if base.Description.IsNull() || base.Description.IsUnknown() {
		base.Description = types.StringNull()
	}
	
	// The API returns the Driver Type in the Type field, so we map it back
	data.Type = types.StringValue(created.Type)

	if r.Res.PopulateModel != nil {
		r.Res.PopulateModel(&base, &data)
	}

	if err := r.MapClientToState(ctx, created.Name, &created.Config, &data); err != nil {
		resp.Diagnostics.AddError("Error mapping client to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceResourceModel
	var base BaseResourceModel
	r.Res.Read(ctx, req, resp, &data, &base)
}

func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceResourceModel
	var base BaseResourceModel
	
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve existing signature from state
	var stateModel DeviceResourceModel
	var stateBase BaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if !resp.Diagnostics.HasError() && r.Res.PopulateBase != nil {
		r.Res.PopulateBase(&stateModel, &stateBase)
	}

	if r.Res.PopulateBase != nil {
		r.Res.PopulateBase(&data, &base)
	}

	config, err := r.MapPlanToClient(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping plan to client", err.Error())
		return
	}

	res := client.ResourceResponse[client.DeviceConfig]{
		Module:       r.Res.Module,
		Type:         data.Type.ValueString(), // Use the driver type from the plan
		Name:         base.Name.ValueString(),
		Enabled:      boolPtr(base.Enabled.ValueBool()),
		Signature:    stateBase.Signature.ValueString(),
		Config:       config,
	}

	if !base.Description.IsNull() {
		res.Description = base.Description.ValueString()
	}

	updated, err := r.Res.UpdateFunc(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	if updated.Signature != "" {
		base.Signature = types.StringValue(updated.Signature)
	} else {
		fresh, err := r.Res.GetFunc(ctx, base.Name.ValueString())
		if err == nil && fresh.Signature != "" {
			base.Signature = types.StringValue(fresh.Signature)
			updated = fresh
		} else if !stateBase.Signature.IsNull() && !stateBase.Signature.IsUnknown() {
			base.Signature = stateBase.Signature
			resp.Diagnostics.AddWarning("Missing Signature on Update", 
				"The API returned an empty signature after update and refresh failed. Preserving the existing signature.")
		}
	}

	base.Id = types.StringValue(updated.Name)
	if base.Name.IsNull() || base.Name.IsUnknown() || base.Name.ValueString() == "" {
		base.Name = types.StringValue(updated.Name)
	}

	if updated.Enabled != nil {
		base.Enabled = types.BoolValue(*updated.Enabled)
	} else {
		base.Enabled = types.BoolValue(true)
	}
	if updated.Description != "" {
		base.Description = types.StringValue(updated.Description)
	} else if base.Description.IsNull() || base.Description.IsUnknown() {
		base.Description = types.StringNull()
	}
	
	// Map the Type back
	data.Type = types.StringValue(updated.Type)
	
	if r.Res.PopulateModel != nil {
		r.Res.PopulateModel(&base, &data)
	}

	if err := r.MapClientToState(ctx, updated.Name, &updated.Config, &data); err != nil {
		resp.Diagnostics.AddError("Error mapping client to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceResourceModel
	var base BaseResourceModel
	r.Res.Delete(ctx, req, resp, &data, &base)
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &DeviceResourceModel{
		Id:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	})...)
}
package provider

import (
	"context"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BaseResourceModel includes the common fields for Ignition resources
type BaseResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Description types.String `tfsdk:"description"`
	Signature   types.String `tfsdk:"signature"`
}

// IgnitionResourceHandler defines the unique mapping logic for a specific resource
type IgnitionResourceHandler[T any, M any] interface {
	// MapPlanToClient converts the Terraform model (plus plan data) to the API config struct
	MapPlanToClient(ctx context.Context, model *M) (T, error)
	// MapClientToState updates the Terraform model from the API response
	MapClientToState(ctx context.Context, name string, config *T, model *M) error
}

// GenericIgnitionResource implements the core CRUD logic
type GenericIgnitionResource[T any, M any] struct {
	Client       client.IgnitionClient
	Handler      IgnitionResourceHandler[T, M]
	Module       string
	ResourceType string
	
	// API Methods
	CreateFunc func(context.Context, client.ResourceResponse[T]) (*client.ResourceResponse[T], error)
	GetFunc    func(context.Context, string) (*client.ResourceResponse[T], error)
	UpdateFunc func(context.Context, client.ResourceResponse[T]) (*client.ResourceResponse[T], error)
	DeleteFunc func(context.Context, string, string) error

	// Helper to extract base fields from the model
	PopulateBase func(*M, *BaseResourceModel)
	// Helper to update the model from base fields
	PopulateModel func(*BaseResourceModel, *M)
}

func (r *GenericIgnitionResource[T, M]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, data *M, base *BaseResourceModel) {
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.PopulateBase != nil {
		r.PopulateBase(data, base)
	}

	config, err := r.Handler.MapPlanToClient(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping plan to client", err.Error())
		return
	}

	res := client.ResourceResponse[T]{
		Module:       r.Module,
		Type:        r.ResourceType,
		Name:         base.Name.ValueString(),
		Enabled:      boolPtr(base.Enabled.ValueBool()),
		Config:       config,
	}

	if !base.Description.IsNull() {
		res.Description = base.Description.ValueString()
	}

	created, err := r.CreateFunc(ctx, res)
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
	
	if r.PopulateModel != nil {
		r.PopulateModel(base, data)
	}

	if err := r.Handler.MapClientToState(ctx, created.Name, &created.Config, data); err != nil {
		resp.Diagnostics.AddError("Error mapping client to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *GenericIgnitionResource[T, M]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, data *M, base *BaseResourceModel) {
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.PopulateBase != nil {
		r.PopulateBase(data, base)
	}

	if base.Name.ValueString() == "" {
		return
	}

	res, err := r.GetFunc(ctx, base.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	base.Signature = types.StringValue(res.Signature)
	base.Id = types.StringValue(res.Name)
	if base.Name.IsNull() || base.Name.IsUnknown() || base.Name.ValueString() == "" {
		base.Name = types.StringValue(res.Name)
	}
	if res.Enabled != nil {
		base.Enabled = types.BoolValue(*res.Enabled)
	} else {
		base.Enabled = types.BoolValue(true)
	}
	
	if res.Description != "" {
		base.Description = types.StringValue(res.Description)
	} else if base.Description.IsNull() || base.Description.IsUnknown() {
		base.Description = types.StringNull()
	}

	if r.PopulateModel != nil {
		r.PopulateModel(base, data)
	}

	if err := r.Handler.MapClientToState(ctx, res.Name, &res.Config, data); err != nil {
		resp.Diagnostics.AddError("Error mapping client to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *GenericIgnitionResource[T, M]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, data *M, base *BaseResourceModel) {
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve existing signature from state
	var stateModel M
	var stateBase BaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if !resp.Diagnostics.HasError() && r.PopulateBase != nil {
		r.PopulateBase(&stateModel, &stateBase)
	}

	if r.PopulateBase != nil {
		r.PopulateBase(data, base)
	}

	config, err := r.Handler.MapPlanToClient(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping plan to client", err.Error())
		return
	}

	res := client.ResourceResponse[T]{
		Module:       r.Module,
		Type:        r.ResourceType,
		Name:         base.Name.ValueString(),
		Enabled:      boolPtr(base.Enabled.ValueBool()),
		Signature:    stateBase.Signature.ValueString(),
		Config:       config,
	}

	if !base.Description.IsNull() {
		res.Description = base.Description.ValueString()
	}

	updated, err := r.UpdateFunc(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	if updated.Signature != "" {
		base.Signature = types.StringValue(updated.Signature)
	} else {
		// Attempt to fetch the latest signature if not returned by Update
		fresh, err := r.GetFunc(ctx, base.Name.ValueString())
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
	
	if r.PopulateModel != nil {
		r.PopulateModel(base, data)
	}

	if err := r.Handler.MapClientToState(ctx, updated.Name, &updated.Config, data); err != nil {
		resp.Diagnostics.AddError("Error mapping client to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *GenericIgnitionResource[T, M]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, data *M, base *BaseResourceModel) {
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.PopulateBase != nil {
		r.PopulateBase(data, base)
	}

	err := r.DeleteFunc(ctx, base.Name.ValueString(), base.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}
}

// stringToNullableString returns a types.StringNull if the input is empty,
// otherwise returns types.StringValue.
func stringToNullableString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func boolPtr(b bool) *bool {
	return &b
}

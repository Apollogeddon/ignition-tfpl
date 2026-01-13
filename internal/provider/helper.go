package provider

import (
	"context"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
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
	MapClientToState(ctx context.Context, config *T, model *M) error
}

// GenericIgnitionResource implements the core CRUD logic
type GenericIgnitionResource[T any, M any] struct {
	Client  client.IgnitionClient
	Handler IgnitionResourceHandler[T, M]
	
	// API Methods
	CreateFunc func(context.Context, client.ResourceResponse[T]) (*client.ResourceResponse[T], error)
	GetFunc    func(context.Context, string) (*client.ResourceResponse[T], error)
	UpdateFunc func(context.Context, client.ResourceResponse[T]) (*client.ResourceResponse[T], error)
	DeleteFunc func(context.Context, string, string) error
}

func (r *GenericIgnitionResource[T, M]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, data *M, base *BaseResourceModel) {
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.Handler.MapPlanToClient(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping plan to client", err.Error())
		return
	}

	res := client.ResourceResponse[T]{
		Name:    base.Name.ValueString(),
		Enabled: base.Enabled.ValueBool(),
		Config:  config,
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
	
	if err := r.Handler.MapClientToState(ctx, &created.Config, data); err != nil {
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

	res, err := r.GetFunc(ctx, base.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	base.Signature = types.StringValue(res.Signature)
	base.Id = types.StringValue(res.Name)
	base.Enabled = types.BoolValue(res.Enabled)
	
	if res.Description != "" {
		base.Description = types.StringValue(res.Description)
	} else {
		base.Description = types.StringNull()
	}

	if err := r.Handler.MapClientToState(ctx, &res.Config, data); err != nil {
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

	config, err := r.Handler.MapPlanToClient(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping plan to client", err.Error())
		return
	}

	res := client.ResourceResponse[T]{
		Name:      base.Name.ValueString(),
		Enabled:   base.Enabled.ValueBool(),
		Signature: base.Signature.ValueString(),
		Config:    config,
	}

	if !base.Description.IsNull() {
		res.Description = base.Description.ValueString()
	}

	updated, err := r.UpdateFunc(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	base.Signature = types.StringValue(updated.Signature)
	
	if err := r.Handler.MapClientToState(ctx, &updated.Config, data); err != nil {
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

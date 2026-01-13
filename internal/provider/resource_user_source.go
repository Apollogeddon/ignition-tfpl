package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserSourceResource{}
var _ resource.ResourceWithImportState = &UserSourceResource{}

func NewUserSourceResource() resource.Resource {
	return &UserSourceResource{}
}

// UserSourceResource defines the resource implementation.
type UserSourceResource struct {
	client client.IgnitionClient
}

// UserSourceResourceModel describes the resource data model.
type UserSourceResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	Description        types.String `tfsdk:"description"`
	FailoverProfile    types.String `tfsdk:"failover_profile"`
	FailoverMode       types.String `tfsdk:"failover_mode"`
	ScheduleRestricted types.Bool   `tfsdk:"schedule_restricted"`
	Signature          types.String `tfsdk:"signature"`
}

func (r *UserSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_source"
}

func (r *UserSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a User Source in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the user source.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the user source (e.g., INTERNAL, ADEASY, ADHYBRID, AD_DB_HYBRID, DATASOURCE).",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the user source.",
				Optional:    true,
			},
			"failover_profile": schema.StringAttribute{
				Description: "If this source is unreachable for authentication, this failover source will be used instead.",
				Optional:    true,
			},
			"failover_mode": schema.StringAttribute{
				Description: "The failover mode to use if a failover source is set. Hard - failover only if this source is unreachable. Soft - try the failover source when a user fails to authenticate with this source.",
				Optional:    true,
			},
			"schedule_restricted": schema.BoolAttribute{
				Description: "Users are only able to log in when their assigned schedule is active.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *UserSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := client.UserSourceProfile{
		Type: data.Type.ValueString(),
	}

	if !data.FailoverProfile.IsNull() {
		profile.FailoverProfile = data.FailoverProfile.ValueString()
	}
	if !data.FailoverMode.IsNull() {
		profile.FailoverMode = data.FailoverMode.ValueString()
	}
	if !data.ScheduleRestricted.IsNull() {
		profile.ScheduleRestricted = data.ScheduleRestricted.ValueBool()
	}

	config := client.UserSourceConfig{
		Profile: profile,
	}

	res := client.ResourceResponse[client.UserSourceConfig]{
		Name:    data.Name.ValueString(),
		Enabled: true,
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateUserSource(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user source", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(data.Name.ValueString())
	if created.Config.Profile.Type != "" {
		data.Type = types.StringValue(created.Config.Profile.Type)
	}
	data.Description = stringToNullableString(created.Description)

	if created.Config.Profile.FailoverProfile != "" {
		data.FailoverProfile = types.StringValue(created.Config.Profile.FailoverProfile)
	}
	if created.Config.Profile.FailoverMode != "" {
		data.FailoverMode = types.StringValue(created.Config.Profile.FailoverMode)
	}
	data.ScheduleRestricted = types.BoolValue(created.Config.Profile.ScheduleRestricted)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetUserSource(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user source", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(data.Name.ValueString())
	if res.Config.Profile.Type != "" {
		data.Type = types.StringValue(res.Config.Profile.Type)
	}
	data.Description = stringToNullableString(res.Description)

	if res.Config.Profile.FailoverProfile != "" {
		data.FailoverProfile = types.StringValue(res.Config.Profile.FailoverProfile)
	} else {
		data.FailoverProfile = types.StringNull()
	}

	if res.Config.Profile.FailoverMode != "" {
		data.FailoverMode = types.StringValue(res.Config.Profile.FailoverMode)
	} else {
		data.FailoverMode = types.StringNull()
	}
	
	data.ScheduleRestricted = types.BoolValue(res.Config.Profile.ScheduleRestricted)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := client.UserSourceProfile{
		Type: data.Type.ValueString(),
	}

	if !data.FailoverProfile.IsNull() {
		profile.FailoverProfile = data.FailoverProfile.ValueString()
	}
	if !data.FailoverMode.IsNull() {
		profile.FailoverMode = data.FailoverMode.ValueString()
	}
	if !data.ScheduleRestricted.IsNull() {
		profile.ScheduleRestricted = data.ScheduleRestricted.ValueBool()
	}

	config := client.UserSourceConfig{
		Profile: profile,
	}

	res := client.ResourceResponse[client.UserSourceConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   true,
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateUserSource(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user source", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	if updated.Config.Profile.Type != "" {
		data.Type = types.StringValue(updated.Config.Profile.Type)
	}
	data.Description = stringToNullableString(updated.Description)

	if updated.Config.Profile.FailoverProfile != "" {
		data.FailoverProfile = types.StringValue(updated.Config.Profile.FailoverProfile)
	}
	if updated.Config.Profile.FailoverMode != "" {
		data.FailoverMode = types.StringValue(updated.Config.Profile.FailoverMode)
	}
	data.ScheduleRestricted = types.BoolValue(updated.Config.Profile.ScheduleRestricted)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUserSource(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user source", err.Error())
		return
	}
}

func (r *UserSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

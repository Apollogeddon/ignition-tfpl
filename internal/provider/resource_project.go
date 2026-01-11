package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client client.IgnitionClient
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Title            types.String `tfsdk:"title"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Parent           types.String `tfsdk:"parent"`
	Inheritable      types.Bool   `tfsdk:"inheritable"`
	DefaultDB        types.String `tfsdk:"default_db"`
	TagProvider      types.String `tfsdk:"tag_provider"`
	UserSource       types.String `tfsdk:"user_source"`
	IdentityProvider types.String `tfsdk:"identity_provider"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ignition Project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the project.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the project.",
				Optional:    true,
			},
			"title": schema.StringAttribute{
				Description: "The title of the project.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the project is enabled.",
				Optional:    true,
				Computed:    true,
			},
			"parent": schema.StringAttribute{
				Description: "The parent project name, if any.",
				Optional:    true,
			},
			"inheritable": schema.BoolAttribute{
				Description: "Whether the project is inheritable.",
				Optional:    true,
				Computed:    true,
			},
			"default_db": schema.StringAttribute{
				Description: "The default database connection for the project.",
				Optional:    true,
			},
			"tag_provider": schema.StringAttribute{
				Description: "The default tag provider for the project.",
				Optional:    true,
			},
			"user_source": schema.StringAttribute{
				Description: "The default user source for the project.",
				Optional:    true,
			},
			"identity_provider": schema.StringAttribute{
				Description: "The default identity provider for the project.",
				Optional:    true,
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := client.Project{
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		project.Description = data.Description.ValueString()
	}
	if !data.Title.IsNull() {
		project.Title = data.Title.ValueString()
	}
	if !data.Enabled.IsNull() {
		project.Enabled = data.Enabled.ValueBool()
	}
	if !data.Parent.IsNull() {
		project.Parent = data.Parent.ValueString()
	}
	if !data.Inheritable.IsNull() {
		project.Inheritable = data.Inheritable.ValueBool()
	}
	if !data.DefaultDB.IsNull() {
		project.DefaultDB = data.DefaultDB.ValueString()
	}
	if !data.TagProvider.IsNull() {
		project.TagProvider = data.TagProvider.ValueString()
	}
	if !data.UserSource.IsNull() {
		project.UserSource = data.UserSource.ValueString()
	}
	if !data.IdentityProvider.IsNull() {
		project.IdentityProvider = data.IdentityProvider.ValueString()
	}

	created, err := r.client.CreateProject(ctx, project)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	data.Description = types.StringValue(created.Description)
	data.Id = types.StringValue(data.Name.ValueString())
	if created.Title != "" {
		data.Title = types.StringValue(created.Title)
	}
	data.Enabled = types.BoolValue(created.Enabled)
	if created.Parent != "" {
		data.Parent = types.StringValue(created.Parent)
	}
	data.Inheritable = types.BoolValue(created.Inheritable)
	if created.DefaultDB != "" {
		data.DefaultDB = types.StringValue(created.DefaultDB)
	}
	if created.TagProvider != "" {
		data.TagProvider = types.StringValue(created.TagProvider)
	}
	if created.UserSource != "" {
		data.UserSource = types.StringValue(created.UserSource)
	}
	if created.IdentityProvider != "" {
		data.IdentityProvider = types.StringValue(created.IdentityProvider)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetProject(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	data.Description = types.StringValue(res.Description)
	data.Id = types.StringValue(data.Name.ValueString())
	if res.Title != "" {
		data.Title = types.StringValue(res.Title)
	}
	data.Enabled = types.BoolValue(res.Enabled)
	if res.Parent != "" {
		data.Parent = types.StringValue(res.Parent)
	}
	data.Inheritable = types.BoolValue(res.Inheritable)
	if res.DefaultDB != "" {
		data.DefaultDB = types.StringValue(res.DefaultDB)
	}
	if res.TagProvider != "" {
		data.TagProvider = types.StringValue(res.TagProvider)
	}
	if res.UserSource != "" {
		data.UserSource = types.StringValue(res.UserSource)
	}
	if res.IdentityProvider != "" {
		data.IdentityProvider = types.StringValue(res.IdentityProvider)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := client.Project{
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		project.Description = data.Description.ValueString()
	}
	if !data.Title.IsNull() {
		project.Title = data.Title.ValueString()
	}
	if !data.Enabled.IsNull() {
		project.Enabled = data.Enabled.ValueBool()
	}
	if !data.Parent.IsNull() {
		project.Parent = data.Parent.ValueString()
	}
	if !data.Inheritable.IsNull() {
		project.Inheritable = data.Inheritable.ValueBool()
	}
	if !data.DefaultDB.IsNull() {
		project.DefaultDB = data.DefaultDB.ValueString()
	}
	if !data.TagProvider.IsNull() {
		project.TagProvider = data.TagProvider.ValueString()
	}
	if !data.UserSource.IsNull() {
		project.UserSource = data.UserSource.ValueString()
	}
	if !data.IdentityProvider.IsNull() {
		project.IdentityProvider = data.IdentityProvider.ValueString()
	}

	updated, err := r.client.UpdateProject(ctx, project)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	data.Description = types.StringValue(updated.Description)
	if updated.Title != "" {
		data.Title = types.StringValue(updated.Title)
	}
	data.Enabled = types.BoolValue(updated.Enabled)
	if updated.Parent != "" {
		data.Parent = types.StringValue(updated.Parent)
	}
	data.Inheritable = types.BoolValue(updated.Inheritable)
	if updated.DefaultDB != "" {
		data.DefaultDB = types.StringValue(updated.DefaultDB)
	}
	if updated.TagProvider != "" {
		data.TagProvider = types.StringValue(updated.TagProvider)
	}
	if updated.UserSource != "" {
		data.UserSource = types.StringValue(updated.UserSource)
	}
	if updated.IdentityProvider != "" {
		data.IdentityProvider = types.StringValue(updated.IdentityProvider)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProject(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

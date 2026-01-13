package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
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
	GenericIgnitionResource[client.Project, ProjectResourceModel]
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Title            types.String `tfsdk:"title"`
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
	
	// Wrap project methods to match ResourceResponse pattern
	r.CreateFunc = func(ctx context.Context, res client.ResourceResponse[client.Project]) (*client.ResourceResponse[client.Project], error) {
		p, err := apiClient.CreateProject(ctx, res.Config)
		if err != nil {
			return nil, err
		}
		return &client.ResourceResponse[client.Project]{
			Name:    p.Name,
			Enabled: &p.Enabled,
			Config:  *p,
		}, nil
	}
	r.GetFunc = func(ctx context.Context, name string) (*client.ResourceResponse[client.Project], error) {
		p, err := apiClient.GetProject(ctx, name)
		if err != nil {
			return nil, err
		}
		return &client.ResourceResponse[client.Project]{
			Name:    p.Name,
			Enabled: &p.Enabled,
			Config:  *p,
		}, nil
	}
	r.UpdateFunc = func(ctx context.Context, res client.ResourceResponse[client.Project]) (*client.ResourceResponse[client.Project], error) {
		p, err := apiClient.UpdateProject(ctx, res.Config)
		if err != nil {
			return nil, err
		}
		return &client.ResourceResponse[client.Project]{
			Name:    p.Name,
			Enabled: &p.Enabled,
			Config:  *p,
		}, nil
	}
	r.DeleteFunc = func(ctx context.Context, name, signature string) error {
		return apiClient.DeleteProject(ctx, name)
	}

	r.PopulateBase = func(m *ProjectResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Enabled = m.Enabled
		b.Description = m.Description
		b.Id = m.Id
		// Signature is left as null/unknown in base
	}
	r.PopulateModel = func(b *BaseResourceModel, m *ProjectResourceModel) {
		m.Name = b.Name
		m.Enabled = b.Enabled
		m.Description = b.Description
		m.Id = b.Id
	}
}

func (r *ProjectResource) MapPlanToClient(ctx context.Context, model *ProjectResourceModel) (client.Project, error) {
	p := client.Project{
		Name: model.Name.ValueString(),
	}
	if !model.Description.IsNull() {
		p.Description = model.Description.ValueString()
	}
	if !model.Title.IsNull() {
		p.Title = model.Title.ValueString()
	}
	if !model.Enabled.IsNull() {
		p.Enabled = model.Enabled.ValueBool()
	}
	if !model.Parent.IsNull() {
		p.Parent = model.Parent.ValueString()
	}
	if !model.Inheritable.IsNull() {
		p.Inheritable = model.Inheritable.ValueBool()
	}
	if !model.DefaultDB.IsNull() {
		p.DefaultDB = model.DefaultDB.ValueString()
	}
	if !model.TagProvider.IsNull() {
		p.TagProvider = model.TagProvider.ValueString()
	}
	if !model.UserSource.IsNull() {
		p.UserSource = model.UserSource.ValueString()
	}
	if !model.IdentityProvider.IsNull() {
		p.IdentityProvider = model.IdentityProvider.ValueString()
	}
	return p, nil
}

func (r *ProjectResource) MapClientToState(ctx context.Context, p *client.Project, model *ProjectResourceModel) error {
	model.Description = stringToNullableString(p.Description)
	model.Title = stringToNullableString(p.Title)
	model.Enabled = types.BoolValue(p.Enabled)
	model.Parent = stringToNullableString(p.Parent)
	model.Inheritable = types.BoolValue(p.Inheritable)
	model.DefaultDB = stringToNullableString(p.DefaultDB)
	model.TagProvider = stringToNullableString(p.TagProvider)
	model.UserSource = stringToNullableString(p.UserSource)
	model.IdentityProvider = stringToNullableString(p.IdentityProvider)
	return nil
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Create(ctx, req, resp, &data, &base)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Read(ctx, req, resp, &data, &base)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Update(ctx, req, resp, &data, &base)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Delete(ctx, req, resp, &data, &base)
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &ProjectResourceModel{
		Name: types.StringValue(req.ID),
	})...)
}
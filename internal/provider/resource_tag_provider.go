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
var _ resource.Resource = &TagProviderResource{}
var _ resource.ResourceWithImportState = &TagProviderResource{}

func NewTagProviderResource() resource.Resource {
	return &TagProviderResource{}
}

// TagProviderResource defines the resource implementation.
type TagProviderResource struct {
	generic GenericIgnitionResource[client.TagProviderConfig, TagProviderResourceModel]
}

// TagProviderResourceModel describes the resource data model.
type TagProviderResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Signature   types.String `tfsdk:"signature"`
}

func (r *TagProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag_provider"
}

func (r *TagProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Tag Provider in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the tag provider.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the tag provider (e.g., standard).",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the tag provider.",
				Optional:    true,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource.",
				Computed:    true,
			},
		},
	}
}

func (r *TagProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(client.IgnitionClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected client.IgnitionClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.generic = GenericIgnitionResource[client.TagProviderConfig, TagProviderResourceModel]{
		Client:       c,
		Handler:      r,
		Module:       "ignition",
		ResourceType: "tag-provider",
		CreateFunc:   c.CreateTagProvider,
		GetFunc:      c.GetTagProvider,
		UpdateFunc:   c.UpdateTagProvider,
		DeleteFunc:   c.DeleteTagProvider,
	}
}

func (r *TagProviderResource) MapPlanToClient(ctx context.Context, model *TagProviderResourceModel) (client.TagProviderConfig, error) {
	return client.TagProviderConfig{
		Type:        model.Type.ValueString(),
		Description: model.Description.ValueString(),
	}, nil
}

func (r *TagProviderResource) MapClientToState(ctx context.Context, name string, config *client.TagProviderConfig, model *TagProviderResourceModel) error {
	model.Name = types.StringValue(name)
	model.Type = types.StringValue(config.Type)
	model.Description = stringToNullableString(config.Description)
	return nil
}

func (r *TagProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TagProviderResourceModel
	var base BaseResourceModel
	r.generic.PopulateBase = func(m *TagProviderResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.generic.PopulateModel = func(b *BaseResourceModel, m *TagProviderResourceModel) {
		m.Name = b.Name
		m.Description = b.Description
		m.Signature = b.Signature
		m.Id = b.Id
	}
	r.generic.Create(ctx, req, resp, &data, &base)
}

func (r *TagProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TagProviderResourceModel
	var base BaseResourceModel
	r.generic.PopulateBase = func(m *TagProviderResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.generic.PopulateModel = func(b *BaseResourceModel, m *TagProviderResourceModel) {
		m.Name = b.Name
		m.Description = b.Description
		m.Signature = b.Signature
		m.Id = b.Id
	}
	r.generic.Read(ctx, req, resp, &data, &base)
}

func (r *TagProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TagProviderResourceModel
	var base BaseResourceModel
	r.generic.PopulateBase = func(m *TagProviderResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.generic.PopulateModel = func(b *BaseResourceModel, m *TagProviderResourceModel) {
		m.Name = b.Name
		m.Description = b.Description
		m.Signature = b.Signature
		m.Id = b.Id
	}
	r.generic.Update(ctx, req, resp, &data, &base)
}

func (r *TagProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TagProviderResourceModel
	var base BaseResourceModel
	r.generic.PopulateBase = func(m *TagProviderResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.generic.Delete(ctx, req, resp, &data, &base)
}

func (r *TagProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &TagProviderResourceModel{
		Id:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	})...)
}
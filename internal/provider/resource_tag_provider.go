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
var _ resource.Resource = &TagProviderResource{}
var _ resource.ResourceWithImportState = &TagProviderResource{}

func NewTagProviderResource() resource.Resource {
	return &TagProviderResource{}
}

// TagProviderResource defines the resource implementation.
type TagProviderResource struct {
	client client.IgnitionClient
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

func (r *TagProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TagProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.TagProviderConfig{
		Type: data.Type.ValueString(),
	}

	if !data.Description.IsNull() {
		config.Description = data.Description.ValueString()
	}

	res := client.ResourceResponse[client.TagProviderConfig]{
		Name:    data.Name.ValueString(),
		Enabled: boolPtr(true),
		Config:  config,
	}

	created, err := r.client.CreateTagProvider(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating tag provider", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	data.Name = types.StringValue(created.Name)
	data.Type = types.StringValue(created.Config.Type)
	data.Description = stringToNullableString(created.Config.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TagProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetTagProvider(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading tag provider", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Name = types.StringValue(res.Name)
	data.Type = types.StringValue(res.Config.Type)
	data.Description = stringToNullableString(res.Config.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TagProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.TagProviderConfig{
		Type: data.Type.ValueString(),
	}

	if !data.Description.IsNull() {
		config.Description = data.Description.ValueString()
	}

	res := client.ResourceResponse[client.TagProviderConfig]{
		Name:      data.Name.ValueString(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	updated, err := r.client.UpdateTagProvider(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tag provider", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	data.Description = stringToNullableString(updated.Config.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TagProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTagProvider(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tag provider", err.Error())
		return
	}
}

func (r *TagProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
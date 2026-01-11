package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseConnectionResource{}
var _ resource.ResourceWithImportState = &DatabaseConnectionResource{}

func NewDatabaseConnectionResource() resource.Resource {
	return &DatabaseConnectionResource{}
}

// DatabaseConnectionResource defines the resource implementation.
type DatabaseConnectionResource struct {
	client client.IgnitionClient
}

// DatabaseConnectionResourceModel describes the resource data model.
type DatabaseConnectionResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Translator types.String `tfsdk:"translator"`
	ConnectURL types.String `tfsdk:"connect_url"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Signature  types.String `tfsdk:"signature"`
}

func (r *DatabaseConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_connection"
}

func (r *DatabaseConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Database Connection in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the database connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the database connection (e.g., MariaDB, PostgreSQL). Maps to 'driver' in Ignition.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"MariaDB",
						"PostgreSQL",
						"MySQL",
						"SQLServer",
						"Oracle",
					),
				},
			},
			"translator": schema.StringAttribute{
				Description: "The SQL translator used to negotiate variances in syntax (e.g., MYSQL, POSTGRESQL).",
				Required:    true,
			},
			"connect_url": schema.StringAttribute{
				Description: "The JDBC connection URL.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username for the database connection.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the database connection.",
				Optional:    true,
				Sensitive:   true,
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *DatabaseConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.DatabaseConfig{
		Driver:     data.Type.ValueString(),
		Translator: data.Translator.ValueString(),
		ConnectURL: data.ConnectURL.ValueString(),
	}

	if !data.Username.IsNull() {
		config.Username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		config.Password = data.Password.ValueString()
	}

	res := client.ResourceResponse[client.DatabaseConfig]{
		Name:    data.Name.ValueString(),
		Enabled: true,
		Config:  config,
	}

	created, err := r.client.CreateDatabaseConnection(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating database connection", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(data.Name.ValueString())
	if created.Config.Driver != "" {
		data.Type = types.StringValue(created.Config.Driver)
	}
	if created.Config.Translator != "" {
		data.Translator = types.StringValue(created.Config.Translator)
	}
	if created.Config.ConnectURL != "" {
		data.ConnectURL = types.StringValue(created.Config.ConnectURL)
	}
	if created.Config.Username != "" {
		data.Username = types.StringValue(created.Config.Username)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetDatabaseConnection(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading database connection", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(data.Name.ValueString())
	if res.Config.Driver != "" {
		data.Type = types.StringValue(res.Config.Driver)
	}
	if res.Config.Translator != "" {
		data.Translator = types.StringValue(res.Config.Translator)
	}
	if res.Config.ConnectURL != "" {
		data.ConnectURL = types.StringValue(res.Config.ConnectURL)
	}
	if res.Config.Username != "" {
		data.Username = types.StringValue(res.Config.Username)
	} else {
		data.Username = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseConnectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.DatabaseConfig{
		Driver:     data.Type.ValueString(),
		Translator: data.Translator.ValueString(),
		ConnectURL: data.ConnectURL.ValueString(),
	}

	if !data.Username.IsNull() {
		config.Username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		config.Password = data.Password.ValueString()
	}

	res := client.ResourceResponse[client.DatabaseConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   true,
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	updated, err := r.client.UpdateDatabaseConnection(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating database connection", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	if updated.Config.Driver != "" {
		data.Type = types.StringValue(updated.Config.Driver)
	}
	if updated.Config.Translator != "" {
		data.Translator = types.StringValue(updated.Config.Translator)
	}
	if updated.Config.ConnectURL != "" {
		data.ConnectURL = types.StringValue(updated.Config.ConnectURL)
	}
	if updated.Config.Username != "" {
		data.Username = types.StringValue(updated.Config.Username)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDatabaseConnection(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting database connection", err.Error())
		return
	}
}

func (r *DatabaseConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
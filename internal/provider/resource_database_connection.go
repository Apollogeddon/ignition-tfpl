package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	GenericIgnitionResource[client.DatabaseConfig, DatabaseConnectionResourceModel]
}

// DatabaseConnectionResourceModel describes the resource data model.
type DatabaseConnectionResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Type        types.String `tfsdk:"type"`
	Translator  types.String `tfsdk:"translator"`
	ConnectURL  types.String `tfsdk:"connect_url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Signature   types.String `tfsdk:"signature"`
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
			"description": schema.StringAttribute{
				Description: "The description of the database connection.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the database connection is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
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

	r.Client = client
	r.Handler = r
	r.Module = "ignition"
	r.ResourceType = "database-connection"
	r.CreateFunc = client.CreateDatabaseConnection
	r.GetFunc = client.GetDatabaseConnection
	r.UpdateFunc = client.UpdateDatabaseConnection
	r.DeleteFunc = client.DeleteDatabaseConnection
	r.PopulateBase = func(m *DatabaseConnectionResourceModel, b *BaseResourceModel) {
		b.Name = m.Name
		b.Enabled = m.Enabled
		b.Description = m.Description
		b.Signature = m.Signature
		b.Id = m.Id
	}
	r.PopulateModel = func(b *BaseResourceModel, m *DatabaseConnectionResourceModel) {
		m.Name = b.Name
		m.Enabled = b.Enabled
		m.Description = b.Description
		m.Signature = b.Signature
		m.Id = b.Id
	}
}

func (r *DatabaseConnectionResource) MapPlanToClient(ctx context.Context, model *DatabaseConnectionResourceModel) (client.DatabaseConfig, error) {
	config := client.DatabaseConfig{
		Driver:     model.Type.ValueString(),
		Translator: model.Translator.ValueString(),
		ConnectURL: model.ConnectURL.ValueString(),
	}

	if !model.Username.IsNull() {
		config.Username = model.Username.ValueString()
	}
	if !model.Password.IsNull() {
		encrypted, err := r.Client.EncryptSecret(ctx, model.Password.ValueString())
		if err != nil {
			return client.DatabaseConfig{}, err
		}
		config.Password = encrypted
	}

	return config, nil
}

func (r *DatabaseConnectionResource) MapClientToState(ctx context.Context, name string, config *client.DatabaseConfig, model *DatabaseConnectionResourceModel) error {
	model.Name = types.StringValue(name)

	if config.Driver != "" {
		model.Type = types.StringValue(config.Driver)
	}
	if config.Translator != "" {
		model.Translator = types.StringValue(config.Translator)
	}
	if config.ConnectURL != "" {
		model.ConnectURL = types.StringValue(config.ConnectURL)
	}
	if config.Username != "" {
		model.Username = types.StringValue(config.Username)
	} else {
		model.Username = types.StringNull()
	}
	// The Name is usually the ID, which is returned in the wrap but let's assume 
	// it's consistent. However, the generic base handles Name.
	// If the name is in the state, we keep it.
	return nil
}

func (r *DatabaseConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseConnectionResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Create(ctx, req, resp, &data, &base)
}

func (r *DatabaseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseConnectionResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Read(ctx, req, resp, &data, &base)
}

func (r *DatabaseConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseConnectionResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Update(ctx, req, resp, &data, &base)
}

func (r *DatabaseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseConnectionResourceModel
	var base BaseResourceModel
	r.GenericIgnitionResource.Delete(ctx, req, resp, &data, &base)
}

func (r *DatabaseConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &DatabaseConnectionResourceModel{
		Id:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	})...)
}
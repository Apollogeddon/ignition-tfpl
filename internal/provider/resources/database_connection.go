package resources

import (
	"context"
	"fmt"

	"github.com/apollogeddon/ignition-tfpl/internal/client"
	"github.com/apollogeddon/ignition-tfpl/internal/provider/base"
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
	base.GenericIgnitionResource[client.DatabaseConfig, DatabaseConnectionResourceModel]
}

// DatabaseConnectionResourceModel describes the resource data model.
type DatabaseConnectionResourceModel struct {
	base.BaseResourceModel
	Type        types.String `tfsdk:"type"`
	Translator  types.String `tfsdk:"translator"`
	ConnectURL  types.String `tfsdk:"connect_url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
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
	r.Module = "ignition"
	r.ResourceType = "database-connection"
	r.CreateFunc = apiClient.CreateDatabaseConnection
	r.GetFunc = apiClient.GetDatabaseConnection
	r.UpdateFunc = apiClient.UpdateDatabaseConnection
	r.DeleteFunc = apiClient.DeleteDatabaseConnection
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
	
	// Ensure signature is preserved
	// The signature is handled by the generic base if provided in the response
	return nil
}

func (r *DatabaseConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseConnectionResourceModel
	r.GenericIgnitionResource.Create(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *DatabaseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseConnectionResourceModel
	r.GenericIgnitionResource.Read(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *DatabaseConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseConnectionResourceModel
	r.GenericIgnitionResource.Update(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *DatabaseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseConnectionResourceModel
	r.GenericIgnitionResource.Delete(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *DatabaseConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &DatabaseConnectionResourceModel{
		BaseResourceModel: base.BaseResourceModel{
			Id:   types.StringValue(req.ID),
			Name: types.StringValue(req.ID),
		},
	})...)
}
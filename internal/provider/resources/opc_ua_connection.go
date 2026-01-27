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
var _ resource.Resource = &OpcUaConnectionResource{}
var _ resource.ResourceWithImportState = &OpcUaConnectionResource{}

func NewOpcUaConnectionResource() resource.Resource {
	return &OpcUaConnectionResource{}
}

// OpcUaConnectionResource defines the resource implementation.
type OpcUaConnectionResource struct {
	base.GenericIgnitionResource[client.OpcUaConnectionConfig, OpcUaConnectionResourceModel]
}

// OpcUaConnectionResourceModel describes the resource data model.
type OpcUaConnectionResourceModel struct {
	base.BaseResourceModel
	Type           types.String `tfsdk:"type"`
	DiscoveryURL   types.String `tfsdk:"discovery_url"`
	EndpointURL    types.String `tfsdk:"endpoint_url"`
	SecurityPolicy types.String `tfsdk:"security_policy"`
	SecurityMode   types.String `tfsdk:"security_mode"`
}

func (r *OpcUaConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_opc_ua_connection"
}

func (r *OpcUaConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OPC UA Connection in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the OPC UA connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the OPC UA connection.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the OPC UA connection is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"type": schema.StringAttribute{
				Description: "The type of the OPC UA connection (e.g., com.inductiveautomation.OpcUaServerType).",
				Optional:    true,
				Computed:    true,
				// Default to standard OPC UA Server type if not provided
			},
			"discovery_url": schema.StringAttribute{
				Description: "The discovery URL of the OPC UA server.",
				Required:    true,
			},
			"endpoint_url": schema.StringAttribute{
				Description: "The endpoint URL of the OPC UA server.",
				Required:    true,
			},
			"security_policy": schema.StringAttribute{
				Description: "The security policy to use (e.g., None, Basic256Sha256).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"None",
						"Basic128Rsa15",
						"Basic256",
						"Basic256Sha256",
						"Aes128_Sha256_RsaOaep",
						"Aes256_Sha256_RsaPss",
					),
				},
			},
			"security_mode": schema.StringAttribute{
				Description: "The security mode to use (e.g., None, SignAndEncrypt).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"None",
						"Sign",
						"SignAndEncrypt",
					),
				},
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource, used for updates and deletes.",
				Computed:    true,
			},
		},
	}
}

func (r *OpcUaConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.GenericIgnitionResource = base.GenericIgnitionResource[client.OpcUaConnectionConfig, OpcUaConnectionResourceModel]{
		Client:       c,
		Handler:      r,
		Module:       "ignition",
		ResourceType: "opc-connection",
		CreateFunc:   c.CreateOpcUaConnection,
		GetFunc:      c.GetOpcUaConnection,
		UpdateFunc:   c.UpdateOpcUaConnection,
		DeleteFunc:   c.DeleteOpcUaConnection,
	}
}

func (r *OpcUaConnectionResource) MapPlanToClient(ctx context.Context, model *OpcUaConnectionResourceModel) (client.OpcUaConnectionConfig, error) {
	profileType := "com.inductiveautomation.OpcUaServerType"
	if !model.Type.IsNull() {
		profileType = model.Type.ValueString()
	}

	return client.OpcUaConnectionConfig{
		Profile: client.OpcUaConnectionProfile{
			Type: profileType,
		},
		Settings: client.OpcUaConnectionSettings{
			Endpoint: client.OpcUaConnectionEndpoint{
				DiscoveryURL:   model.DiscoveryURL.ValueString(),
				EndpointURL:    model.EndpointURL.ValueString(),
				SecurityPolicy: model.SecurityPolicy.ValueString(),
				SecurityMode:   model.SecurityMode.ValueString(),
			},
		},
	}, nil
}

func (r *OpcUaConnectionResource) MapClientToState(ctx context.Context, name string, config *client.OpcUaConnectionConfig, model *OpcUaConnectionResourceModel) error {
	model.Name = types.StringValue(name)
	model.Type = types.StringValue(config.Profile.Type)
	model.DiscoveryURL = types.StringValue(config.Settings.Endpoint.DiscoveryURL)
	model.EndpointURL = types.StringValue(config.Settings.Endpoint.EndpointURL)
	model.SecurityPolicy = types.StringValue(config.Settings.Endpoint.SecurityPolicy)
	model.SecurityMode = types.StringValue(config.Settings.Endpoint.SecurityMode)
	return nil
}

func (r *OpcUaConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpcUaConnectionResourceModel
	r.GenericIgnitionResource.Create(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *OpcUaConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpcUaConnectionResourceModel
	r.GenericIgnitionResource.Read(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *OpcUaConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OpcUaConnectionResourceModel
	r.GenericIgnitionResource.Update(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *OpcUaConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpcUaConnectionResourceModel
	r.GenericIgnitionResource.Delete(ctx, req, resp, &data, &data.BaseResourceModel)
}

func (r *OpcUaConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, &OpcUaConnectionResourceModel{
		BaseResourceModel: base.BaseResourceModel{
			Id:   types.StringValue(req.ID),
			Name: types.StringValue(req.ID),
		},
	})...)
}

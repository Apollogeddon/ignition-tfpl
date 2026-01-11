package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GanOutgoingResource{}
var _ resource.ResourceWithImportState = &GanOutgoingResource{}

func NewGanOutgoingResource() resource.Resource {
	return &GanOutgoingResource{}
}

// GanOutgoingResource defines the resource implementation.
type GanOutgoingResource struct {
	client client.IgnitionClient
}

// GanOutgoingResourceModel describes the resource data model.
type GanOutgoingResourceModel struct {
	Id                       types.String  `tfsdk:"id"`
	Name                     types.String  `tfsdk:"name"`
	Description              types.String  `tfsdk:"description"`
	Enabled                  types.Bool    `tfsdk:"enabled"`
	Host                     types.String  `tfsdk:"host"`
	Port                     types.Int64   `tfsdk:"port"`
	UseSSL                   types.Bool    `tfsdk:"use_ssl"`
	PingRateMillis           types.Float64 `tfsdk:"ping_rate_millis"`
	PingTimeoutMillis        types.Float64 `tfsdk:"ping_timeout_millis"`
	PingMaxMissed            types.Float64 `tfsdk:"ping_max_missed"`
	WsTimeoutMillis          types.Float64 `tfsdk:"ws_timeout_millis"`
	HttpConnectTimeoutMillis types.Float64 `tfsdk:"http_connect_timeout_millis"`
	HttpReadTimeoutMillis    types.Float64 `tfsdk:"http_read_timeout_millis"`
	SendThreads              types.Float64 `tfsdk:"send_threads"`
	ReceiveThreads           types.Float64 `tfsdk:"receive_threads"`
	Signature                types.String  `tfsdk:"signature"`
}

func (r *GanOutgoingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gan_outgoing"
}

func (r *GanOutgoingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Outgoing Gateway Network Connection in Ignition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the connection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the connection.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the connection is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"host": schema.StringAttribute{
				Description: "The address of the remote server.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The port of the remote server.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(8060),
			},
			"use_ssl": schema.BoolAttribute{
				Description: "Whether to use SSL for the connection.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ping_rate_millis": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(2000),
			},
			"ping_timeout_millis": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(60000),
			},
			"ping_max_missed": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(3),
			},
			"ws_timeout_millis": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(10000),
			},
			"http_connect_timeout_millis": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(10000),
			},
			"http_read_timeout_millis": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(30000),
			},
			"send_threads": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(1),
			},
			"receive_threads": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Default:  float64default.StaticFloat64(1),
			},
			"signature": schema.StringAttribute{
				Description: "The signature of the resource.",
				Computed:    true,
			},
		},
	}
}

func (r *GanOutgoingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GanOutgoingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GanOutgoingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.GanOutgoingConfig{
		Host:                     data.Host.ValueString(),
		Port:                     int(data.Port.ValueInt64()),
		UseSSL:                   data.UseSSL.ValueBool(),
		PingRateMillis:           data.PingRateMillis.ValueFloat64(),
		PingTimeoutMillis:        data.PingTimeoutMillis.ValueFloat64(),
		PingMaxMissed:            data.PingMaxMissed.ValueFloat64(),
		WsTimeoutMillis:          data.WsTimeoutMillis.ValueFloat64(),
		HttpConnectTimeoutMillis: data.HttpConnectTimeoutMillis.ValueFloat64(),
		HttpReadTimeoutMillis:    data.HttpReadTimeoutMillis.ValueFloat64(),
		SendThreads:              data.SendThreads.ValueFloat64(),
		ReceiveThreads:           data.ReceiveThreads.ValueFloat64(),
	}

	res := client.ResourceResponse[client.GanOutgoingConfig]{
		Name:    data.Name.ValueString(),
		Enabled: data.Enabled.ValueBool(),
		Config:  config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	created, err := r.client.CreateGanOutgoing(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error creating GAN connection", err.Error())
		return
	}

	data.Signature = types.StringValue(created.Signature)
	data.Id = types.StringValue(created.Name)
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GanOutgoingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GanOutgoingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetGanOutgoing(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading GAN connection", err.Error())
		return
	}

	data.Signature = types.StringValue(res.Signature)
	data.Id = types.StringValue(res.Name)
	data.Enabled = types.BoolValue(res.Enabled)
	data.Description = types.StringValue(res.Description)
	
	data.Host = types.StringValue(res.Config.Host)
	data.Port = types.Int64Value(int64(res.Config.Port))
	data.UseSSL = types.BoolValue(res.Config.UseSSL)
	data.PingRateMillis = types.Float64Value(res.Config.PingRateMillis)
	data.PingTimeoutMillis = types.Float64Value(res.Config.PingTimeoutMillis)
	data.PingMaxMissed = types.Float64Value(res.Config.PingMaxMissed)
	data.WsTimeoutMillis = types.Float64Value(res.Config.WsTimeoutMillis)
	data.HttpConnectTimeoutMillis = types.Float64Value(res.Config.HttpConnectTimeoutMillis)
	data.HttpReadTimeoutMillis = types.Float64Value(res.Config.HttpReadTimeoutMillis)
	data.SendThreads = types.Float64Value(res.Config.SendThreads)
	data.ReceiveThreads = types.Float64Value(res.Config.ReceiveThreads)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GanOutgoingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GanOutgoingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.GanOutgoingConfig{
		Host:                     data.Host.ValueString(),
		Port:                     int(data.Port.ValueInt64()),
		UseSSL:                   data.UseSSL.ValueBool(),
		PingRateMillis:           data.PingRateMillis.ValueFloat64(),
		PingTimeoutMillis:        data.PingTimeoutMillis.ValueFloat64(),
		PingMaxMissed:            data.PingMaxMissed.ValueFloat64(),
		WsTimeoutMillis:          data.WsTimeoutMillis.ValueFloat64(),
		HttpConnectTimeoutMillis: data.HttpConnectTimeoutMillis.ValueFloat64(),
		HttpReadTimeoutMillis:    data.HttpReadTimeoutMillis.ValueFloat64(),
		SendThreads:              data.SendThreads.ValueFloat64(),
		ReceiveThreads:           data.ReceiveThreads.ValueFloat64(),
	}

	res := client.ResourceResponse[client.GanOutgoingConfig]{
		Name:      data.Name.ValueString(),
		Enabled:   data.Enabled.ValueBool(),
		Signature: data.Signature.ValueString(),
		Config:    config,
	}

	if !data.Description.IsNull() {
		res.Description = data.Description.ValueString()
	}

	updated, err := r.client.UpdateGanOutgoing(ctx, res)
	if err != nil {
		resp.Diagnostics.AddError("Error updating GAN connection", err.Error())
		return
	}

	data.Signature = types.StringValue(updated.Signature)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GanOutgoingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GanOutgoingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGanOutgoing(ctx, data.Name.ValueString(), data.Signature.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting GAN connection", err.Error())
		return
	}
}

func (r *GanOutgoingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

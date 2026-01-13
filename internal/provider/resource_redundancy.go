package provider

import (
	"context"
	"fmt"

	"github.com/apollogeddon/terraform-provider-ignition/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RedundancyResource{}

func NewRedundancyResource() resource.Resource {
	return &RedundancyResource{}
}

// RedundancyResource defines the resource implementation.
type RedundancyResource struct {
	client client.IgnitionClient
}

// RedundancyResourceModel describes the resource data model.
type RedundancyResourceModel struct {
	Id                  types.String         `tfsdk:"id"`
	Role                types.String         `tfsdk:"role"`
	ActiveHistoryLevel  types.String         `tfsdk:"active_history_level"`
	JoinWaitTime        types.Int64          `tfsdk:"join_wait_time"`
	RecoveryMode        types.String         `tfsdk:"recovery_mode"`
	AllowHistoryCleanup types.Bool           `tfsdk:"allow_history_cleanup"`
	GatewayNetworkSetup *GatewayNetworkSetup `tfsdk:"gateway_network_setup"`
}

type GatewayNetworkSetup struct {
	Host               types.String  `tfsdk:"host"`
	Port               types.Int64   `tfsdk:"port"`
	EnableSsl          types.Bool    `tfsdk:"enable_ssl"`
	PingRate           types.Float64 `tfsdk:"ping_rate"`
	PingTimeout        types.Float64 `tfsdk:"ping_timeout"`
	PingMaxMissed      types.Float64 `tfsdk:"ping_max_missed"`
	WebsocketTimeout   types.Float64 `tfsdk:"websocket_timeout"`
	HttpConnectTimeout types.Float64 `tfsdk:"http_connect_timeout"`
	HttpReadTimeout    types.Float64 `tfsdk:"http_read_timeout"`
	SendThreads        types.Float64 `tfsdk:"send_threads"`
	ReceiveThreads     types.Float64 `tfsdk:"receive_threads"`
}

func (r *RedundancyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redundancy"
}

func (r *RedundancyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Gateway Redundancy Settings. This is a singleton resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"role": schema.StringAttribute{
				Description: "The node's redundancy role (Independent, Backup, Master).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Independent", "Backup", "Master"),
				},
			},
			"active_history_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Full"),
				Validators: []validator.String{
					stringvalidator.OneOf("Partial", "Full"),
				},
			},
			"join_wait_time": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(10000),
			},
			"recovery_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Automatic"),
			},
			"allow_history_cleanup": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"gateway_network_setup": schema.SingleNestedAttribute{
				Description: "Gateway network settings to establish a connection to the redundant master. (Only applies to Backup)",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Optional: true,
					},
					"port": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(8060),
					},
					"enable_ssl": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"ping_rate":            schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(2000)},
					"ping_timeout":         schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(60000)},
					"ping_max_missed":      schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(3)},
					"websocket_timeout":    schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(10000)},
					"http_connect_timeout": schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(10000)},
					"http_read_timeout":    schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(30000)},
					"send_threads":         schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(1)},
					"receive_threads":      schema.Float64Attribute{Optional: true, Computed: true, Default: float64default.StaticFloat64(1)},
				},
			},
		},
	}
}

func (r *RedundancyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RedundancyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RedundancyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.RedundancyConfig{
		Role:                data.Role.ValueString(),
		ActiveHistoryLevel:  data.ActiveHistoryLevel.ValueString(),
		JoinWaitTime:        int(data.JoinWaitTime.ValueInt64()),
		RecoveryMode:        data.RecoveryMode.ValueString(),
		AllowHistoryCleanup: data.AllowHistoryCleanup.ValueBool(),
	}

	if data.GatewayNetworkSetup != nil {
		config.GatewayNetworkSetup = &struct {
			Host               string  `json:"host,omitempty"`
			Port               int     `json:"port,omitempty"`
			EnableSsl          bool    `json:"enableSsl,omitempty"`
			PingRate           float64 `json:"pingRate,omitempty"`
			PingTimeout        float64 `json:"pingTimeout,omitempty"`
			PingMaxMissed      float64 `json:"pingMaxMissed,omitempty"`
			WebsocketTimeout   float64 `json:"websocketTimeout,omitempty"`
			HttpConnectTimeout float64 `json:"httpConnectTimeout,omitempty"`
			HttpReadTimeout    float64 `json:"httpReadTimeout,omitempty"`
			SendThreads        float64 `json:"sendThreads,omitempty"`
			ReceiveThreads     float64 `json:"receiveThreads,omitempty"`
		}{
			Host:               data.GatewayNetworkSetup.Host.ValueString(),
			Port:               int(data.GatewayNetworkSetup.Port.ValueInt64()),
			EnableSsl:          data.GatewayNetworkSetup.EnableSsl.ValueBool(),
			PingRate:           data.GatewayNetworkSetup.PingRate.ValueFloat64(),
			PingTimeout:        data.GatewayNetworkSetup.PingTimeout.ValueFloat64(),
			PingMaxMissed:      data.GatewayNetworkSetup.PingMaxMissed.ValueFloat64(),
			WebsocketTimeout:   data.GatewayNetworkSetup.WebsocketTimeout.ValueFloat64(),
			HttpConnectTimeout: data.GatewayNetworkSetup.HttpConnectTimeout.ValueFloat64(),
			HttpReadTimeout:    data.GatewayNetworkSetup.HttpReadTimeout.ValueFloat64(),
			SendThreads:        data.GatewayNetworkSetup.SendThreads.ValueFloat64(),
			ReceiveThreads:     data.GatewayNetworkSetup.ReceiveThreads.ValueFloat64(),
		}
	}

	err := r.client.UpdateRedundancyConfig(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError("Error updating redundancy config", err.Error())
		return
	}

	data.Id = types.StringValue("gateway_redundancy")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RedundancyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RedundancyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetRedundancyConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading redundancy config", err.Error())
		return
	}

	data.Id = types.StringValue("gateway_redundancy")
	data.Role = types.StringValue(res.Role)
	data.ActiveHistoryLevel = types.StringValue(res.ActiveHistoryLevel)
	data.JoinWaitTime = types.Int64Value(int64(res.JoinWaitTime))
	data.RecoveryMode = types.StringValue(res.RecoveryMode)
	data.AllowHistoryCleanup = types.BoolValue(res.AllowHistoryCleanup)

	if res.GatewayNetworkSetup != nil {
		data.GatewayNetworkSetup = &GatewayNetworkSetup{
			Host:               types.StringValue(res.GatewayNetworkSetup.Host),
			Port:               types.Int64Value(int64(res.GatewayNetworkSetup.Port)),
			EnableSsl:          types.BoolValue(res.GatewayNetworkSetup.EnableSsl),
			PingRate:           types.Float64Value(res.GatewayNetworkSetup.PingRate),
			PingTimeout:        types.Float64Value(res.GatewayNetworkSetup.PingTimeout),
			PingMaxMissed:      types.Float64Value(res.GatewayNetworkSetup.PingMaxMissed),
			WebsocketTimeout:   types.Float64Value(res.GatewayNetworkSetup.WebsocketTimeout),
			HttpConnectTimeout: types.Float64Value(res.GatewayNetworkSetup.HttpConnectTimeout),
			HttpReadTimeout:    types.Float64Value(res.GatewayNetworkSetup.HttpReadTimeout),
			SendThreads:        types.Float64Value(res.GatewayNetworkSetup.SendThreads),
			ReceiveThreads:     types.Float64Value(res.GatewayNetworkSetup.ReceiveThreads),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RedundancyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RedundancyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := client.RedundancyConfig{
		Role:                data.Role.ValueString(),
		ActiveHistoryLevel:  data.ActiveHistoryLevel.ValueString(),
		JoinWaitTime:        int(data.JoinWaitTime.ValueInt64()),
		RecoveryMode:        data.RecoveryMode.ValueString(),
		AllowHistoryCleanup: data.AllowHistoryCleanup.ValueBool(),
	}

	if data.GatewayNetworkSetup != nil {
		config.GatewayNetworkSetup = &struct {
			Host               string  `json:"host,omitempty"`
			Port               int     `json:"port,omitempty"`
			EnableSsl          bool    `json:"enableSsl,omitempty"`
			PingRate           float64 `json:"pingRate,omitempty"`
			PingTimeout        float64 `json:"pingTimeout,omitempty"`
			PingMaxMissed      float64 `json:"pingMaxMissed,omitempty"`
			WebsocketTimeout   float64 `json:"websocketTimeout,omitempty"`
			HttpConnectTimeout float64 `json:"httpConnectTimeout,omitempty"`
			HttpReadTimeout    float64 `json:"httpReadTimeout,omitempty"`
			SendThreads        float64 `json:"sendThreads,omitempty"`
			ReceiveThreads     float64 `json:"receiveThreads,omitempty"`
		}{
			Host:               data.GatewayNetworkSetup.Host.ValueString(),
			Port:               int(data.GatewayNetworkSetup.Port.ValueInt64()),
			EnableSsl:          data.GatewayNetworkSetup.EnableSsl.ValueBool(),
			PingRate:           data.GatewayNetworkSetup.PingRate.ValueFloat64(),
			PingTimeout:        data.GatewayNetworkSetup.PingTimeout.ValueFloat64(),
			PingMaxMissed:      data.GatewayNetworkSetup.PingMaxMissed.ValueFloat64(),
			WebsocketTimeout:   data.GatewayNetworkSetup.WebsocketTimeout.ValueFloat64(),
			HttpConnectTimeout: data.GatewayNetworkSetup.HttpConnectTimeout.ValueFloat64(),
			HttpReadTimeout:    data.GatewayNetworkSetup.HttpReadTimeout.ValueFloat64(),
			SendThreads:        data.GatewayNetworkSetup.SendThreads.ValueFloat64(),
			ReceiveThreads:     data.GatewayNetworkSetup.ReceiveThreads.ValueFloat64(),
		}
	}

	err := r.client.UpdateRedundancyConfig(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError("Error updating redundancy config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RedundancyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// For singletons, we usually don't "delete" the config, just disable it or leave it as is.
	// We'll set the role to Independent to effectively "disable" redundancy.
	config := client.RedundancyConfig{
		Role: "Independent",
	}
	_ = r.client.UpdateRedundancyConfig(ctx, config)
}

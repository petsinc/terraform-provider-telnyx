package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

var (
	_ resource.Resource              = &CredentialConnectionResource{}
	_ resource.ResourceWithConfigure = &CredentialConnectionResource{}
)

func NewCredentialConnectionResource() resource.Resource {
	return &CredentialConnectionResource{}
}

type CredentialConnectionResource struct {
	client *telnyx.TelnyxClient
}

type CredentialConnectionResourceModel struct {
	ID                               types.String `tfsdk:"id"`
	ConnectionName                   types.String `tfsdk:"connection_name"`
	Username                         types.String `tfsdk:"username"`
	Password                         types.String `tfsdk:"password"`
	Active                           types.Bool   `tfsdk:"active"`
	AnchorsiteOverride               types.String `tfsdk:"anchorsite_override"`
	DefaultOnHoldComfortNoiseEnabled types.Bool   `tfsdk:"default_on_hold_comfort_noise_enabled"`
	DTMFType                         types.String `tfsdk:"dtmf_type"`
	EncodeContactHeaderEnabled       types.Bool   `tfsdk:"encode_contact_header_enabled"`
	OnnetT38PassthroughEnabled       types.Bool   `tfsdk:"onnet_t38_passthrough_enabled"`
	MicrosoftTeamsSBC                types.Bool   `tfsdk:"microsoft_teams_sbc"`
	WebhookEventURL                  types.String `tfsdk:"webhook_event_url"`
	WebhookEventFailoverURL          types.String `tfsdk:"webhook_event_failover_url"`
	WebhookAPIVersion                types.String `tfsdk:"webhook_api_version"`
	WebhookTimeoutSecs               types.Int64  `tfsdk:"webhook_timeout_secs"`
	RTCPSettings                     types.Object `tfsdk:"rtcp_settings"`
	Inbound                          types.Object `tfsdk:"inbound"`
	Outbound                         types.Object `tfsdk:"outbound"`
}

func (r *CredentialConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_connection"
}

func (r *CredentialConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx Credential Connections",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the credential connection",
				Computed:    true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"connection_name": schema.StringAttribute{
				Description: "Name of the credential connection",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for the credential connection",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the credential connection",
				Required:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Specifies whether the credential connection is active or not",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"anchorsite_override": schema.StringAttribute{
				Description: "Anchorsite override setting",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Latency"),
			},
			"default_on_hold_comfort_noise_enabled": schema.BoolAttribute{
				Description: "Default on-hold comfort noise enabled setting",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"dtmf_type": schema.StringAttribute{
				Description: "DTMF type",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("RFC 2833"),
			},
			"encode_contact_header_enabled": schema.BoolAttribute{
				Description: "Encode contact header enabled setting",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"onnet_t38_passthrough_enabled": schema.BoolAttribute{
				Description: "On-net T38 passthrough enabled setting",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"microsoft_teams_sbc": schema.BoolAttribute{
				Description: "Microsoft Teams SBC setting",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"webhook_event_url": schema.StringAttribute{
				Description: "Webhook event URL",
				Optional:    true,
				Computed:    true,
			},
			"webhook_event_failover_url": schema.StringAttribute{
				Description: "Webhook event failover URL",
				Optional:    true,
				Computed:    true,
			},
			"webhook_api_version": schema.StringAttribute{
				Description: "Webhook API version",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1"),
			},
			"webhook_timeout_secs": schema.Int64Attribute{
				Description: "Webhook timeout in seconds",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(25),
			},
			"rtcp_settings": schema.SingleNestedAttribute{
				Description: "RTCP settings",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"port":                  types.StringType,
						"capture_enabled":       types.BoolType,
						"report_frequency_secs": types.Int64Type,
					},
					map[string]attr.Value{
						"port":                  types.StringValue("rtp+1"),
						"capture_enabled":       types.BoolValue(false),
						"report_frequency_secs": types.Int64Value(5),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"port": schema.StringAttribute{
						Description: "Port for RTCP",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("rtp+1"),
					},
					"capture_enabled": schema.BoolAttribute{
						Description: "Capture enabled for RTCP",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"report_frequency_secs": schema.Int64Attribute{
						Description: "Report frequency for RTCP in seconds",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
					},
				},
			},
			"inbound": schema.SingleNestedAttribute{
				Description: "Inbound settings",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"ani_number_format":           types.StringType,
						"dnis_number_format":          types.StringType,
						"codecs":                      types.ListType{ElemType: types.StringType},
						"default_routing_method":      types.StringType,
						"channel_limit":               types.Int64Type,
						"generate_ringback_tone":      types.BoolType,
						"isup_headers_enabled":        types.BoolType,
						"prack_enabled":               types.BoolType,
						"privacy_zone_enabled":        types.BoolType,
						"sip_compact_headers_enabled": types.BoolType,
						"timeout_1xx_secs":            types.Int64Type,
						"timeout_2xx_secs":            types.Int64Type,
						"shaken_stir_enabled":         types.BoolType,
					},
					map[string]attr.Value{
						"ani_number_format":           types.StringValue("E.164-national"),
						"dnis_number_format":          types.StringValue("e164"),
						"codecs":                      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("G722"), types.StringValue("G711U"), types.StringValue("G711A"), types.StringValue("G729"), types.StringValue("OPUS"), types.StringValue("H.264")}),
						"default_routing_method":      types.StringValue("sequential"),
						"channel_limit":               types.Int64Value(10),
						"generate_ringback_tone":      types.BoolValue(true),
						"isup_headers_enabled":        types.BoolValue(true),
						"prack_enabled":               types.BoolValue(true),
						"privacy_zone_enabled":        types.BoolValue(true),
						"sip_compact_headers_enabled": types.BoolValue(true),
						"timeout_1xx_secs":            types.Int64Value(3),
						"timeout_2xx_secs":            types.Int64Value(90),
						"shaken_stir_enabled":         types.BoolValue(true),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"ani_number_format": schema.StringAttribute{
						Description: "ANI number format",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("E.164-national"),
					},
					"dnis_number_format": schema.StringAttribute{
						Description: "DNIS number format",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("e164"),
					},
					"codecs": schema.ListAttribute{
						Description: "List of codecs",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("G722"), types.StringValue("G711U"), types.StringValue("G711A"), types.StringValue("G729"), types.StringValue("OPUS"), types.StringValue("H.264")})),
					},
					"default_routing_method": schema.StringAttribute{
						Description: "Default routing method",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("sequential"),
					},
					"channel_limit": schema.Int64Attribute{
						Description: "Channel limit",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(10),
					},
					"generate_ringback_tone": schema.BoolAttribute{
						Description: "Generate ringback tone",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"isup_headers_enabled": schema.BoolAttribute{
						Description: "ISUP headers enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"prack_enabled": schema.BoolAttribute{
						Description: "PRACK enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"privacy_zone_enabled": schema.BoolAttribute{
						Description: "Privacy zone enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"sip_compact_headers_enabled": schema.BoolAttribute{
						Description: "SIP compact headers enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"timeout_1xx_secs": schema.Int64Attribute{
						Description: "Timeout for 1xx responses in seconds",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3),
					},
					"timeout_2xx_secs": schema.Int64Attribute{
						Description: "Timeout for 2xx responses in seconds",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(90),
					},
					"shaken_stir_enabled": schema.BoolAttribute{
						Description: "SHAKEN/STIR enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
			},
			"outbound": schema.SingleNestedAttribute{
				Description: "Outbound settings",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"ani_override":              types.StringType,
						"ani_override_type":         types.StringType,
						"call_parking_enabled":      types.BoolType,
						"channel_limit":             types.Int64Type,
						"generate_ringback_tone":    types.BoolType,
						"instant_ringback_enabled":  types.BoolType,
						"localization":              types.StringType,
						"outbound_voice_profile_id": types.StringType,
						"t38_reinvite_source":       types.StringType,
					},
					map[string]attr.Value{
						"ani_override":              types.StringValue("+12345678901"),
						"ani_override_type":         types.StringValue("always"),
						"call_parking_enabled":      types.BoolValue(true),
						"channel_limit":             types.Int64Value(10),
						"generate_ringback_tone":    types.BoolValue(true),
						"instant_ringback_enabled":  types.BoolValue(false),
						"localization":              types.StringValue("US"),
						"outbound_voice_profile_id": types.StringValue(""),
						"t38_reinvite_source":       types.StringValue("customer"),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"ani_override": schema.StringAttribute{
						Description: "ANI override",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("+12345678901"),
					},
					"ani_override_type": schema.StringAttribute{
						Description: "ANI override type",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("always"),
					},
					"call_parking_enabled": schema.BoolAttribute{
						Description: "Call parking enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"channel_limit": schema.Int64Attribute{
						Description: "Channel limit",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(10),
					},
					"generate_ringback_tone": schema.BoolAttribute{
						Description: "Generate ringback tone",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"instant_ringback_enabled": schema.BoolAttribute{
						Description: "Instant ringback enabled",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"localization": schema.StringAttribute{
						Description: "Localization",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("US"),
					},
					"outbound_voice_profile_id": schema.StringAttribute{
						Description: "Outbound voice profile ID",
						Optional:    true,
					},
					"t38_reinvite_source": schema.StringAttribute{
						Description: "T38 reinvite source",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("customer"),
					},
				},
			},
		},
	}
}

func (r *CredentialConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		client, ok := req.ProviderData.(*telnyx.TelnyxClient)
		if !ok {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected *telnyx.TelnyxClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
			return
		}
		r.client = client
		tflog.Info(ctx, "Configured Telnyx client for CredentialConnectionResource")
	}
}

func (r *CredentialConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CredentialConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Credential Connection", map[string]interface{}{
		"connection_name": plan.ConnectionName.ValueString(),
	})

	rtcpSettingsAttributes := plan.RTCPSettings.Attributes()
	inboundAttributes := plan.Inbound.Attributes()
	outboundAttributes := plan.Outbound.Attributes()

	codecs, diagCodecs := convertListToStrings(ctx, inboundAttributes["codecs"].(types.List))
	resp.Diagnostics.Append(diagCodecs...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Converted inbound codecs", map[string]interface{}{
		"codecs": codecs,
	})

	connection := telnyx.CredentialConnection{
		ConnectionName:                   plan.ConnectionName.ValueString(),
		Username:                         plan.Username.ValueString(),
		Password:                         plan.Password.ValueString(),
		Active:                           plan.Active.ValueBool(),
		AnchorsiteOverride:               plan.AnchorsiteOverride.ValueString(),
		DefaultOnHoldComfortNoiseEnabled: plan.DefaultOnHoldComfortNoiseEnabled.ValueBool(),
		DTMFType:                         plan.DTMFType.ValueString(),
		EncodeContactHeaderEnabled:       plan.EncodeContactHeaderEnabled.ValueBool(),
		OnnetT38PassthroughEnabled:       plan.OnnetT38PassthroughEnabled.ValueBool(),
		MicrosoftTeamsSbc:                plan.MicrosoftTeamsSBC.ValueBool(),
		WebhookEventURL:                  plan.WebhookEventURL.ValueString(),
		WebhookEventFailoverURL:          plan.WebhookEventFailoverURL.ValueString(),
		WebhookAPIVersion:                plan.WebhookAPIVersion.ValueString(),
		WebhookTimeoutSecs:               int(plan.WebhookTimeoutSecs.ValueInt64()),
		RTCPSettings: telnyx.RTCPSettings{
			Port:                rtcpSettingsAttributes["port"].(types.String).ValueString(),
			CaptureEnabled:      rtcpSettingsAttributes["capture_enabled"].(types.Bool).ValueBool(),
			ReportFrequencySecs: int(rtcpSettingsAttributes["report_frequency_secs"].(types.Int64).ValueInt64()),
		},
		Inbound: telnyx.InboundSettings{
			ANINumberFormat:          inboundAttributes["ani_number_format"].(types.String).ValueString(),
			DNISNumberFormat:         inboundAttributes["dnis_number_format"].(types.String).ValueString(),
			Codecs:                   codecs,
			DefaultRoutingMethod:     inboundAttributes["default_routing_method"].(types.String).ValueString(),
			ChannelLimit:             int(inboundAttributes["channel_limit"].(types.Int64).ValueInt64()),
			GenerateRingbackTone:     inboundAttributes["generate_ringback_tone"].(types.Bool).ValueBool(),
			ISUPHeadersEnabled:       inboundAttributes["isup_headers_enabled"].(types.Bool).ValueBool(),
			PRACKEnabled:             inboundAttributes["prack_enabled"].(types.Bool).ValueBool(),
			PrivacyZoneEnabled:       inboundAttributes["privacy_zone_enabled"].(types.Bool).ValueBool(),
			SIPCompactHeadersEnabled: inboundAttributes["sip_compact_headers_enabled"].(types.Bool).ValueBool(),
			Timeout1xxSecs:           int(inboundAttributes["timeout_1xx_secs"].(types.Int64).ValueInt64()),
			Timeout2xxSecs:           int(inboundAttributes["timeout_2xx_secs"].(types.Int64).ValueInt64()),
			ShakenSTIREnabled:        inboundAttributes["shaken_stir_enabled"].(types.Bool).ValueBool(),
		},
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            outboundAttributes["ani_override"].(types.String).ValueString(),
			ANIOverrideType:        outboundAttributes["ani_override_type"].(types.String).ValueString(),
			CallParkingEnabled:     outboundAttributes["call_parking_enabled"].(types.Bool).ValueBool(),
			ChannelLimit:           int(outboundAttributes["channel_limit"].(types.Int64).ValueInt64()),
			GenerateRingbackTone:   outboundAttributes["generate_ringback_tone"].(types.Bool).ValueBool(),
			InstantRingbackEnabled: outboundAttributes["instant_ringback_enabled"].(types.Bool).ValueBool(),
			Localization:           outboundAttributes["localization"].(types.String).ValueString(),
			OutboundVoiceProfileID: outboundAttributes["outbound_voice_profile_id"].(types.String).ValueString(),
			T38ReinviteSource:      outboundAttributes["t38_reinvite_source"].(types.String).ValueString(),
		},
	}

	createdConnection, err := r.client.CreateCredentialConnection(connection)
	if err != nil {
		resp.Diagnostics.AddError("Error creating credential connection", err.Error())
		return
	}

	plan.ID = types.StringValue(createdConnection.ID)

	tflog.Info(ctx, "Created Credential Connection", map[string]interface{}{
		"id": createdConnection.ID,
	})

	setCredentialConnectionState(ctx, &plan, createdConnection)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CredentialConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CredentialConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.client.GetCredentialConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading credential connection", err.Error())
		return
	}

	setCredentialConnectionState(ctx, &state, connection)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *CredentialConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CredentialConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CredentialConnectionResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	codecs, diagCodecs := convertListToStrings(ctx, plan.Inbound.Attributes()["codecs"].(types.List))
	resp.Diagnostics.Append(diagCodecs...)
	if resp.Diagnostics.HasError() {
		return
	}

	connection := telnyx.CredentialConnection{
		ConnectionName:                   plan.ConnectionName.ValueString(),
		Username:                         plan.Username.ValueString(),
		Password:                         plan.Password.ValueString(),
		Active:                           plan.Active.ValueBool(),
		AnchorsiteOverride:               plan.AnchorsiteOverride.ValueString(),
		DefaultOnHoldComfortNoiseEnabled: plan.DefaultOnHoldComfortNoiseEnabled.ValueBool(),
		DTMFType:                         plan.DTMFType.ValueString(),
		EncodeContactHeaderEnabled:       plan.EncodeContactHeaderEnabled.ValueBool(),
		OnnetT38PassthroughEnabled:       plan.OnnetT38PassthroughEnabled.ValueBool(),
		MicrosoftTeamsSbc:                plan.MicrosoftTeamsSBC.ValueBool(),
		WebhookEventURL:                  plan.WebhookEventURL.ValueString(),
		WebhookEventFailoverURL:          plan.WebhookEventFailoverURL.ValueString(),
		WebhookAPIVersion:                plan.WebhookAPIVersion.ValueString(),
		WebhookTimeoutSecs:               int(plan.WebhookTimeoutSecs.ValueInt64()),
		RTCPSettings: telnyx.RTCPSettings{
			Port:                plan.RTCPSettings.Attributes()["port"].(types.String).ValueString(),
			CaptureEnabled:      plan.RTCPSettings.Attributes()["capture_enabled"].(types.Bool).ValueBool(),
			ReportFrequencySecs: int(plan.RTCPSettings.Attributes()["report_frequency_secs"].(types.Int64).ValueInt64()),
		},
		Inbound: telnyx.InboundSettings{
			ANINumberFormat:          plan.Inbound.Attributes()["ani_number_format"].(types.String).ValueString(),
			DNISNumberFormat:         plan.Inbound.Attributes()["dnis_number_format"].(types.String).ValueString(),
			Codecs:                   codecs,
			DefaultRoutingMethod:     plan.Inbound.Attributes()["default_routing_method"].(types.String).ValueString(),
			ChannelLimit:             int(plan.Inbound.Attributes()["channel_limit"].(types.Int64).ValueInt64()),
			GenerateRingbackTone:     plan.Inbound.Attributes()["generate_ringback_tone"].(types.Bool).ValueBool(),
			ISUPHeadersEnabled:       plan.Inbound.Attributes()["isup_headers_enabled"].(types.Bool).ValueBool(),
			PRACKEnabled:             plan.Inbound.Attributes()["prack_enabled"].(types.Bool).ValueBool(),
			PrivacyZoneEnabled:       plan.Inbound.Attributes()["privacy_zone_enabled"].(types.Bool).ValueBool(),
			SIPCompactHeadersEnabled: plan.Inbound.Attributes()["sip_compact_headers_enabled"].(types.Bool).ValueBool(),
			Timeout1xxSecs:           int(plan.Inbound.Attributes()["timeout_1xx_secs"].(types.Int64).ValueInt64()),
			Timeout2xxSecs:           int(plan.Inbound.Attributes()["timeout_2xx_secs"].(types.Int64).ValueInt64()),
			ShakenSTIREnabled:        plan.Inbound.Attributes()["shaken_stir_enabled"].(types.Bool).ValueBool(),
		},
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            plan.Outbound.Attributes()["ani_override"].(types.String).ValueString(),
			ANIOverrideType:        plan.Outbound.Attributes()["ani_override_type"].(types.String).ValueString(),
			CallParkingEnabled:     plan.Outbound.Attributes()["call_parking_enabled"].(types.Bool).ValueBool(),
			ChannelLimit:           int(plan.Outbound.Attributes()["channel_limit"].(types.Int64).ValueInt64()),
			GenerateRingbackTone:   plan.Outbound.Attributes()["generate_ringback_tone"].(types.Bool).ValueBool(),
			InstantRingbackEnabled: plan.Outbound.Attributes()["instant_ringback_enabled"].(types.Bool).ValueBool(),
			Localization:           plan.Outbound.Attributes()["localization"].(types.String).ValueString(),
			OutboundVoiceProfileID: plan.Outbound.Attributes()["outbound_voice_profile_id"].(types.String).ValueString(),
			T38ReinviteSource:      plan.Outbound.Attributes()["t38_reinvite_source"].(types.String).ValueString(),
		},
	}

	// Use state ID in update call
	updatedConnection, err := r.client.UpdateCredentialConnection(state.ID.ValueString(), connection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating credential connection",
			"Could not update credential connection, unexpected error: "+err.Error(),
		)
		return
	}

	setCredentialConnectionState(ctx, &plan, updatedConnection)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CredentialConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CredentialConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCredentialConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting credential connection", err.Error())
	}

	tflog.Info(ctx, "Deleted Credential Connection", map[string]interface{}{"id": state.ID.ValueString()})
}

func (r *CredentialConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setCredentialConnectionState(ctx context.Context, state *CredentialConnectionResourceModel, connection *telnyx.CredentialConnection) {
	state.ID = types.StringValue(connection.ID)
	state.ConnectionName = types.StringValue(connection.ConnectionName)
	state.Username = types.StringValue(connection.Username)
	state.Password = types.StringValue(connection.Password)
	state.Active = types.BoolValue(connection.Active)
	state.AnchorsiteOverride = types.StringValue(connection.AnchorsiteOverride)
	state.DefaultOnHoldComfortNoiseEnabled = types.BoolValue(connection.DefaultOnHoldComfortNoiseEnabled)
	state.DTMFType = types.StringValue(connection.DTMFType)
	state.EncodeContactHeaderEnabled = types.BoolValue(connection.EncodeContactHeaderEnabled)
	state.OnnetT38PassthroughEnabled = types.BoolValue(connection.OnnetT38PassthroughEnabled)
	state.MicrosoftTeamsSBC = types.BoolValue(connection.MicrosoftTeamsSbc)
	state.WebhookEventURL = types.StringValue(connection.WebhookEventURL)
	state.WebhookEventFailoverURL = types.StringValue(connection.WebhookEventFailoverURL)
	state.WebhookAPIVersion = types.StringValue(connection.WebhookAPIVersion)
	state.WebhookTimeoutSecs = types.Int64Value(int64(connection.WebhookTimeoutSecs))

	if connection.RTCPSettings != (telnyx.RTCPSettings{}) {
		state.RTCPSettings = types.ObjectValueMust(
			map[string]attr.Type{
				"port":                  types.StringType,
				"capture_enabled":       types.BoolType,
				"report_frequency_secs": types.Int64Type,
			},
			map[string]attr.Value{
				"port":                  types.StringValue(connection.RTCPSettings.Port),
				"capture_enabled":       types.BoolValue(connection.RTCPSettings.CaptureEnabled),
				"report_frequency_secs": types.Int64Value(int64(connection.RTCPSettings.ReportFrequencySecs)),
			},
		)
	} else {
		state.RTCPSettings = types.ObjectNull(map[string]attr.Type{
			"port":                  types.StringType,
			"capture_enabled":       types.BoolType,
			"report_frequency_secs": types.Int64Type,
		})
	}

	state.Inbound = types.ObjectValueMust(
		map[string]attr.Type{
			"ani_number_format":           types.StringType,
			"dnis_number_format":          types.StringType,
			"codecs":                      types.ListType{ElemType: types.StringType},
			"default_routing_method":      types.StringType,
			"channel_limit":               types.Int64Type,
			"generate_ringback_tone":      types.BoolType,
			"isup_headers_enabled":        types.BoolType,
			"prack_enabled":               types.BoolType,
			"privacy_zone_enabled":        types.BoolType,
			"sip_compact_headers_enabled": types.BoolType,
			"timeout_1xx_secs":            types.Int64Type,
			"timeout_2xx_secs":            types.Int64Type,
			"shaken_stir_enabled":         types.BoolType,
		},
		map[string]attr.Value{
			"ani_number_format":           types.StringValue(connection.Inbound.ANINumberFormat),
			"dnis_number_format":          types.StringValue(connection.Inbound.DNISNumberFormat),
			"codecs":                      convertStringsToList(connection.Inbound.Codecs),
			"default_routing_method":      types.StringValue(connection.Inbound.DefaultRoutingMethod),
			"channel_limit":               types.Int64Value(int64(connection.Inbound.ChannelLimit)),
			"generate_ringback_tone":      types.BoolValue(connection.Inbound.GenerateRingbackTone),
			"isup_headers_enabled":        types.BoolValue(connection.Inbound.ISUPHeadersEnabled),
			"prack_enabled":               types.BoolValue(connection.Inbound.PRACKEnabled),
			"privacy_zone_enabled":        types.BoolValue(connection.Inbound.PrivacyZoneEnabled),
			"sip_compact_headers_enabled": types.BoolValue(connection.Inbound.SIPCompactHeadersEnabled),
			"timeout_1xx_secs":            types.Int64Value(int64(connection.Inbound.Timeout1xxSecs)),
			"timeout_2xx_secs":            types.Int64Value(int64(connection.Inbound.Timeout2xxSecs)),
			"shaken_stir_enabled":         types.BoolValue(connection.Inbound.ShakenSTIREnabled),
		},
	)

	state.Outbound = types.ObjectValueMust(
		map[string]attr.Type{
			"ani_override":              types.StringType,
			"ani_override_type":         types.StringType,
			"call_parking_enabled":      types.BoolType,
			"channel_limit":             types.Int64Type,
			"generate_ringback_tone":    types.BoolType,
			"instant_ringback_enabled":  types.BoolType,
			"localization":              types.StringType,
			"outbound_voice_profile_id": types.StringType,
			"t38_reinvite_source":       types.StringType,
		},
		map[string]attr.Value{
			"ani_override":              types.StringValue(connection.Outbound.ANIOverride),
			"ani_override_type":         types.StringValue(connection.Outbound.ANIOverrideType),
			"call_parking_enabled":      types.BoolValue(connection.Outbound.CallParkingEnabled),
			"channel_limit":             types.Int64Value(int64(connection.Outbound.ChannelLimit)),
			"generate_ringback_tone":    types.BoolValue(connection.Outbound.GenerateRingbackTone),
			"instant_ringback_enabled":  types.BoolValue(connection.Outbound.InstantRingbackEnabled),
			"localization":              types.StringValue(connection.Outbound.Localization),
			"outbound_voice_profile_id": types.StringValue(connection.Outbound.OutboundVoiceProfileID),
			"t38_reinvite_source":       types.StringValue(connection.Outbound.T38ReinviteSource),
		},
	)
}

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	_ resource.Resource                = &FQDNConnectionResource{}
	_ resource.ResourceWithConfigure   = &FQDNConnectionResource{}
	_ resource.ResourceWithImportState = &FQDNConnectionResource{}
)

func NewFQDNConnectionResource() resource.Resource {
	return &FQDNConnectionResource{}
}

type FQDNConnectionResource struct {
	client *telnyx.TelnyxClient
}

type FQDNConnectionResourceModel struct {
	ID                               types.String `tfsdk:"id"`
	ConnectionName                   types.String `tfsdk:"connection_name"`
	Username                         types.String `tfsdk:"username"`
	Password                         types.String `tfsdk:"password"`
	Active                           types.Bool   `tfsdk:"active"`
	AnchorsiteOverride               types.String `tfsdk:"anchorsite_override"`
	TransportProtocol                types.String `tfsdk:"transport_protocol"`
	DefaultOnHoldComfortNoiseEnabled types.Bool   `tfsdk:"default_on_hold_comfort_noise_enabled"`
	DTMFType                         types.String `tfsdk:"dtmf_type"`
	EncodeContactHeaderEnabled       types.Bool   `tfsdk:"encode_contact_header_enabled"`
	EncryptedMedia                   types.String `tfsdk:"encrypted_media"`
	OnnetT38PassthroughEnabled       types.Bool   `tfsdk:"onnet_t38_passthrough_enabled"`
	MicrosoftTeamsSBC                types.Bool   `tfsdk:"microsoft_teams_sbc"`
	WebhookEventURL                  types.String `tfsdk:"webhook_event_url"`
	WebhookEventFailoverURL          types.String `tfsdk:"webhook_event_failover_url"`
	WebhookAPIVersion                types.String `tfsdk:"webhook_api_version"`
	WebhookTimeoutSecs               types.Int64  `tfsdk:"webhook_timeout_secs"`
	RTCPSettings                     types.Object `tfsdk:"rtcp_settings"`
	Inbound                          types.Object `tfsdk:"inbound"`
	Outbound                         types.Object `tfsdk:"outbound"`
	SipUriCallingPreference          types.String `tfsdk:"sip_uri_calling_preference"`
}

func (r *FQDNConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fqdn_connection"
}

func (r *FQDNConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx FQDN Connections",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the FQDN connection",
				Computed:    true,
			},
			"connection_name": schema.StringAttribute{
				Description: "Name of the FQDN connection",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for the FQDN connection",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the FQDN connection",
				Required:    true,
				Sensitive:   true,
			},
			"active": schema.BoolAttribute{
				Description: "Specifies whether the FQDN connection is active or not",
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
			"transport_protocol": schema.StringAttribute{
				Description: "Transport protocol",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("UDP"),
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
			"encrypted_media": schema.StringAttribute{
				Description: "Encrypted media",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("SRTP"),
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
				Default:     stringdefault.StaticString(""),
			},
			"webhook_event_failover_url": schema.StringAttribute{
				Description: "Webhook event failover URL",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"webhook_api_version": schema.StringAttribute{
				Description: "Webhook API version",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("2"),
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
						"ani_number_format":              types.StringType,
						"dnis_number_format":             types.StringType,
						"codecs":                         types.ListType{ElemType: types.StringType},
						"default_routing_method":         types.StringType,
						"channel_limit":                  types.Int64Type,
						"generate_ringback_tone":         types.BoolType,
						"isup_headers_enabled":           types.BoolType,
						"prack_enabled":                  types.BoolType,
						"privacy_zone_enabled":           types.BoolType,
						"sip_compact_headers_enabled":    types.BoolType,
						"sip_region":                     types.StringType,
						"sip_subdomain":                  types.StringType,
						"sip_subdomain_receive_settings": types.StringType,
						"timeout_1xx_secs":               types.Int64Type,
						"timeout_2xx_secs":               types.Int64Type,
						"shaken_stir_enabled":            types.BoolType,
					},
					map[string]attr.Value{
						"ani_number_format":              types.StringValue("E.164-national"),
						"dnis_number_format":             types.StringValue("e164"),
						"codecs":                         types.ListValueMust(types.StringType, []attr.Value{types.StringValue("G722"), types.StringValue("G711U"), types.StringValue("G711A"), types.StringValue("G729"), types.StringValue("OPUS"), types.StringValue("H.264")}),
						"default_routing_method":         types.StringValue("sequential"),
						"channel_limit":                  types.Int64Null(),
						"generate_ringback_tone":         types.BoolNull(),
						"isup_headers_enabled":           types.BoolNull(),
						"prack_enabled":                  types.BoolNull(),
						"privacy_zone_enabled":           types.BoolNull(),
						"sip_compact_headers_enabled":    types.BoolNull(),
						"sip_region":                     types.StringValue("US"),
						"sip_subdomain_receive_settings": types.StringValue(""),
						"sip_subdomain":                  types.StringValue(""),
						"timeout_1xx_secs":               types.Int64Null(),
						"timeout_2xx_secs":               types.Int64Null(),
						"shaken_stir_enabled":            types.BoolNull(),
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
					},
					"generate_ringback_tone": schema.BoolAttribute{
						Description: "Generate ringback tone",
						Optional:    true,
						Computed:    true,
					},
					"isup_headers_enabled": schema.BoolAttribute{
						Description: "ISUP headers enabled",
						Optional:    true,
						Computed:    true,
					},
					"prack_enabled": schema.BoolAttribute{
						Description: "PRACK enabled",
						Optional:    true,
						Computed:    true,
					},
					"privacy_zone_enabled": schema.BoolAttribute{
						Description: "Privacy zone enabled",
						Optional:    true,
						Computed:    true,
					},
					"sip_compact_headers_enabled": schema.BoolAttribute{
						Description: "SIP compact headers enabled",
						Optional:    true,
						Computed:    true,
					},
					"sip_region": schema.StringAttribute{
						Description: "SIP region",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("US"),
					},
					"sip_subdomain": schema.StringAttribute{
						Description: "SIP subdomain",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"sip_subdomain_receive_settings": schema.StringAttribute{
						Description: "SIP subdomain receive settings",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("only_my_connections"),
					},
					"timeout_1xx_secs": schema.Int64Attribute{
						Description: "Timeout for 1xx responses in seconds",
						Optional:    true,
						Computed:    true,
					},
					"timeout_2xx_secs": schema.Int64Attribute{
						Description: "Timeout for 2xx responses in seconds",
						Optional:    true,
						Computed:    true,
					},
					"shaken_stir_enabled": schema.BoolAttribute{
						Description: "SHAKEN/STIR enabled",
						Optional:    true,
						Computed:    true,
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
						"ip_authentication_method":  types.StringType,
						"ip_authentication_token":   types.StringType,
						"localization":              types.StringType,
						"outbound_voice_profile_id": types.StringType,
						"t38_reinvite_source":       types.StringType,
					},
					map[string]attr.Value{
						"ani_override":              types.StringValue(""),
						"ani_override_type":         types.StringValue("always"),
						"call_parking_enabled":      types.BoolValue(false),
						"channel_limit":             types.Int64Null(),
						"generate_ringback_tone":    types.BoolNull(),
						"instant_ringback_enabled":  types.BoolNull(),
						"ip_authentication_method":  types.StringValue("token"),
						"ip_authentication_token":   types.StringNull(),
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
						Default:     stringdefault.StaticString(""),
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
						Default:     booldefault.StaticBool(false),
					},
					"channel_limit": schema.Int64Attribute{
						Description: "Channel limit",
						Optional:    true,
						Computed:    true,
					},
					"generate_ringback_tone": schema.BoolAttribute{
						Description: "Generate ringback tone",
						Optional:    true,
						Computed:    true,
					},
					"instant_ringback_enabled": schema.BoolAttribute{
						Description: "Instant ringback enabled",
						Optional:    true,
						Computed:    true,
					},
					"ip_authentication_method": schema.StringAttribute{
						Description: "IP authentication method",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("token"),
					},
					"ip_authentication_token": schema.StringAttribute{
						Description: "IP authentication token",
						Optional:    true,
						Computed:    true,
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
			"sip_uri_calling_preference": schema.StringAttribute{
				Description: "SIP URI calling preference",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *FQDNConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for FQDNConnectionResource")
	}
}

func (r *FQDNConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FQDNConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating FQDN Connection", map[string]interface{}{
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

	connection := telnyx.FQDNConnection{
		ConnectionName:                   plan.ConnectionName.ValueString(),
		Username:                         telnyx.StringPtr(plan.Username.ValueString()),
		Password:                         telnyx.StringPtr(plan.Password.ValueString()),
		Active:                           plan.Active.ValueBool(),
		AnchorsiteOverride:               plan.AnchorsiteOverride.ValueString(),
		TransportProtocol:                plan.TransportProtocol.ValueString(),
		DefaultOnHoldComfortNoiseEnabled: plan.DefaultOnHoldComfortNoiseEnabled.ValueBool(),
		DTMFType:                         plan.DTMFType.ValueString(),
		EncodeContactHeaderEnabled:       plan.EncodeContactHeaderEnabled.ValueBool(),
		EncryptedMedia:                   nil, // As specified in the Terraform config
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
			ANINumberFormat:             inboundAttributes["ani_number_format"].(types.String).ValueString(),
			DNISNumberFormat:            inboundAttributes["dnis_number_format"].(types.String).ValueString(),
			Codecs:                      codecs,
			DefaultRoutingMethod:        inboundAttributes["default_routing_method"].(types.String).ValueString(),
			ChannelLimit:                getIntPointer(inboundAttributes["channel_limit"].(types.Int64)),
			GenerateRingbackTone:        getBoolPointer(inboundAttributes["generate_ringback_tone"].(types.Bool)),
			ISUPHeadersEnabled:          getBoolPointer(inboundAttributes["isup_headers_enabled"].(types.Bool)),
			PRACKEnabled:                getBoolPointer(inboundAttributes["prack_enabled"].(types.Bool)),
			PrivacyZoneEnabled:          getBoolPointer(inboundAttributes["privacy_zone_enabled"].(types.Bool)),
			SIPCompactHeadersEnabled:    getBoolPointer(inboundAttributes["sip_compact_headers_enabled"].(types.Bool)),
			SIPRegion:                   inboundAttributes["sip_region"].(types.String).ValueString(),
			SIPSubdomain:                inboundAttributes["sip_subdomain"].(types.String).ValueString(),
			SIPSubdomainReceiveSettings: inboundAttributes["sip_subdomain_receive_settings"].(types.String).ValueString(),
			Timeout1xxSecs:              getIntPointer(inboundAttributes["timeout_1xx_secs"].(types.Int64)),
			Timeout2xxSecs:              getIntPointer(inboundAttributes["timeout_2xx_secs"].(types.Int64)),
			ShakenSTIREnabled:           getBoolPointer(inboundAttributes["shaken_stir_enabled"].(types.Bool)),
		},
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            outboundAttributes["ani_override"].(types.String).ValueString(),
			ANIOverrideType:        outboundAttributes["ani_override_type"].(types.String).ValueString(),
			CallParkingEnabled:     getBoolPointer(outboundAttributes["call_parking_enabled"].(types.Bool)),
			ChannelLimit:           getIntPointer(outboundAttributes["channel_limit"].(types.Int64)),
			GenerateRingbackTone:   getBoolPointer(outboundAttributes["generate_ringback_tone"].(types.Bool)),
			InstantRingbackEnabled: getBoolPointer(outboundAttributes["instant_ringback_enabled"].(types.Bool)),
			IPAuthenticationMethod: outboundAttributes["ip_authentication_method"].(types.String).ValueString(),
			IPAuthenticationToken:  getStringPointer(outboundAttributes["ip_authentication_token"].(types.String)),
			Localization:           outboundAttributes["localization"].(types.String).ValueString(),
			OutboundVoiceProfileID: outboundAttributes["outbound_voice_profile_id"].(types.String).ValueString(),
			T38ReinviteSource:      outboundAttributes["t38_reinvite_source"].(types.String).ValueString(),
		},
	}

	createdConnection, err := r.client.CreateFQDNConnection(connection)
	if err != nil {
		resp.Diagnostics.AddError("Error creating FQDN connection", err.Error())
		return
	}

	plan.ID = types.StringValue(createdConnection.ID)

	tflog.Info(ctx, "Created FQDN Connection", map[string]interface{}{
		"id": createdConnection.ID,
	})

	setFQDNConnectionState(ctx, &plan, createdConnection)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FQDNConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FQDNConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.client.GetFQDNConnection(state.ID.ValueString())
	if err == nil {
		setFQDNConnectionState(ctx, &state, connection)
		diags = resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		return
	}
	if telnyxErr, ok := err.(*telnyx.TelnyxError); ok && telnyxErr.IsResourceNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.AddError("Error reading FQDN connection", err.Error())
}

func (r *FQDNConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FQDNConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FQDNConnectionResourceModel
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

	connection := telnyx.FQDNConnection{
		ConnectionName:                   plan.ConnectionName.ValueString(),
		Username:                         telnyx.StringPtr(plan.Username.ValueString()),
		Password:                         telnyx.StringPtr(plan.Password.ValueString()),
		Active:                           plan.Active.ValueBool(),
		AnchorsiteOverride:               plan.AnchorsiteOverride.ValueString(),
		TransportProtocol:                plan.TransportProtocol.ValueString(),
		DefaultOnHoldComfortNoiseEnabled: plan.DefaultOnHoldComfortNoiseEnabled.ValueBool(),
		DTMFType:                         plan.DTMFType.ValueString(),
		EncodeContactHeaderEnabled:       plan.EncodeContactHeaderEnabled.ValueBool(),
		EncryptedMedia:                   nil, // As specified in the Terraform config
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
			ANINumberFormat:             plan.Inbound.Attributes()["ani_number_format"].(types.String).ValueString(),
			DNISNumberFormat:            plan.Inbound.Attributes()["dnis_number_format"].(types.String).ValueString(),
			Codecs:                      codecs,
			DefaultRoutingMethod:        plan.Inbound.Attributes()["default_routing_method"].(types.String).ValueString(),
			ChannelLimit:                getIntPointer(plan.Inbound.Attributes()["channel_limit"].(types.Int64)),
			GenerateRingbackTone:        getBoolPointer(plan.Inbound.Attributes()["generate_ringback_tone"].(types.Bool)),
			ISUPHeadersEnabled:          getBoolPointer(plan.Inbound.Attributes()["isup_headers_enabled"].(types.Bool)),
			PRACKEnabled:                getBoolPointer(plan.Inbound.Attributes()["prack_enabled"].(types.Bool)),
			PrivacyZoneEnabled:          getBoolPointer(plan.Inbound.Attributes()["privacy_zone_enabled"].(types.Bool)),
			SIPCompactHeadersEnabled:    getBoolPointer(plan.Inbound.Attributes()["sip_compact_headers_enabled"].(types.Bool)),
			SIPRegion:                   plan.Inbound.Attributes()["sip_region"].(types.String).ValueString(),
			SIPSubdomain:                plan.Inbound.Attributes()["sip_subdomain"].(types.String).ValueString(),
			SIPSubdomainReceiveSettings: plan.Inbound.Attributes()["sip_subdomain_receive_settings"].(types.String).ValueString(),
			Timeout1xxSecs:              getIntPointer(plan.Inbound.Attributes()["timeout_1xx_secs"].(types.Int64)),
			Timeout2xxSecs:              getIntPointer(plan.Inbound.Attributes()["timeout_2xx_secs"].(types.Int64)),
			ShakenSTIREnabled:           getBoolPointer(plan.Inbound.Attributes()["shaken_stir_enabled"].(types.Bool)),
		},
		Outbound: telnyx.OutboundSettings{
			ANIOverride:            plan.Outbound.Attributes()["ani_override"].(types.String).ValueString(),
			ANIOverrideType:        plan.Outbound.Attributes()["ani_override_type"].(types.String).ValueString(),
			CallParkingEnabled:     getBoolPointer(plan.Outbound.Attributes()["call_parking_enabled"].(types.Bool)),
			ChannelLimit:           getIntPointer(plan.Outbound.Attributes()["channel_limit"].(types.Int64)),
			GenerateRingbackTone:   getBoolPointer(plan.Outbound.Attributes()["generate_ringback_tone"].(types.Bool)),
			InstantRingbackEnabled: getBoolPointer(plan.Outbound.Attributes()["instant_ringback_enabled"].(types.Bool)),
			IPAuthenticationMethod: plan.Outbound.Attributes()["ip_authentication_method"].(types.String).ValueString(),
			IPAuthenticationToken:  getStringPointer(plan.Outbound.Attributes()["ip_authentication_token"].(types.String)),
			Localization:           plan.Outbound.Attributes()["localization"].(types.String).ValueString(),
			OutboundVoiceProfileID: plan.Outbound.Attributes()["outbound_voice_profile_id"].(types.String).ValueString(),
			T38ReinviteSource:      plan.Outbound.Attributes()["t38_reinvite_source"].(types.String).ValueString(),
		},
	}

	updatedConnection, err := r.client.UpdateFQDNConnection(state.ID.ValueString(), connection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating FQDN connection",
			"Could not update FQDN connection, unexpected error: "+err.Error(),
		)
		return
	}

	setFQDNConnectionState(ctx, &plan, updatedConnection)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FQDNConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FQDNConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFQDNConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting FQDN connection", err.Error())
	}

	tflog.Info(ctx, "Deleted FQDN Connection", map[string]interface{}{"id": state.ID.ValueString()})
}

func (r *FQDNConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setFQDNConnectionState(ctx context.Context, state *FQDNConnectionResourceModel, connection *telnyx.FQDNConnection) {
	state.ID = types.StringValue(connection.ID)
	state.ConnectionName = types.StringValue(connection.ConnectionName)
	state.Username = types.StringValue(*connection.Username)
	state.Password = types.StringValue(*connection.Password)
	state.Active = types.BoolValue(connection.Active)
	state.AnchorsiteOverride = types.StringValue(connection.AnchorsiteOverride)
	state.TransportProtocol = types.StringValue(connection.TransportProtocol)
	state.DefaultOnHoldComfortNoiseEnabled = types.BoolValue(connection.DefaultOnHoldComfortNoiseEnabled)
	state.DTMFType = types.StringValue(connection.DTMFType)
	state.EncodeContactHeaderEnabled = types.BoolValue(connection.EncodeContactHeaderEnabled)
	// state.EncryptedMedia = types.StringValue("") // Conforming to null value in Terraform config
	state.OnnetT38PassthroughEnabled = types.BoolValue(connection.OnnetT38PassthroughEnabled)
	// state.MicrosoftTeamsSBC = types.BoolValue(connection.MicrosoftTeamsSbc)
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
			"ani_number_format":              types.StringType,
			"dnis_number_format":             types.StringType,
			"codecs":                         types.ListType{ElemType: types.StringType},
			"default_routing_method":         types.StringType,
			"channel_limit":                  types.Int64Type,
			"generate_ringback_tone":         types.BoolType,
			"isup_headers_enabled":           types.BoolType,
			"prack_enabled":                  types.BoolType,
			"privacy_zone_enabled":           types.BoolType,
			"sip_compact_headers_enabled":    types.BoolType,
			"sip_region":                     types.StringType,
			"sip_subdomain":                  types.StringType,
			"sip_subdomain_receive_settings": types.StringType,
			"timeout_1xx_secs":               types.Int64Type,
			"timeout_2xx_secs":               types.Int64Type,
			"shaken_stir_enabled":            types.BoolType,
		},
		map[string]attr.Value{
			"ani_number_format":              types.StringValue(connection.Inbound.ANINumberFormat),
			"dnis_number_format":             types.StringValue(connection.Inbound.DNISNumberFormat),
			"codecs":                         convertStringsToList(connection.Inbound.Codecs),
			"default_routing_method":         types.StringValue(connection.Inbound.DefaultRoutingMethod),
			"channel_limit":                  types.Int64Value(getInt64(connection.Inbound.ChannelLimit)),
			"generate_ringback_tone":         types.BoolValue(getBool(connection.Inbound.GenerateRingbackTone)),
			"isup_headers_enabled":           types.BoolValue(getBool(connection.Inbound.ISUPHeadersEnabled)),
			"prack_enabled":                  types.BoolValue(getBool(connection.Inbound.PRACKEnabled)),
			"privacy_zone_enabled":           types.BoolValue(getBool(connection.Inbound.PrivacyZoneEnabled)),
			"sip_compact_headers_enabled":    types.BoolValue(getBool(connection.Inbound.SIPCompactHeadersEnabled)),
			"sip_region":                     types.StringValue(connection.Inbound.SIPRegion),
			"sip_subdomain":                  types.StringValue(connection.Inbound.SIPSubdomain),
			"sip_subdomain_receive_settings": types.StringValue(connection.Inbound.SIPSubdomainReceiveSettings),
			"timeout_1xx_secs":               types.Int64Value(getInt64(connection.Inbound.Timeout1xxSecs)),
			"timeout_2xx_secs":               types.Int64Value(getInt64(connection.Inbound.Timeout2xxSecs)),
			"shaken_stir_enabled":            types.BoolValue(getBool(connection.Inbound.ShakenSTIREnabled)),
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
			"ip_authentication_method":  types.StringType,
			"ip_authentication_token":   types.StringType,
			"localization":              types.StringType,
			"outbound_voice_profile_id": types.StringType,
			"t38_reinvite_source":       types.StringType,
		},
		map[string]attr.Value{
			"ani_override":              types.StringValue(connection.Outbound.ANIOverride),
			"ani_override_type":         types.StringValue(connection.Outbound.ANIOverrideType),
			"call_parking_enabled":      types.BoolValue(getBool(connection.Outbound.CallParkingEnabled)),
			"channel_limit":             types.Int64Value(getInt64(connection.Outbound.ChannelLimit)),
			"generate_ringback_tone":    types.BoolValue(getBool(connection.Outbound.GenerateRingbackTone)),
			"instant_ringback_enabled":  types.BoolValue(getBool(connection.Outbound.InstantRingbackEnabled)),
			"ip_authentication_method":  types.StringValue(connection.Outbound.IPAuthenticationMethod),
			"ip_authentication_token":   types.StringValue(getString(connection.Outbound.IPAuthenticationToken)),
			"localization":              types.StringValue(connection.Outbound.Localization),
			"outbound_voice_profile_id": types.StringValue(connection.Outbound.OutboundVoiceProfileID),
			"t38_reinvite_source":       types.StringValue(connection.Outbound.T38ReinviteSource),
		},
	)

	state.SipUriCallingPreference = types.StringValue("")
}

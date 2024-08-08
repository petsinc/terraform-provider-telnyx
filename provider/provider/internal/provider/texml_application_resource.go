package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

var (
	_ resource.Resource              = &TeXMLApplicationResource{}
	_ resource.ResourceWithConfigure = &TeXMLApplicationResource{}
)

func NewTeXMLApplicationResource() resource.Resource {
	return &TeXMLApplicationResource{}
}

type TeXMLApplicationResource struct {
	client *telnyx.TelnyxClient
}

type TeXMLApplicationResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	FriendlyName            types.String `tfsdk:"friendly_name"`
	Active                  types.Bool   `tfsdk:"active"`
	AnchorsiteOverride      types.String `tfsdk:"anchorsite_override"`
	DTMFType                types.String `tfsdk:"dtmf_type"`
	FirstCommandTimeout     types.Bool   `tfsdk:"first_command_timeout"`
	FirstCommandTimeoutSecs types.Int64  `tfsdk:"first_command_timeout_secs"`
	VoiceURL                types.String `tfsdk:"voice_url"`
	VoiceFallbackURL        types.String `tfsdk:"voice_fallback_url"`
	VoiceMethod             types.String `tfsdk:"voice_method"`
	StatusCallback          types.String `tfsdk:"status_callback"`
	StatusCallbackMethod    types.String `tfsdk:"status_callback_method"`
	Inbound                 types.Object `tfsdk:"inbound"`
	Outbound                types.Object `tfsdk:"outbound"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
}

func (r *TeXMLApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_texml_application"
}

func (r *TeXMLApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx TeXML applications",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the TeXML application",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"friendly_name": schema.StringAttribute{
				Description: "User-assigned name for the application",
				Required:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Specifies whether the connection can be used",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"anchorsite_override": schema.StringAttribute{
				Description: "Anchorsite Override",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Latency"),
			},
			"dtmf_type": schema.StringAttribute{
				Description: "DTMF Type",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("RFC 2833"),
			},
			"first_command_timeout": schema.BoolAttribute{
				Description: "Specifies whether calls should hang up after timing out",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"first_command_timeout_secs": schema.Int64Attribute{
				Description: "Specifies how many seconds to wait before timing out a dial command",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
			},
			"voice_url": schema.StringAttribute{
				Description: "URL to deliver XML Translator webhooks",
				Required:    true,
			},
			"voice_fallback_url": schema.StringAttribute{
				Description: "Fallback URL to deliver XML Translator webhooks if the primary URL fails",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"voice_method": schema.StringAttribute{
				Description: "HTTP request method for voice webhooks",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("post"),
			},
			"status_callback": schema.StringAttribute{
				Description: "URL for status callback",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"status_callback_method": schema.StringAttribute{
				Description: "HTTP request method for status callback",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("post"),
			},
			"inbound": schema.SingleNestedAttribute{
				Description: "Inbound settings for the TeXML application",
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
						"ani_number_format":              types.StringValue(""),
						"dnis_number_format":             types.StringValue(""),
						"codecs":                         types.ListValueMust(types.StringType, []attr.Value{types.StringValue("G722"), types.StringValue("G711U"), types.StringValue("G711A"), types.StringValue("G729"), types.StringValue("OPUS"), types.StringValue("H.264")}),
						"default_routing_method":         types.StringValue("sequential"),
						"channel_limit":                  types.Int64Null(),
						"generate_ringback_tone":         types.BoolNull(),
						"isup_headers_enabled":           types.BoolNull(),
						"prack_enabled":                  types.BoolNull(),
						"privacy_zone_enabled":           types.BoolNull(),
						"sip_compact_headers_enabled":    types.BoolNull(),
						"sip_region":                     types.StringValue(""),
						"sip_subdomain":                  types.StringValue(""),
						"sip_subdomain_receive_settings": types.StringValue("only_my_connections"),
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
						Default:     stringdefault.StaticString(""),
					},
					"dnis_number_format": schema.StringAttribute{
						Description: "DNIS number format",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
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
						Default:     stringdefault.StaticString(""),
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
						Default:     stringdefault.StaticString(""),
					},
					"sip_subdomain": schema.StringAttribute{
						Description: "Subdomain for receiving inbound calls",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"sip_subdomain_receive_settings": schema.StringAttribute{
						Description: "Receive calls from specified endpoints",
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
				Description: "Outbound settings for the TeXML application",
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
						"ani_override_type":         types.StringValue(""),
						"call_parking_enabled":      types.BoolValue(false),
						"channel_limit":             types.Int64Null(),
						"generate_ringback_tone":    types.BoolNull(),
						"instant_ringback_enabled":  types.BoolNull(),
						"ip_authentication_method":  types.StringValue(""),
						"ip_authentication_token":   types.StringNull(),
						"localization":              types.StringValue(""),
						"outbound_voice_profile_id": types.StringValue(""),
						"t38_reinvite_source":       types.StringValue(""),
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
						Default:     stringdefault.StaticString(""),
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
						Default:     stringdefault.StaticString(""),
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
						Default:     stringdefault.StaticString(""),
					},
					"outbound_voice_profile_id": schema.StringAttribute{
						Description: "Outbound voice profile ID",
						Optional:    true,
					},
					"t38_reinvite_source": schema.StringAttribute{
						Description: "T38 reinvite source",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Creation time of the TeXML application",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update time of the TeXML application",
				Computed:    true,
			},
		},
	}
}

func (r *TeXMLApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for TeXMLApplicationResource")
	}
}

func (r *TeXMLApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeXMLApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inboundAttributes := plan.Inbound.Attributes()
	outboundAttributes := plan.Outbound.Attributes()

	codecs, diagCodecs := convertListToStrings(ctx, inboundAttributes["codecs"].(types.List))
	resp.Diagnostics.Append(diagCodecs...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationRequest := telnyx.TeXMLApplicationRequest{
		FriendlyName:            plan.FriendlyName.ValueString(),
		Active:                  plan.Active.ValueBool(),
		AnchorsiteOverride:      plan.AnchorsiteOverride.ValueString(),
		DTMFType:                plan.DTMFType.ValueString(),
		FirstCommandTimeout:     plan.FirstCommandTimeout.ValueBool(),
		FirstCommandTimeoutSecs: int(plan.FirstCommandTimeoutSecs.ValueInt64()),
		VoiceURL:                plan.VoiceURL.ValueString(),
		VoiceFallbackURL:        plan.VoiceFallbackURL.ValueString(),
		VoiceMethod:             plan.VoiceMethod.ValueString(),
		StatusCallback:          plan.StatusCallback.ValueString(),
		StatusCallbackMethod:    plan.StatusCallbackMethod.ValueString(),
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

	application, err := r.client.CreateTeXMLApplication(applicationRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating TeXML application", err.Error())
		return
	}

	// Set state based on response from the API
	setStateFromTeXMLApplicationResponse(&plan, application)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *TeXMLApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeXMLApplicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.GetTeXMLApplication(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading TeXML application", err.Error())
		return
	}
	// Update state based on response from the API
	setStateFromTeXMLApplicationResponse(&state, application)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *TeXMLApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeXMLApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inboundAttributes := plan.Inbound.Attributes()
	outboundAttributes := plan.Outbound.Attributes()

	codecs, diagCodecs := convertListToStrings(ctx, inboundAttributes["codecs"].(types.List))
	resp.Diagnostics.Append(diagCodecs...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationRequest := telnyx.TeXMLApplicationRequest{
		FriendlyName:            plan.FriendlyName.ValueString(),
		Active:                  plan.Active.ValueBool(),
		AnchorsiteOverride:      plan.AnchorsiteOverride.ValueString(),
		DTMFType:                plan.DTMFType.ValueString(),
		FirstCommandTimeout:     plan.FirstCommandTimeout.ValueBool(),
		FirstCommandTimeoutSecs: int(plan.FirstCommandTimeoutSecs.ValueInt64()),
		VoiceURL:                plan.VoiceURL.ValueString(),
		VoiceFallbackURL:        plan.VoiceFallbackURL.ValueString(),
		VoiceMethod:             plan.VoiceMethod.ValueString(),
		StatusCallback:          plan.StatusCallback.ValueString(),
		StatusCallbackMethod:    plan.StatusCallbackMethod.ValueString(),
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

	application, err := r.client.UpdateTeXMLApplication(plan.ID.ValueString(), applicationRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating TeXML application", err.Error())
		return
	}

	// Update state based on response from the API
	setStateFromTeXMLApplicationResponse(&plan, application)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *TeXMLApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeXMLApplicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeXMLApplication(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting TeXML application", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

func setStateFromTeXMLApplicationResponse(state *TeXMLApplicationResourceModel, application *telnyx.TeXMLApplication) {
	state.ID = types.StringValue(application.ID)
	state.FriendlyName = types.StringValue(application.FriendlyName)
	state.Active = types.BoolValue(application.Active)
	state.AnchorsiteOverride = types.StringValue(application.AnchorsiteOverride)
	state.DTMFType = types.StringValue(application.DTMFType)
	state.FirstCommandTimeout = types.BoolValue(application.FirstCommandTimeout)
	state.FirstCommandTimeoutSecs = types.Int64Value(int64(application.FirstCommandTimeoutSecs))
	state.VoiceURL = types.StringValue(application.VoiceURL)
	state.VoiceFallbackURL = types.StringValue(application.VoiceFallbackURL)
	state.VoiceMethod = types.StringValue(application.VoiceMethod)
	state.StatusCallback = types.StringValue(application.StatusCallback)
	state.StatusCallbackMethod = types.StringValue(application.StatusCallbackMethod)

	codecsList := convertStringsToList(application.Inbound.Codecs)

	state.Inbound, _ = types.ObjectValue(map[string]attr.Type{
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
	}, map[string]attr.Value{
		"ani_number_format":              types.StringValue(application.Inbound.ANINumberFormat),
		"dnis_number_format":             types.StringValue(application.Inbound.DNISNumberFormat),
		"codecs":                         codecsList,
		"default_routing_method":         types.StringValue(application.Inbound.DefaultRoutingMethod),
		"channel_limit":                  types.Int64Value(getInt64(application.Inbound.ChannelLimit)),
		"generate_ringback_tone":         types.BoolValue(getBool(application.Inbound.GenerateRingbackTone)),
		"isup_headers_enabled":           types.BoolValue(getBool(application.Inbound.ISUPHeadersEnabled)),
		"prack_enabled":                  types.BoolValue(getBool(application.Inbound.PRACKEnabled)),
		"privacy_zone_enabled":           types.BoolValue(getBool(application.Inbound.PrivacyZoneEnabled)),
		"sip_compact_headers_enabled":    types.BoolValue(getBool(application.Inbound.SIPCompactHeadersEnabled)),
		"sip_region":                     types.StringValue(application.Inbound.SIPRegion),
		"sip_subdomain":                  types.StringValue(application.Inbound.SIPSubdomain),
		"sip_subdomain_receive_settings": types.StringValue(application.Inbound.SIPSubdomainReceiveSettings),
		"timeout_1xx_secs":               types.Int64Value(getInt64(application.Inbound.Timeout1xxSecs)),
		"timeout_2xx_secs":               types.Int64Value(getInt64(application.Inbound.Timeout2xxSecs)),
		"shaken_stir_enabled":            types.BoolValue(getBool(application.Inbound.ShakenSTIREnabled)),
	})

	state.Outbound, _ = types.ObjectValue(map[string]attr.Type{
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
	}, map[string]attr.Value{
		"ani_override":              types.StringValue(application.Outbound.ANIOverride),
		"ani_override_type":         types.StringValue(application.Outbound.ANIOverrideType),
		"call_parking_enabled":      types.BoolValue(getBool(application.Outbound.CallParkingEnabled)),
		"channel_limit":             types.Int64Value(getInt64(application.Outbound.ChannelLimit)),
		"generate_ringback_tone":    types.BoolValue(getBool(application.Outbound.GenerateRingbackTone)),
		"instant_ringback_enabled":  types.BoolValue(getBool(application.Outbound.InstantRingbackEnabled)),
		"ip_authentication_method":  types.StringValue(application.Outbound.IPAuthenticationMethod),
		"ip_authentication_token":   types.StringValue(getString(application.Outbound.IPAuthenticationToken)),
		"localization":              types.StringValue(application.Outbound.Localization),
		"outbound_voice_profile_id": types.StringValue(application.Outbound.OutboundVoiceProfileID),
		"t38_reinvite_source":       types.StringValue(application.Outbound.T38ReinviteSource),
	})
	state.CreatedAt = types.StringValue(application.CreatedAt.String())
	state.UpdatedAt = types.StringValue(application.UpdatedAt.String())
}

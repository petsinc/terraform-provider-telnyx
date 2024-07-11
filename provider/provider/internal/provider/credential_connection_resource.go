package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

type CredentialConnectionResource struct {
	client *telnyx.TelnyxClient
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
	}
}

func NewCredentialConnectionResource() resource.Resource {
	return &CredentialConnectionResource{}
}

type CredentialConnectionResourceModel struct {
	ID                               types.String                                `tfsdk:"id"`
	ConnectionName                   types.String                                `tfsdk:"connection_name"`
	Username                         types.String                                `tfsdk:"username"`
	Password                         types.String                                `tfsdk:"password"`
	Active                           types.Bool                                  `tfsdk:"active"`
	AnchorsiteOverride               types.String                                `tfsdk:"anchorsite_override"`
	DefaultOnHoldComfortNoiseEnabled types.Bool                                  `tfsdk:"default_on_hold_comfort_noise_enabled"`
	DTMFType                         types.String                                `tfsdk:"dtmf_type"`
	EncodeContactHeaderEnabled       types.Bool                                  `tfsdk:"encode_contact_header_enabled"`
	OnnetT38PassthroughEnabled       types.Bool                                  `tfsdk:"onnet_t38_passthrough_enabled"`
	MicrosoftTeamsSBC                types.Bool                                  `tfsdk:"microsoft_teams_sbc"`
	WebhookEventURL                  types.String                                `tfsdk:"webhook_event_url"`
	WebhookEventFailoverURL          types.String                                `tfsdk:"webhook_event_failover_url"`
	WebhookAPIVersion                types.String                                `tfsdk:"webhook_api_version"`
	WebhookTimeoutSecs               types.Int64                                 `tfsdk:"webhook_timeout_secs"`
	RTCPSettings                     *RTCPSettingsModel                          `tfsdk:"rtcp_settings"`
	Inbound                          *CredentialConnectionInboundSettingsModel   `tfsdk:"inbound"`
	Outbound                         *CredentialConnectionOutboundSettingsModel  `tfsdk:"outbound"`
}

type RTCPSettingsModel struct {
	Port                types.String `tfsdk:"port"`
	CaptureEnabled      types.Bool   `tfsdk:"capture_enabled"`
	ReportFrequencySecs types.Int64  `tfsdk:"report_frequency_secs"`
}

type CredentialConnectionInboundSettingsModel struct {
	ANINumberFormat          types.String `tfsdk:"ani_number_format"`
	DNISNumberFormat         types.String `tfsdk:"dnis_number_format"`
	Codecs                   types.List   `tfsdk:"codecs"`
	DefaultRoutingMethod     types.String `tfsdk:"default_routing_method"`
	ChannelLimit             types.Int64  `tfsdk:"channel_limit"`
	GenerateRingbackTone     types.Bool   `tfsdk:"generate_ringback_tone"`
	ISUPHeadersEnabled       types.Bool   `tfsdk:"isup_headers_enabled"`
	PRACKEnabled             types.Bool   `tfsdk:"prack_enabled"`
	PrivacyZoneEnabled       types.Bool   `tfsdk:"privacy_zone_enabled"`
	SIPCompactHeadersEnabled types.Bool   `tfsdk:"sip_compact_headers_enabled"`
	Timeout1xxSecs           types.Int64  `tfsdk:"timeout_1xx_secs"`
	Timeout2xxSecs           types.Int64  `tfsdk:"timeout_2xx_secs"`
	ShakenSTIREnabled        types.Bool   `tfsdk:"shaken_stir_enabled"`
}

type CredentialConnectionOutboundSettingsModel struct {
	ANIOverride            types.String `tfsdk:"ani_override"`
	ANIOverrideType        types.String `tfsdk:"ani_override_type"`
	CallParkingEnabled     types.Bool   `tfsdk:"call_parking_enabled"`
	ChannelLimit           types.Int64  `tfsdk:"channel_limit"`
	GenerateRingbackTone   types.Bool   `tfsdk:"generate_ringback_tone"`
	InstantRingbackEnabled types.Bool   `tfsdk:"instant_ringback_enabled"`
	Localization           types.String `tfsdk:"localization"`
	OutboundVoiceProfileID types.String `tfsdk:"outbound_voice_profile_id"`
	T38ReinviteSource      types.String `tfsdk:"t38_reinvite_source"`
}

func (r *CredentialConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_connection"
}

func (r *CredentialConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"connection_name": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"password": schema.StringAttribute{
				Required: true,
			},
			"active": schema.BoolAttribute{
				Required: true,
			},
			"anchorsite_override": schema.StringAttribute{
				Optional: true,
			},
			"default_on_hold_comfort_noise_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"dtmf_type": schema.StringAttribute{
				Optional: true,
			},
			"encode_contact_header_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"onnet_t38_passthrough_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"microsoft_teams_sbc": schema.BoolAttribute{
				Optional: true,
			},
			"webhook_event_url": schema.StringAttribute{
				Optional: true,
			},
			"webhook_event_failover_url": schema.StringAttribute{
				Optional: true,
			},
			"webhook_api_version": schema.StringAttribute{
				Optional: true,
			},
			"webhook_timeout_secs": schema.Int64Attribute{
				Optional: true,
			},
			"rtcp_settings": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"port": schema.StringAttribute{
						Optional: true,
					},
					"capture_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"report_frequency_secs": schema.Int64Attribute{
						Optional: true,
					},
				},
			},
			"inbound": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"ani_number_format": schema.StringAttribute{
						Optional: true,
					},
					"dnis_number_format": schema.StringAttribute{
						Optional: true,
					},
					"codecs": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"default_routing_method": schema.StringAttribute{
						Optional: true,
					},
					"channel_limit": schema.Int64Attribute{
						Optional: true,
					},
					"generate_ringback_tone": schema.BoolAttribute{
						Optional: true,
					},
					"isup_headers_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"prack_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"privacy_zone_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"sip_compact_headers_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"timeout_1xx_secs": schema.Int64Attribute{
						Optional: true,
					},
					"timeout_2xx_secs": schema.Int64Attribute{
						Optional: true,
					},
					"shaken_stir_enabled": schema.BoolAttribute{
						Optional: true,
					},
				},
			},
			"outbound": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"ani_override": schema.StringAttribute{
						Optional: true,
					},
					"ani_override_type": schema.StringAttribute{
						Optional: true,
					},
					"call_parking_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"channel_limit": schema.Int64Attribute{
						Optional: true,
					},
					"generate_ringback_tone": schema.BoolAttribute{
						Optional: true,
					},
					"instant_ringback_enabled": schema.BoolAttribute{
						Optional: true,
					},
					"localization": schema.StringAttribute{
						Optional: true,
					},
					"outbound_voice_profile_id": schema.StringAttribute{
						Optional: true,
					},
					"t38_reinvite_source": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *CredentialConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CredentialConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client

	codecs, diags := convertListToStrings(ctx, plan.Inbound.Codecs)
	resp.Diagnostics.Append(diags...)
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
			Port:                plan.RTCPSettings.Port.ValueString(),
			CaptureEnabled:      plan.RTCPSettings.CaptureEnabled.ValueBool(),
			ReportFrequencySecs: int(plan.RTCPSettings.ReportFrequencySecs.ValueInt64()),
		},
		Inbound: telnyx.CredentialConnectionInboundSettings{
			ANINumberFormat:          plan.Inbound.ANINumberFormat.ValueString(),
			DNISNumberFormat:         plan.Inbound.DNISNumberFormat.ValueString(),
			Codecs:                   codecs,
			DefaultRoutingMethod:     plan.Inbound.DefaultRoutingMethod.ValueString(),
			ChannelLimit:             int(plan.Inbound.ChannelLimit.ValueInt64()),
			GenerateRingbackTone:     plan.Inbound.GenerateRingbackTone.ValueBool(),
			ISUPHeadersEnabled:       plan.Inbound.ISUPHeadersEnabled.ValueBool(),
			PRACKEnabled:             plan.Inbound.PRACKEnabled.ValueBool(),
			PrivacyZoneEnabled:       plan.Inbound.PrivacyZoneEnabled.ValueBool(),
			SIPCompactHeadersEnabled: plan.Inbound.SIPCompactHeadersEnabled.ValueBool(),
			Timeout1xxSecs:           int(plan.Inbound.Timeout1xxSecs.ValueInt64()),
			Timeout2xxSecs:           int(plan.Inbound.Timeout2xxSecs.ValueInt64()),
			ShakenSTIREnabled:        plan.Inbound.ShakenSTIREnabled.ValueBool(),
		},
		Outbound: telnyx.CredentialConnectionOutboundSettings{
			ANIOverride:            plan.Outbound.ANIOverride.ValueString(),
			ANIOverrideType:        plan.Outbound.ANIOverrideType.ValueString(),
			CallParkingEnabled:     plan.Outbound.CallParkingEnabled.ValueBool(),
			ChannelLimit:           int(plan.Outbound.ChannelLimit.ValueInt64()),
			GenerateRingbackTone:   plan.Outbound.GenerateRingbackTone.ValueBool(),
			InstantRingbackEnabled: plan.Outbound.InstantRingbackEnabled.ValueBool(),
			Localization:           plan.Outbound.Localization.ValueString(),
			OutboundVoiceProfileID: plan.Outbound.OutboundVoiceProfileID.ValueString(),
			T38ReinviteSource:      plan.Outbound.T38ReinviteSource.ValueString(),
		},
	}

	createdConnection, err := client.CreateCredentialConnection(connection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating credential connection",
			"Could not create credential connection, unexpected error: "+err.Error(),
		)
		return
	}

	// fmt.Printf("\n\n--- Created Connection ---\n%+v\n\n", createdConnection)

	setCredentialConnectionState(ctx, &plan, createdConnection)

	// fmt.Printf("\n\n--- State After Create ---\n%+v\n\n", plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// fmt.Printf("\n\n--- Error setting state after create ---\n%+v\n\n", resp.Diagnostics)
		return
	}
}

func (r *CredentialConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CredentialConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// fmt.Printf("\n\n--- Error getting state ---\n%+v\n\n", resp.Diagnostics)
		return
	}

	// fmt.Printf("\n\n--- State Before Read ---\n%+v\n\n", state)

	client := r.client

	connection, err := client.GetCredentialConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading credential connection",
			"Could not read credential connection, unexpected error: "+err.Error(),
		)
		// fmt.Printf("\n\n--- Error reading credential connection ---\n%s\n\n", err.Error())
		return
	}

	// fmt.Printf("\n\n--- Response from API ---\n%+v\n\n", connection)

	setCredentialConnectionState(ctx, &state, connection)

	// fmt.Printf("\n\n--- State After Setting Connection ---\n%+v\n\n", state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// fmt.Printf("\n\n--- Error setting state after read ---\n%+v\n\n", resp.Diagnostics)
		return
	}

	// fmt.Printf("\n\n--- Final State After Read ---\n%+v\n\n", state)
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

	client := r.client

	// Convert necessary fields from plan
	codecs, diags := convertListToStrings(ctx, plan.Inbound.Codecs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare connection data
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
			Port:                plan.RTCPSettings.Port.ValueString(),
			CaptureEnabled:      plan.RTCPSettings.CaptureEnabled.ValueBool(),
			ReportFrequencySecs: int(plan.RTCPSettings.ReportFrequencySecs.ValueInt64()),
		},
		Inbound: telnyx.CredentialConnectionInboundSettings{
			ANINumberFormat:          plan.Inbound.ANINumberFormat.ValueString(),
			DNISNumberFormat:         plan.Inbound.DNISNumberFormat.ValueString(),
			Codecs:                   codecs,
			DefaultRoutingMethod:     plan.Inbound.DefaultRoutingMethod.ValueString(),
			ChannelLimit:             int(plan.Inbound.ChannelLimit.ValueInt64()),
			GenerateRingbackTone:     plan.Inbound.GenerateRingbackTone.ValueBool(),
			ISUPHeadersEnabled:       plan.Inbound.ISUPHeadersEnabled.ValueBool(),
			PRACKEnabled:             plan.Inbound.PRACKEnabled.ValueBool(),
			PrivacyZoneEnabled:       plan.Inbound.PrivacyZoneEnabled.ValueBool(),
			SIPCompactHeadersEnabled: plan.Inbound.SIPCompactHeadersEnabled.ValueBool(),
			Timeout1xxSecs:           int(plan.Inbound.Timeout1xxSecs.ValueInt64()),
			Timeout2xxSecs:           int(plan.Inbound.Timeout2xxSecs.ValueInt64()),
			ShakenSTIREnabled:        plan.Inbound.ShakenSTIREnabled.ValueBool(),
		},
		Outbound: telnyx.CredentialConnectionOutboundSettings{
			ANIOverride:            plan.Outbound.ANIOverride.ValueString(),
			ANIOverrideType:        plan.Outbound.ANIOverrideType.ValueString(),
			CallParkingEnabled:     plan.Outbound.CallParkingEnabled.ValueBool(),
			ChannelLimit:           int(plan.Outbound.ChannelLimit.ValueInt64()),
			GenerateRingbackTone:   plan.Outbound.GenerateRingbackTone.ValueBool(),
			InstantRingbackEnabled: plan.Outbound.InstantRingbackEnabled.ValueBool(),
			Localization:           plan.Outbound.Localization.ValueString(),
			OutboundVoiceProfileID: plan.Outbound.OutboundVoiceProfileID.ValueString(),
			T38ReinviteSource:      plan.Outbound.T38ReinviteSource.ValueString(),
		},
	}

	// Use state ID in update call
	updatedConnection, err := client.UpdateCredentialConnection(state.ID.ValueString(), connection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating credential connection",
			"Could not update credential connection, unexpected error: "+err.Error(),
		)
		return
	}

	// fmt.Printf("\n\n--- Updated Connection ---\n%+v\n\n", updatedConnection)

	setCredentialConnectionState(ctx, &plan, updatedConnection)

	// fmt.Printf("\n\n--- State After Update ---\n%+v\n\n", plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// fmt.Printf("\n\n--- Error setting state after update ---\n%+v\n\n", resp.Diagnostics)
		return
	}
}

func (r *CredentialConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CredentialConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client

	err := client.DeleteCredentialConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting credential connection",
			"Could not delete credential connection, unexpected error: "+err.Error(),
		)
		return
	}
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
		state.RTCPSettings = &RTCPSettingsModel{
			Port:                types.StringValue(connection.RTCPSettings.Port),
			CaptureEnabled:      types.BoolValue(connection.RTCPSettings.CaptureEnabled),
			ReportFrequencySecs: types.Int64Value(int64(connection.RTCPSettings.ReportFrequencySecs)),
		}
	} else {
		state.RTCPSettings = nil
	}

	// Set inbound settings with defaults
	state.Inbound = &CredentialConnectionInboundSettingsModel{
		ANINumberFormat:          types.StringValue(connection.Inbound.ANINumberFormat),
		DNISNumberFormat:         types.StringValue(connection.Inbound.DNISNumberFormat),
		Codecs:                   convertStringsToList(connection.Inbound.Codecs),
		DefaultRoutingMethod:     types.StringValue(connection.Inbound.DefaultRoutingMethod),
		ChannelLimit:             types.Int64Value(int64(connection.Inbound.ChannelLimit)),
		GenerateRingbackTone:     types.BoolValue(connection.Inbound.GenerateRingbackTone),
		ISUPHeadersEnabled:       types.BoolValue(connection.Inbound.ISUPHeadersEnabled),
		PRACKEnabled:             types.BoolValue(connection.Inbound.PRACKEnabled),
		PrivacyZoneEnabled:       types.BoolValue(connection.Inbound.PrivacyZoneEnabled),
		SIPCompactHeadersEnabled: types.BoolValue(connection.Inbound.SIPCompactHeadersEnabled),
		Timeout1xxSecs:           types.Int64Value(int64(connection.Inbound.Timeout1xxSecs)),
		Timeout2xxSecs:           types.Int64Value(int64(connection.Inbound.Timeout2xxSecs)),
		ShakenSTIREnabled:        types.BoolValue(connection.Inbound.ShakenSTIREnabled),
	}

	// Set outbound settings with defaults
	state.Outbound = &CredentialConnectionOutboundSettingsModel{
		ANIOverride:            types.StringValue(connection.Outbound.ANIOverride),
		ANIOverrideType:        types.StringValue(connection.Outbound.ANIOverrideType),
		CallParkingEnabled:     types.BoolValue(connection.Outbound.CallParkingEnabled),
		ChannelLimit:           types.Int64Value(int64(connection.Outbound.ChannelLimit)),
		GenerateRingbackTone:   types.BoolValue(connection.Outbound.GenerateRingbackTone),
		InstantRingbackEnabled: types.BoolValue(connection.Outbound.InstantRingbackEnabled),
		Localization:           types.StringValue(connection.Outbound.Localization),
		OutboundVoiceProfileID: types.StringValue(connection.Outbound.OutboundVoiceProfileID),
		T38ReinviteSource:      types.StringValue(connection.Outbound.T38ReinviteSource),
	}

	// Log the state for debugging
	// fmt.Printf("\n\n--- State Set in setCredentialConnectionState ---\n%+v\n\n", state)
}

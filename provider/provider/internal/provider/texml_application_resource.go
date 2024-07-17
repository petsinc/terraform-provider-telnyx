package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

var (
	_ resource.Resource = &TeXMLApplicationResource{}
)

func NewTeXMLApplicationResource() resource.Resource {
	return &TeXMLApplicationResource{}
}

type TeXMLApplicationResource struct {
	client *telnyx.TelnyxClient
}

type TeXMLApplicationResourceModel struct {
	ID                      types.String               `tfsdk:"id"`
	FriendlyName            types.String               `tfsdk:"friendly_name"`
	Active                  types.Bool                 `tfsdk:"active"`
	AnchorsiteOverride      types.String               `tfsdk:"anchorsite_override"`
	DTMFType                types.String               `tfsdk:"dtmf_type"`
	FirstCommandTimeout     types.Bool                 `tfsdk:"first_command_timeout"`
	FirstCommandTimeoutSecs types.Int64                `tfsdk:"first_command_timeout_secs"`
	VoiceURL                types.String               `tfsdk:"voice_url"`
	VoiceFallbackURL        types.String               `tfsdk:"voice_fallback_url"`
	VoiceMethod             types.String               `tfsdk:"voice_method"`
	StatusCallback          types.String               `tfsdk:"status_callback"`
	StatusCallbackMethod    types.String               `tfsdk:"status_callback_method"`
	Inbound                 InboundTeXMLSettingsModel  `tfsdk:"inbound"`
	Outbound                OutboundTeXMLSettingsModel `tfsdk:"outbound"`
	CreatedAt               types.String               `tfsdk:"created_at"`
	UpdatedAt               types.String               `tfsdk:"updated_at"`
}

type InboundTeXMLSettingsModel struct {
	ChannelLimit                types.Int64  `tfsdk:"channel_limit"`
	ShakenStirEnabled           types.Bool   `tfsdk:"shaken_stir_enabled"`
	SIPSubdomain                types.String `tfsdk:"sip_subdomain"`
	SIPSubdomainReceiveSettings types.String `tfsdk:"sip_subdomain_receive_settings"`
}

type OutboundTeXMLSettingsModel struct {
	ChannelLimit           types.Int64  `tfsdk:"channel_limit"`
	OutboundVoiceProfileID types.String `tfsdk:"outbound_voice_profile_id"`
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
			},
			"anchorsite_override": schema.StringAttribute{
				Description: "Anchorsite Override",
				Optional:    true,
				Computed:    true,
			},
			"dtmf_type": schema.StringAttribute{
				Description: "DTMF Type",
				Optional:    true,
				Computed:    true,
			},
			"first_command_timeout": schema.BoolAttribute{
				Description: "Specifies whether calls should hangup after timing out",
				Optional:    true,
				Computed:    true,
			},
			"first_command_timeout_secs": schema.Int64Attribute{
				Description: "Specifies how many seconds to wait before timing out a dial command",
				Optional:    true,
				Computed:    true,
			},
			"voice_url": schema.StringAttribute{
				Description: "URL to deliver XML Translator webhooks",
				Required:    true,
			},
			"voice_fallback_url": schema.StringAttribute{
				Description: "Fallback URL to deliver XML Translator webhooks if the primary URL fails",
				Optional:    true,
				Computed:    true,
			},
			"voice_method": schema.StringAttribute{
				Description: "HTTP request method for voice webhooks",
				Optional:    true,
				Computed:    true,
			},
			"status_callback": schema.StringAttribute{
				Description: "URL for status callback",
				Optional:    true,
				Computed:    true,
			},
			"status_callback_method": schema.StringAttribute{
				Description: "HTTP request method for status callback",
				Optional:    true,
				Computed:    true,
			},
			"inbound": schema.SingleNestedAttribute{
				Description: "Inbound settings for the TeXML application",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"channel_limit": schema.Int64Attribute{
						Description: "Limits the total number of inbound calls",
						Optional:    true,
						Computed:    true,
					},
					"shaken_stir_enabled": schema.BoolAttribute{
						Description: "Enables Shaken/Stir data for inbound calls",
						Optional:    true,
						Computed:    true,
					},
					"sip_subdomain": schema.StringAttribute{
						Description: "Subdomain for receiving inbound calls",
						Optional:    true,
						Computed:    true,
					},
					"sip_subdomain_receive_settings": schema.StringAttribute{
						Description: "Receive calls from specified endpoints",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"outbound": schema.SingleNestedAttribute{
				Description: "Outbound settings for the TeXML application",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"channel_limit": schema.Int64Attribute{
						Description: "Limits the total number of outbound calls",
						Optional:    true,
						Computed:    true,
					},
					"outbound_voice_profile_id": schema.StringAttribute{
						Description: "Associated outbound voice profile ID",
						Optional:    true,
						Computed:    true,
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
		Inbound: telnyx.InboundTeXMLSettings{
			ChannelLimit:                int(plan.Inbound.ChannelLimit.ValueInt64()),
			ShakenStirEnabled:           plan.Inbound.ShakenStirEnabled.ValueBool(),
			SIPSubdomain:                plan.Inbound.SIPSubdomain.ValueString(),
			SIPSubdomainReceiveSettings: plan.Inbound.SIPSubdomainReceiveSettings.ValueString(),
		},
		Outbound: telnyx.OutboundTeXMLSettings{
			ChannelLimit:           int(plan.Outbound.ChannelLimit.ValueInt64()),
			OutboundVoiceProfileID: plan.Outbound.OutboundVoiceProfileID.ValueString(),
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
		Inbound: telnyx.InboundTeXMLSettings{
			ChannelLimit:                int(plan.Inbound.ChannelLimit.ValueInt64()),
			ShakenStirEnabled:           plan.Inbound.ShakenStirEnabled.ValueBool(),
			SIPSubdomain:                plan.Inbound.SIPSubdomain.ValueString(),
			SIPSubdomainReceiveSettings: plan.Inbound.SIPSubdomainReceiveSettings.ValueString(),
		},
		Outbound: telnyx.OutboundTeXMLSettings{
			ChannelLimit:           int(plan.Outbound.ChannelLimit.ValueInt64()),
			OutboundVoiceProfileID: plan.Outbound.OutboundVoiceProfileID.ValueString(),
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
	state.Inbound = InboundTeXMLSettingsModel{
		ChannelLimit:                types.Int64Value(int64(application.Inbound.ChannelLimit)),
		ShakenStirEnabled:           types.BoolValue(application.Inbound.ShakenStirEnabled),
		SIPSubdomain:                types.StringValue(application.Inbound.SIPSubdomain),
		SIPSubdomainReceiveSettings: types.StringValue(application.Inbound.SIPSubdomainReceiveSettings),
	}
	state.Outbound = OutboundTeXMLSettingsModel{
		ChannelLimit:           types.Int64Value(int64(application.Outbound.ChannelLimit)),
		OutboundVoiceProfileID: types.StringValue(application.Outbound.OutboundVoiceProfileID),
	}
	state.CreatedAt = types.StringValue(application.CreatedAt.String())
	state.UpdatedAt = types.StringValue(application.UpdatedAt.String())
}

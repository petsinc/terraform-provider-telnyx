package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

var (
	_ resource.Resource              = &CallControlApplicationResource{}
	_ resource.ResourceWithConfigure = &CallControlApplicationResource{}
)

func NewCallControlApplicationResource() resource.Resource {
	return &CallControlApplicationResource{}
}

type CallControlApplicationResource struct {
	client *telnyx.TelnyxClient
}

type CallControlApplicationResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Active                  types.Bool   `tfsdk:"active"`
	AnchorsiteOverride      types.String `tfsdk:"anchorsite_override"`
	ApplicationName         types.String `tfsdk:"application_name"`
	DTMFType                types.String `tfsdk:"dtmf_type"`
	FirstCommandTimeout     types.Bool   `tfsdk:"first_command_timeout"`
	FirstCommandTimeoutSecs types.Int64  `tfsdk:"first_command_timeout_secs"`
	Inbound                 types.Object `tfsdk:"inbound"`
	Outbound                types.Object `tfsdk:"outbound"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
	WebhookAPIVersion       types.String `tfsdk:"webhook_api_version"`
	WebhookEventFailoverURL types.String `tfsdk:"webhook_event_failover_url"`
	WebhookEventURL         types.String `tfsdk:"webhook_event_url"`
	WebhookTimeoutSecs      types.Int64  `tfsdk:"webhook_timeout_secs"`
}

func (r *CallControlApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_call_control_application"
}

func (r *CallControlApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx Call Control Applications",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Call Control Application",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_name": schema.StringAttribute{
				Description: "User-assigned name for the application",
				Required:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Specifies whether the application is active",
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
				Description: "How many seconds to wait before timing out a dial command",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
			},
			"inbound": schema.SingleNestedAttribute{
				Description: "Inbound settings for the call control application",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"channel_limit": schema.Int64Attribute{
						Description: "Channel limit",
						Optional:    true,
						Computed:    true,
					},
					"shaken_stir_enabled": schema.BoolAttribute{
						Description: "shaken sir enabled",
						Optional:    true,
						Computed:    true,
					},
					"sip_subdomain": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"sip_subdomain_receive_settings": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
				},
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"channel_limit":                  types.Int64Type,
						"shaken_stir_enabled":            types.BoolType,
						"sip_subdomain":                  types.StringType,
						"sip_subdomain_receive_settings": types.StringType,
					},
					map[string]attr.Value{
						"channel_limit":                  types.Int64Null(),
						"shaken_stir_enabled":            types.BoolNull(),
						"sip_subdomain":                  types.StringValue(""),
						"sip_subdomain_receive_settings": types.StringValue(""),
					},
				)),
			},
			"outbound": schema.SingleNestedAttribute{
				Description: "Outbound settings for the call control application",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"channel_limit": schema.Int64Attribute{
						Description: "Channel limit",
						Optional:    true,
						Computed:    true,
					},
					"outbound_voice_profile_id": schema.StringAttribute{
						Description: "Outbound voice profile ID",
						Optional:    true,
					},
				},
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"channel_limit":             types.Int64Type,
						"outbound_voice_profile_id": types.StringType,
					},
					map[string]attr.Value{
						"channel_limit":             types.Int64Null(),
						"outbound_voice_profile_id": types.StringValue(""),
					},
				)),
			},
			"webhook_api_version": schema.StringAttribute{
				Description: "Webhook API version",
				Optional:    true,
			},
			"webhook_event_failover_url": schema.StringAttribute{
				Description: "Webhook event failover URL",
				Optional:    true,
			},
			"webhook_event_url": schema.StringAttribute{
				Description: "The URL where webhooks related to this connection will be sent. Must include a scheme, such as 'https'.",
				Required:    true,
			},
			"webhook_timeout_secs": schema.Int64Attribute{
				Description: "Webhook timeout in seconds",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the application was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the application was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *CallControlApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CallControlApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CallControlApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inboundAttributes := plan.Inbound.Attributes()
	outboundAttributes := plan.Outbound.Attributes()

	request := telnyx.CallControlApplicationRequest{
		Active:                  plan.Active.ValueBool(),
		AnchorsiteOverride:      plan.AnchorsiteOverride.ValueString(),
		ApplicationName:         plan.ApplicationName.ValueString(),
		DTMFType:                plan.DTMFType.ValueString(),
		FirstCommandTimeout:     plan.FirstCommandTimeout.ValueBool(),
		FirstCommandTimeoutSecs: int(plan.FirstCommandTimeoutSecs.ValueInt64()),
		WebhookAPIVersion:       plan.WebhookAPIVersion.ValueString(),
		WebhookEventFailoverURL: plan.WebhookEventFailoverURL.ValueString(),
		WebhookEventURL:         plan.WebhookEventURL.ValueString(),
		WebhookTimeoutSecs:      int(plan.WebhookTimeoutSecs.ValueInt64()),
		Inbound: telnyx.CallControlInboundSettings{
			ChannelLimit:                getIntPointer(inboundAttributes["channel_limit"].(types.Int64)),
			ShakenSTIREnabled:           getBoolPointer(inboundAttributes["shaken_stir_enabled"].(types.Bool)),
			SIPSubdomain:                inboundAttributes["sip_subdomain"].(types.String).ValueString(),
			SIPSubdomainReceiveSettings: inboundAttributes["sip_subdomain_receive_settings"].(types.String).ValueString(),
		},
		Outbound: telnyx.CallControlOutboundSettings{
			ChannelLimit:           getIntPointer(outboundAttributes["channel_limit"].(types.Int64)),
			OutboundVoiceProfileID: outboundAttributes["outbound_voice_profile_id"].(types.String).ValueString(),
		},
	}
	app, err := r.client.CreateCallControlApplication(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Call Control Application", err.Error())
		return
	}
	setStateResponse(&plan, app)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CallControlApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CallControlApplicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	app, err := r.client.GetCallControlApplication(state.ID.ValueString())
	if err == nil {
		setStateResponse(&state, app)
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}
	if telnyxErr, ok := err.(*telnyx.TelnyxError); ok && telnyxErr.IsResourceNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.AddError("Error reading Call Control Application", err.Error())
}

func (r *CallControlApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CallControlApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inboundAttributes := plan.Inbound.Attributes()
	outboundAttributes := plan.Outbound.Attributes()

	request := telnyx.CallControlApplicationRequest{
		Active:                  plan.Active.ValueBool(),
		AnchorsiteOverride:      plan.AnchorsiteOverride.ValueString(),
		ApplicationName:         plan.ApplicationName.ValueString(),
		DTMFType:                plan.DTMFType.ValueString(),
		FirstCommandTimeout:     plan.FirstCommandTimeout.ValueBool(),
		FirstCommandTimeoutSecs: int(plan.FirstCommandTimeoutSecs.ValueInt64()),
		WebhookAPIVersion:       plan.WebhookAPIVersion.ValueString(),
		WebhookEventFailoverURL: plan.WebhookEventFailoverURL.ValueString(),
		WebhookEventURL:         plan.WebhookEventURL.ValueString(),
		WebhookTimeoutSecs:      int(plan.WebhookTimeoutSecs.ValueInt64()),
		Inbound: telnyx.CallControlInboundSettings{
			ChannelLimit:                getIntPointer(inboundAttributes["channel_limit"].(types.Int64)),
			ShakenSTIREnabled:           getBoolPointer(inboundAttributes["shaken_stir_enabled"].(types.Bool)),
			SIPSubdomain:                inboundAttributes["sip_subdomain"].(types.String).ValueString(),
			SIPSubdomainReceiveSettings: inboundAttributes["sip_subdomain_receive_settings"].(types.String).ValueString(),
		},
		Outbound: telnyx.CallControlOutboundSettings{
			ChannelLimit:           getIntPointer(outboundAttributes["channel_limit"].(types.Int64)),
			OutboundVoiceProfileID: outboundAttributes["outbound_voice_profile_id"].(types.String).ValueString(),
		},
	}
	app, err := r.client.UpdateCallControlApplication(plan.ID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Call Control Application", err.Error())
		return
	}
	setStateResponse(&plan, app)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CallControlApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CallControlApplicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteCallControlApplication(state.ID.ValueString())
	if err == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if telnyxErr, ok := err.(*telnyx.TelnyxError); ok && telnyxErr.IsResourceNotFound() {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.AddError("Error deleting Call Control Application", err.Error())
}

func flattenInboundSettings(inbound telnyx.CallControlInboundSettings) types.Object {
	obj, _ := types.ObjectValue(map[string]attr.Type{
		"channel_limit":                  types.Int64Type,
		"sip_subdomain":                  types.StringType,
		"sip_subdomain_receive_settings": types.StringType,
		"shaken_stir_enabled":            types.BoolType,
	}, map[string]attr.Value{
		"channel_limit":                  types.Int64Value(getInt64(inbound.ChannelLimit)),
		"sip_subdomain":                  types.StringValue(inbound.SIPSubdomain),
		"sip_subdomain_receive_settings": types.StringValue(inbound.SIPSubdomainReceiveSettings),
		"shaken_stir_enabled":            types.BoolValue(getBool(inbound.ShakenSTIREnabled)),
	})
	return obj
}

func flattenOutboundSettings(outbound telnyx.CallControlOutboundSettings) types.Object {
	obj, _ := types.ObjectValue(map[string]attr.Type{
		"channel_limit":             types.Int64Type,
		"outbound_voice_profile_id": types.StringType,
	}, map[string]attr.Value{
		"channel_limit":             types.Int64Value(getInt64(outbound.ChannelLimit)),
		"outbound_voice_profile_id": types.StringValue(outbound.OutboundVoiceProfileID),
	})
	return obj
}

func setStateResponse(state *CallControlApplicationResourceModel, application *telnyx.CallControlApplication) {
	state.ID = types.StringValue(application.ID)
	state.ApplicationName = types.StringValue(application.ApplicationName)
	state.Active = types.BoolValue(application.Active)
	state.AnchorsiteOverride = types.StringValue(application.AnchorsiteOverride)
	state.DTMFType = types.StringValue(application.DTMFType)
	state.FirstCommandTimeout = types.BoolValue(application.FirstCommandTimeout)
	state.FirstCommandTimeoutSecs = types.Int64Value(int64(application.FirstCommandTimeoutSecs))

	// Fix for webhook_event_failover_url
	if application.WebhookEventFailoverURL == "" {
		state.WebhookEventFailoverURL = types.StringNull()
	} else {
		state.WebhookEventFailoverURL = types.StringValue(application.WebhookEventFailoverURL)
	}

	// Fix for webhook_timeout_secs
	if application.WebhookTimeoutSecs == 0 {
		state.WebhookTimeoutSecs = types.Int64Null()
	} else {
		state.WebhookTimeoutSecs = types.Int64Value(int64(application.WebhookTimeoutSecs))
	}

	state.WebhookAPIVersion = types.StringValue(application.WebhookAPIVersion)
	state.WebhookEventURL = types.StringValue(application.WebhookEventURL)

	state.Inbound = flattenInboundSettings(application.Inbound)
	state.Outbound = flattenOutboundSettings(application.Outbound)

	state.CreatedAt = types.StringValue(application.CreatedAt.String())
	state.UpdatedAt = types.StringValue(application.UpdatedAt.String())
}

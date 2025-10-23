package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	_ resource.Resource                = &OutboundVoiceProfileResource{}
	_ resource.ResourceWithConfigure   = &OutboundVoiceProfileResource{}
	_ resource.ResourceWithImportState = &OutboundVoiceProfileResource{}
)

func NewOutboundVoiceProfileResource() resource.Resource {
	return &OutboundVoiceProfileResource{}
}

type OutboundVoiceProfileResource struct {
	client *telnyx.TelnyxClient
}

type OutboundVoiceProfileResourceModel struct {
	ID                      types.String  `tfsdk:"id"`
	Name                    types.String  `tfsdk:"name"`
	BillingGroupID          types.String  `tfsdk:"billing_group_id"`
	ConnectionsCount        types.Int64   `tfsdk:"connections_count"`
	TrafficType             types.String  `tfsdk:"traffic_type"`
	ServicePlan             types.String  `tfsdk:"service_plan"`
	ConcurrentCallLimit     types.Int64   `tfsdk:"concurrent_call_limit"`
	Enabled                 types.Bool    `tfsdk:"enabled"`
	Tags                    types.List    `tfsdk:"tags"`
	UsagePaymentMethod      types.String  `tfsdk:"usage_payment_method"`
	WhitelistedDestinations types.List    `tfsdk:"whitelisted_destinations"`
	MaxDestinationRate      types.Float64 `tfsdk:"max_destination_rate"`
	DailySpendLimit         types.String  `tfsdk:"daily_spend_limit"`
	DailySpendLimitEnabled  types.Bool    `tfsdk:"daily_spend_limit_enabled"`
	CallRecording           types.Object  `tfsdk:"call_recording"`
}

func (r *OutboundVoiceProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_voice_profile"
}

func (r *OutboundVoiceProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx Outbound Voice Profiles",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the outbound voice profile",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the outbound voice profile",
				Required:    true,
			},
			"billing_group_id": schema.StringAttribute{
				Description: "Billing group ID associated with the profile",
				Required:    true,
			},
			"traffic_type": schema.StringAttribute{
				Description: "Type of traffic",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("conversational"),
			},
			"service_plan": schema.StringAttribute{
				Description: "Service plan",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("global"),
			},
			"concurrent_call_limit": schema.Int64Attribute{
				Description: "Concurrent call limit",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Is the profile enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"tags": schema.ListAttribute{
				Description: "Tags for the profile",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"usage_payment_method": schema.StringAttribute{
				Description: "Usage payment method",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("rate-deck"),
			},
			"whitelisted_destinations": schema.ListAttribute{
				Description: "Whitelisted destinations",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("US"), types.StringValue("CA")})),
			},
			"max_destination_rate": schema.Float64Attribute{
				Description: "Max destination rate",
				Optional:    true,
				Computed:    true,
			},
			"daily_spend_limit": schema.StringAttribute{
				Description: "Daily spend limit",
				Optional:    true,
				Computed:    true,
			},
			"daily_spend_limit_enabled": schema.BoolAttribute{
				Description: "Is daily spend limit enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"call_recording": schema.SingleNestedAttribute{
				Description: "Call recording settings",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"type":                 types.StringType,
						"caller_phone_numbers": types.ListType{ElemType: types.StringType},
						"channels":             types.StringType,
						"format":               types.StringType,
					},
					map[string]attr.Value{
						"type":                 types.StringValue("none"),
						"caller_phone_numbers": types.ListValueMust(types.StringType, []attr.Value{}),
						"channels":             types.StringValue("single"),
						"format":               types.StringValue("wav"),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Call recording type",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("none"),
					},
					"caller_phone_numbers": schema.ListAttribute{
						Description: "Caller phone numbers for recording",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
					"channels": schema.StringAttribute{
						Description: "Recording channels",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("single"),
					},
					"format": schema.StringAttribute{
						Description: "Recording format",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("wav"),
					},
				},
			},
		},
	}
}

func (r *OutboundVoiceProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for OutboundVoiceProfileResource")
	}
}

func (r *OutboundVoiceProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OutboundVoiceProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, diagTags := convertListToStrings(ctx, plan.Tags)
	resp.Diagnostics.Append(diagTags...)
	if resp.Diagnostics.HasError() {
		return
	}

	whitelistedDestinations, diagWD := convertListToStrings(ctx, plan.WhitelistedDestinations)
	resp.Diagnostics.Append(diagWD...)
	if resp.Diagnostics.HasError() {
		return
	}

	callRecordingAttributes := plan.CallRecording.Attributes()

	callerPhoneNumbers, diagCP := convertListToStrings(ctx, callRecordingAttributes["caller_phone_numbers"].(types.List))
	resp.Diagnostics.Append(diagCP...)
	if resp.Diagnostics.HasError() {
		return
	}

	var concurrentCallLimitPointer *int
	if !plan.ConcurrentCallLimit.IsNull() && plan.ConcurrentCallLimit.ValueInt64() != 0 {
		value := int(plan.ConcurrentCallLimit.ValueInt64())
		concurrentCallLimitPointer = &value
	}

	profile, err := r.client.CreateOutboundVoiceProfile(telnyx.OutboundVoiceProfile{
		Name:                    plan.Name.ValueString(),
		BillingGroupID:          plan.BillingGroupID.ValueString(),
		ConnectionsCount:        plan.ConnectionsCount.ValueInt64(),
		TrafficType:             plan.TrafficType.ValueString(),
		ServicePlan:             plan.ServicePlan.ValueString(),
		ConcurrentCallLimit:     concurrentCallLimitPointer,
		Enabled:                 plan.Enabled.ValueBool(),
		Tags:                    tags,
		UsagePaymentMethod:      plan.UsagePaymentMethod.ValueString(),
		WhitelistedDestinations: whitelistedDestinations,
		MaxDestinationRate:      getFloat64Pointer(plan.MaxDestinationRate),
		DailySpendLimit:         getStringPointer(plan.DailySpendLimit),
		DailySpendLimitEnabled:  plan.DailySpendLimitEnabled.ValueBool(),
		CallRecording: telnyx.CallRecording{
			Type:               callRecordingAttributes["type"].(types.String).ValueString(),
			CallerPhoneNumbers: callerPhoneNumbers,
			Channels:           callRecordingAttributes["channels"].(types.String).ValueString(),
			Format:             callRecordingAttributes["format"].(types.String).ValueString(),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Error creating outbound voice profile", err.Error())
		return
	}

	if profile.DailySpendLimit == nil {
		plan.DailySpendLimit = types.StringNull()
	} else {
		plan.DailySpendLimit = types.StringValue(*profile.DailySpendLimit)
	}

	if profile.MaxDestinationRate == nil {
		plan.MaxDestinationRate = types.Float64Null()
	} else {
		plan.MaxDestinationRate = types.Float64Value(*profile.MaxDestinationRate)
	}

	if profile.ConcurrentCallLimit == nil {
		plan.ConcurrentCallLimit = types.Int64Null()
	} else {
		plan.ConcurrentCallLimit = types.Int64Value(int64(*profile.ConcurrentCallLimit))
	}

	plan.ID = types.StringValue(profile.ID)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *OutboundVoiceProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OutboundVoiceProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.GetOutboundVoiceProfile(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound voice profile", err.Error())
		return
	}

	state.Name = types.StringValue(profile.Name)
	state.BillingGroupID = types.StringValue(profile.BillingGroupID)
	state.TrafficType = types.StringValue(profile.TrafficType)
	state.ServicePlan = types.StringValue(profile.ServicePlan)

	if profile.ConcurrentCallLimit == nil {
		state.ConcurrentCallLimit = types.Int64Null()
	} else {
		state.ConcurrentCallLimit = types.Int64Value(int64(*profile.ConcurrentCallLimit))
	}

	state.Enabled = types.BoolValue(profile.Enabled)
	state.Tags = convertStringsToList(profile.Tags)
	state.UsagePaymentMethod = types.StringValue(profile.UsagePaymentMethod)
	state.WhitelistedDestinations = convertStringsToList(profile.WhitelistedDestinations)

	if profile.MaxDestinationRate == nil {
		state.MaxDestinationRate = types.Float64Null()
	} else {
		state.MaxDestinationRate = types.Float64Value(*profile.MaxDestinationRate)
	}

	if profile.DailySpendLimit == nil {
		state.DailySpendLimit = types.StringValue("")
	} else {
		state.DailySpendLimit = types.StringValue(*profile.DailySpendLimit)
	}

	state.DailySpendLimitEnabled = types.BoolValue(profile.DailySpendLimitEnabled)

	state.CallRecording, diags = types.ObjectValue(map[string]attr.Type{
		"type":                 types.StringType,
		"caller_phone_numbers": types.ListType{ElemType: types.StringType},
		"channels":             types.StringType,
		"format":               types.StringType,
	}, map[string]attr.Value{
		"type":                 types.StringValue(profile.CallRecording.Type),
		"caller_phone_numbers": convertStringsToList(profile.CallRecording.CallerPhoneNumbers),
		"channels":             types.StringValue(profile.CallRecording.Channels),
		"format":               types.StringValue(profile.CallRecording.Format),
	})
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *OutboundVoiceProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OutboundVoiceProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, diagTags := convertListToStrings(ctx, plan.Tags)
	resp.Diagnostics.Append(diagTags...)
	if resp.Diagnostics.HasError() {
		return
	}

	whitelistedDestinations, diagWD := convertListToStrings(ctx, plan.WhitelistedDestinations)
	resp.Diagnostics.Append(diagWD...)
	if resp.Diagnostics.HasError() {
		return
	}

	callRecordingAttributes := plan.CallRecording.Attributes()

	callerPhoneNumbers, diagCP := convertListToStrings(ctx, callRecordingAttributes["caller_phone_numbers"].(types.List))
	resp.Diagnostics.Append(diagCP...)
	if resp.Diagnostics.HasError() {
		return
	}

	var concurrentCallLimitPointer *int
	if !plan.ConcurrentCallLimit.IsNull() && plan.ConcurrentCallLimit.ValueInt64() != 0 {
		value := int(plan.ConcurrentCallLimit.ValueInt64())
		concurrentCallLimitPointer = &value
	}

	var dailySpendLimitPointer *string
	if plan.DailySpendLimit.IsNull() || plan.DailySpendLimit.ValueString() == "" {
		dailySpendLimitPointer = nil
	} else {
		dailySpendLimitPointer = getStringPointer(plan.DailySpendLimit)
	}

	var maxDestinationRatePointer *float64
	if plan.MaxDestinationRate.IsNull() {
		maxDestinationRatePointer = nil
	} else {
		maxDestinationRatePointer = getFloat64Pointer(plan.MaxDestinationRate)
	}

	profile, err := r.client.UpdateOutboundVoiceProfile(plan.ID.ValueString(), telnyx.OutboundVoiceProfile{
		Name:                    plan.Name.ValueString(),
		BillingGroupID:          plan.BillingGroupID.ValueString(),
		ConnectionsCount:        plan.ConnectionsCount.ValueInt64(),
		TrafficType:             plan.TrafficType.ValueString(),
		ServicePlan:             plan.ServicePlan.ValueString(),
		ConcurrentCallLimit:     concurrentCallLimitPointer,
		Enabled:                 plan.Enabled.ValueBool(),
		Tags:                    tags,
		UsagePaymentMethod:      plan.UsagePaymentMethod.ValueString(),
		WhitelistedDestinations: whitelistedDestinations,
		MaxDestinationRate:      maxDestinationRatePointer,
		DailySpendLimit:         dailySpendLimitPointer,
		DailySpendLimitEnabled:  plan.DailySpendLimitEnabled.ValueBool(),
		CallRecording: telnyx.CallRecording{
			Type:               callRecordingAttributes["type"].(types.String).ValueString(),
			CallerPhoneNumbers: callerPhoneNumbers,
			Channels:           callRecordingAttributes["channels"].(types.String).ValueString(),
			Format:             callRecordingAttributes["format"].(types.String).ValueString(),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Error updating outbound voice profile", err.Error())
		return
	}

	if profile.DailySpendLimit == nil {
		plan.DailySpendLimit = types.StringNull()
	} else {
		plan.DailySpendLimit = types.StringValue(*profile.DailySpendLimit)
	}

	if profile.MaxDestinationRate == nil {
		plan.MaxDestinationRate = types.Float64Null()
	} else {
		plan.MaxDestinationRate = types.Float64Value(*profile.MaxDestinationRate)
	}

	if profile.ConcurrentCallLimit == nil {
		plan.ConcurrentCallLimit = types.Int64Null()
	} else {
		plan.ConcurrentCallLimit = types.Int64Value(int64(*profile.ConcurrentCallLimit))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *OutboundVoiceProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OutboundVoiceProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOutboundVoiceProfile(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting outbound voice profile", err.Error())
	}
}

func (r *OutboundVoiceProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

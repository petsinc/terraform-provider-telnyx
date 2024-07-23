package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
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
				Default:     int64default.StaticInt64(10),
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
				Default:     float64default.StaticFloat64(10.0),
			},
			"daily_spend_limit": schema.StringAttribute{
				Description: "Daily spend limit",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("100.00"),
			},
			"daily_spend_limit_enabled": schema.BoolAttribute{
				Description: "Is daily spend limit enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"call_recording": schema.SingleNestedAttribute{
				Description: "Call recording settings",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"call_recording_type":                 types.StringType,
						"call_recording_caller_phone_numbers": types.ListType{ElemType: types.StringType},
						"call_recording_channels":             types.StringType,
						"call_recording_format":               types.StringType,
					},
					map[string]attr.Value{
						"call_recording_type":                 types.StringValue("all"),
						"call_recording_caller_phone_numbers": types.ListValueMust(types.StringType, []attr.Value{}),
						"call_recording_channels":             types.StringValue("single"),
						"call_recording_format":               types.StringValue("wav"),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"call_recording_type": schema.StringAttribute{
						Description: "Call recording type",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("all"),
					},
					"call_recording_caller_phone_numbers": schema.ListAttribute{
						Description: "Caller phone numbers for recording",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
					"call_recording_channels": schema.StringAttribute{
						Description: "Recording channels",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("single"),
					},
					"call_recording_format": schema.StringAttribute{
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
	fmt.Println(plan)
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

	callerPhoneNumbers, diagCP := convertListToStrings(ctx, callRecordingAttributes["call_recording_caller_phone_numbers"].(types.List))
	resp.Diagnostics.Append(diagCP...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.CreateOutboundVoiceProfile(telnyx.OutboundVoiceProfile{
		Name:                    plan.Name.ValueString(),
		BillingGroupID:          plan.BillingGroupID.ValueString(),
		TrafficType:             plan.TrafficType.ValueString(),
		ServicePlan:             plan.ServicePlan.ValueString(),
		ConcurrentCallLimit:     int(plan.ConcurrentCallLimit.ValueInt64()),
		Enabled:                 plan.Enabled.ValueBool(),
		Tags:                    tags,
		UsagePaymentMethod:      plan.UsagePaymentMethod.ValueString(),
		WhitelistedDestinations: whitelistedDestinations,
		MaxDestinationRate:      plan.MaxDestinationRate.ValueFloat64(),
		DailySpendLimit:         plan.DailySpendLimit.ValueString(),
		DailySpendLimitEnabled:  plan.DailySpendLimitEnabled.ValueBool(),
		CallRecording: telnyx.CallRecording{
			Type:               callRecordingAttributes["call_recording_type"].(types.String).ValueString(),
			CallerPhoneNumbers: callerPhoneNumbers,
			Channels:           callRecordingAttributes["call_recording_channels"].(types.String).ValueString(),
			Format:             callRecordingAttributes["call_recording_format"].(types.String).ValueString(),
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating outbound voice profile", err.Error())
		return
	}

	plan.ID = types.StringValue(profile.ID)
	plan.Name = types.StringValue(profile.Name)

	tflog.Info(ctx, "Created Outbound Voice Profile", map[string]interface{}{"id": profile.ID, "name": profile.Name})

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

	tflog.Info(ctx, "Reading Outbound Voice Profile", map[string]interface{}{"id": state.ID.ValueString()})

	profile, err := r.client.GetOutboundVoiceProfile(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound voice profile", err.Error())
		return
	}

	state.Name = types.StringValue(profile.Name)
	state.BillingGroupID = types.StringValue(profile.BillingGroupID)
	state.TrafficType = types.StringValue(profile.TrafficType)
	state.ServicePlan = types.StringValue(profile.ServicePlan)
	state.ConcurrentCallLimit = types.Int64Value(int64(profile.ConcurrentCallLimit))
	state.Enabled = types.BoolValue(profile.Enabled)
	state.Tags = convertStringsToList(profile.Tags)
	state.UsagePaymentMethod = types.StringValue(profile.UsagePaymentMethod)
	state.WhitelistedDestinations = convertStringsToList(profile.WhitelistedDestinations)
	state.MaxDestinationRate = types.Float64Value(profile.MaxDestinationRate)
	state.DailySpendLimit = types.StringValue(profile.DailySpendLimit)
	state.DailySpendLimitEnabled = types.BoolValue(profile.DailySpendLimitEnabled)
	state.CallRecording, diags = types.ObjectValue(map[string]attr.Type{
		"call_recording_type":                 types.StringType,
		"call_recording_caller_phone_numbers": types.ListType{ElemType: types.StringType},
		"call_recording_channels":             types.StringType,
		"call_recording_format":               types.StringType,
	}, map[string]attr.Value{
		"call_recording_type":                 types.StringValue(profile.CallRecording.Type),
		"call_recording_caller_phone_numbers": convertStringsToList(profile.CallRecording.CallerPhoneNumbers),
		"call_recording_channels":             types.StringValue(profile.CallRecording.Channels),
		"call_recording_format":               types.StringValue(profile.CallRecording.Format),
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

	callerPhoneNumbers, diagCP := convertListToStrings(ctx, callRecordingAttributes["call_recording_caller_phone_numbers"].(types.List))
	resp.Diagnostics.Append(diagCP...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateOutboundVoiceProfile(plan.ID.ValueString(), telnyx.OutboundVoiceProfile{
		Name:                    plan.Name.ValueString(),
		BillingGroupID:          plan.BillingGroupID.ValueString(),
		TrafficType:             plan.TrafficType.ValueString(),
		ServicePlan:             plan.ServicePlan.ValueString(),
		ConcurrentCallLimit:     int(plan.ConcurrentCallLimit.ValueInt64()),
		Enabled:                 plan.Enabled.ValueBool(),
		Tags:                    tags,
		UsagePaymentMethod:      plan.UsagePaymentMethod.ValueString(),
		WhitelistedDestinations: whitelistedDestinations,
		MaxDestinationRate:      plan.MaxDestinationRate.ValueFloat64(),
		DailySpendLimit:         plan.DailySpendLimit.ValueString(),
		DailySpendLimitEnabled:  plan.DailySpendLimitEnabled.ValueBool(),
		CallRecording: telnyx.CallRecording{
			Type:               callRecordingAttributes["call_recording_type"].(types.String).ValueString(),
			CallerPhoneNumbers: callerPhoneNumbers,
			Channels:           callRecordingAttributes["call_recording_channels"].(types.String).ValueString(),
			Format:             callRecordingAttributes["call_recording_format"].(types.String).ValueString(),
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating outbound voice profile", err.Error())
		return
	}

	tflog.Info(ctx, "Updated Outbound Voice Profile", map[string]interface{}{
		"id":                        plan.ID.ValueString(),
		"name":                      plan.Name.ValueString(),
		"billing_group_id":          plan.BillingGroupID.ValueString(),
		"traffic_type":              plan.TrafficType.ValueString(),
		"service_plan":              plan.ServicePlan.ValueString(),
		"concurrent_call_limit":     plan.ConcurrentCallLimit.ValueInt64(),
		"enabled":                   plan.Enabled.ValueBool(),
		"tags":                      tags,
		"usage_payment_method":      plan.UsagePaymentMethod.ValueString(),
		"whitelisted_destinations":  whitelistedDestinations,
		"max_destination_rate":      plan.MaxDestinationRate.ValueFloat64(),
		"daily_spend_limit":         plan.DailySpendLimit.ValueString(),
		"daily_spend_limit_enabled": plan.DailySpendLimitEnabled.ValueBool(),
		"call_recording_type":       callRecordingAttributes["call_recording_type"].(types.String).ValueString(),
		"call_recording_channels":   callRecordingAttributes["call_recording_channels"].(types.String).ValueString(),
		"call_recording_format":     callRecordingAttributes["call_recording_format"].(types.String).ValueString(),
		"caller_phone_numbers":      callerPhoneNumbers,
	})

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

	tflog.Info(ctx, "Deleted Outbound Voice Profile", map[string]interface{}{"id": state.ID.ValueString()})
}

func (r *OutboundVoiceProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

var (
	_ resource.Resource = &MessagingProfileResource{}
)

func NewMessagingProfileResource() resource.Resource {
	return &MessagingProfileResource{}
}

type MessagingProfileResource struct {
	client *telnyx.TelnyxClient
}

type MessagingProfileResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	WebhookURL              types.String `tfsdk:"webhook_url"`
	WebhookFailoverURL      types.String `tfsdk:"webhook_failover_url"`
	WebhookAPIVersion       types.String `tfsdk:"webhook_api_version"`
	WhitelistedDestinations types.List   `tfsdk:"whitelisted_destinations"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
	V1Secret                types.String `tfsdk:"v1_secret"`
}

func (r *MessagingProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_messaging_profile"
}

func (r *MessagingProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx Messaging Profiles",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the messaging profile",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the messaging profile",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Specifies whether the messaging profile is enabled or not",
				Required:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "The URL where webhooks related to this messaging profile will be sent",
				Optional:    true,
				Computed:    true,
			},
			"webhook_failover_url": schema.StringAttribute{
				Description: "The failover URL where webhooks related to this messaging profile will be sent if sending to the primary URL fails",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"webhook_api_version": schema.StringAttribute{
				Description: "Determines which webhook format will be used, Telnyx API v1, v2, or a legacy 2010-04-01 format",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("2"),
			},
			"whitelisted_destinations": schema.ListAttribute{
				Description: "Destinations to which the messaging profile is allowed to send",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("US")})),
			},
			"created_at": schema.StringAttribute{
				Description: "ISO 8601 formatted date indicating when the resource was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "ISO 8601 formatted date indicating when the resource was updated",
				Computed:    true,
			},
			"v1_secret": schema.StringAttribute{
				Description: "Secret used to authenticate with v1 endpoints",
				Computed:    true,
			},
		},
	}
}

func (r *MessagingProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for MessagingProfileResource")
	}
}

func (r *MessagingProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MessagingProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Messaging Profile", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	whitelistedDestinations, diagWD := convertListToStrings(ctx, plan.WhitelistedDestinations)
	resp.Diagnostics.Append(diagWD...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Converted whitelisted destinations", map[string]interface{}{
		"whitelisted_destinations": whitelistedDestinations,
	})

	profile, err := r.client.CreateMessagingProfile(telnyx.MessagingProfile{
		Name:                    plan.Name.ValueString(),
		Enabled:                 plan.Enabled.ValueBool(),
		WebhookURL:              plan.WebhookURL.ValueString(),
		WebhookFailoverURL:      plan.WebhookFailoverURL.ValueString(),
		WebhookAPIVersion:       plan.WebhookAPIVersion.ValueString(),
		WhitelistedDestinations: whitelistedDestinations,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating messaging profile", err.Error())
		return
	}

	tflog.Info(ctx, "Created Messaging Profile", map[string]interface{}{
		"id": profile.ID,
	})

	plan.ID = types.StringValue(profile.ID)
	plan.CreatedAt = types.StringValue(profile.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(profile.UpdatedAt.String())
	plan.V1Secret = types.StringValue(profile.V1Secret)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *MessagingProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MessagingProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.GetMessagingProfile(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading messaging profile", err.Error())
		return
	}

	state.Name = types.StringValue(profile.Name)
	state.Enabled = types.BoolValue(profile.Enabled)
	state.WebhookURL = types.StringValue(profile.WebhookURL)
	state.WebhookFailoverURL = types.StringValue(profile.WebhookFailoverURL)
	state.WebhookAPIVersion = types.StringValue(profile.WebhookAPIVersion)
	state.WhitelistedDestinations = convertStringsToList(profile.WhitelistedDestinations)
	state.CreatedAt = types.StringValue(profile.CreatedAt.String())
	state.UpdatedAt = types.StringValue(profile.UpdatedAt.String())
	state.V1Secret = types.StringValue(profile.V1Secret)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *MessagingProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MessagingProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	whitelistedDestinations, diagWD := convertListToStrings(ctx, plan.WhitelistedDestinations)
	resp.Diagnostics.Append(diagWD...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.UpdateMessagingProfile(plan.ID.ValueString(), telnyx.MessagingProfile{
		Name:                    plan.Name.ValueString(),
		Enabled:                 plan.Enabled.ValueBool(),
		WebhookURL:              plan.WebhookURL.ValueString(),
		WebhookFailoverURL:      plan.WebhookFailoverURL.ValueString(),
		WebhookAPIVersion:       plan.WebhookAPIVersion.ValueString(),
		WhitelistedDestinations: whitelistedDestinations,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating messaging profile", err.Error())
		return
	}

	plan.CreatedAt = types.StringValue(profile.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(profile.UpdatedAt.String())
	plan.V1Secret = types.StringValue(profile.V1Secret)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *MessagingProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MessagingProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMessagingProfile(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting messaging profile", err.Error())
	}
}

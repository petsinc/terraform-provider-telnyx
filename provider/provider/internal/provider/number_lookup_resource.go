package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

type PhoneNumberLookupResource struct {
	client *telnyx.TelnyxClient
}

func NewPhoneNumberLookupResource() resource.Resource {
	return &PhoneNumberLookupResource{}
}

func (r *PhoneNumberLookupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_phone_number_lookup"
}

func (r *PhoneNumberLookupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"starts_with": schema.StringAttribute{
				Optional: true,
			},
			"ends_with": schema.StringAttribute{
				Optional: true,
			},
			"contains": schema.StringAttribute{
				Optional: true,
			},
			"locality": schema.StringAttribute{
				Optional: true,
			},
			"administrative_area": schema.StringAttribute{
				Optional: true,
			},
			"country_code": schema.StringAttribute{
				Optional: true,
			},
			"national_destination_code": schema.StringAttribute{
				Optional: true,
			},
			"rate_center": schema.StringAttribute{
				Optional: true,
			},
			"phone_number_type": schema.StringAttribute{
				Optional: true,
			},
			"features": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"limit": schema.Int64Attribute{
				Optional: true,
			},
			"best_effort": schema.BoolAttribute{
				Optional: true,
			},
			"quickship": schema.BoolAttribute{
				Optional: true,
			},
			"reservable": schema.BoolAttribute{
				Optional: true,
			},
			"exclude_held_numbers": schema.BoolAttribute{
				Optional: true,
			},
			"phone_numbers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"phone_number": schema.StringAttribute{
							Computed: true,
						},
						"reservable": schema.BoolAttribute{
							Computed: true,
						},
						"upfront_cost": schema.StringAttribute{
							Computed: true,
						},
						"monthly_cost": schema.StringAttribute{
							Computed: true,
						},
						"features": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *PhoneNumberLookupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for PhoneNumberLookupResource")
	}
}

type PhoneNumberLookupModel struct {
	StartsWith              types.String `tfsdk:"starts_with"`
	EndsWith                types.String `tfsdk:"ends_with"`
	Contains                types.String `tfsdk:"contains"`
	Locality                types.String `tfsdk:"locality"`
	AdministrativeArea      types.String `tfsdk:"administrative_area"`
	CountryCode             types.String `tfsdk:"country_code"`
	NationalDestinationCode types.String `tfsdk:"national_destination_code"`
	RateCenter              types.String `tfsdk:"rate_center"`
	PhoneNumberType         types.String `tfsdk:"phone_number_type"`
	Features                types.List   `tfsdk:"features"`
	Limit                   types.Int64  `tfsdk:"limit"`
	BestEffort              types.Bool   `tfsdk:"best_effort"`
	Quickship               types.Bool   `tfsdk:"quickship"`
	Reservable              types.Bool   `tfsdk:"reservable"`
	ExcludeHeldNumbers      types.Bool   `tfsdk:"exclude_held_numbers"`
	PhoneNumbers            types.List   `tfsdk:"phone_numbers"`
}

func (r *PhoneNumberLookupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PhoneNumberLookupModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := telnyx.AvailablePhoneNumbersRequest{
		StartsWith:              plan.StartsWith.ValueString(),
		EndsWith:                plan.EndsWith.ValueString(),
		Contains:                plan.Contains.ValueString(),
		Locality:                plan.Locality.ValueString(),
		AdministrativeArea:      plan.AdministrativeArea.ValueString(),
		CountryCode:             plan.CountryCode.ValueString(),
		NationalDestinationCode: plan.NationalDestinationCode.ValueString(),
		RateCenter:              plan.RateCenter.ValueString(),
		PhoneNumberType:         plan.PhoneNumberType.ValueString(),
		Limit:                   int(plan.Limit.ValueInt64()),
		BestEffort:              plan.BestEffort.ValueBool(),
		Quickship:               plan.Quickship.ValueBool(),
		Reservable:              plan.Reservable.ValueBool(),
		ExcludeHeldNumbers:      plan.ExcludeHeldNumbers.ValueBool(),
	}

	if !plan.Features.IsNull() && !plan.Features.IsUnknown() {
		var features []string
		diags = plan.Features.ElementsAs(ctx, &features, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		filters.Features = features
	}

	client := r.client
	response, err := client.ListAvailablePhoneNumbers(filters)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving available phone numbers",
			"Could not retrieve available phone numbers, unexpected error: "+err.Error(),
		)
		return
	}

	phoneNumbers := make([]attr.Value, len(response.Data))
	for i, phoneNumber := range response.Data {
		featureNames := make([]string, len(phoneNumber.Features))
		for j, feature := range phoneNumber.Features {
			featureNames[j] = feature.Name
		}

		sort.Strings(featureNames)

		features := make([]attr.Value, len(featureNames))
		for j, featureName := range featureNames {
			features[j] = types.StringValue(featureName)
		}

		phoneNumbers[i], diags = types.ObjectValue(
			map[string]attr.Type{
				"phone_number": types.StringType,
				"reservable":   types.BoolType,
				"upfront_cost": types.StringType,
				"monthly_cost": types.StringType,
				"features":     types.ListType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"phone_number": types.StringValue(phoneNumber.PhoneNumber),
				"reservable":   types.BoolValue(phoneNumber.Reservable),
				"upfront_cost": types.StringValue(phoneNumber.CostInformation.UpfrontCost),
				"monthly_cost": types.StringValue(phoneNumber.CostInformation.MonthlyCost),
				"features":     types.ListValueMust(types.StringType, features),
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	phoneNumbersList := types.ListValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"phone_number": types.StringType,
			"reservable":   types.BoolType,
			"upfront_cost": types.StringType,
			"monthly_cost": types.StringType,
			"features":     types.ListType{ElemType: types.StringType},
		},
	}, phoneNumbers)

	state := PhoneNumberLookupModel{
		StartsWith:              plan.StartsWith,
		EndsWith:                plan.EndsWith,
		Contains:                plan.Contains,
		Locality:                plan.Locality,
		AdministrativeArea:      plan.AdministrativeArea,
		CountryCode:             plan.CountryCode,
		NationalDestinationCode: plan.NationalDestinationCode,
		RateCenter:              plan.RateCenter,
		PhoneNumberType:         plan.PhoneNumberType,
		Features:                plan.Features,
		Limit:                   plan.Limit,
		 BestEffort:             plan.BestEffort,
		Quickship:               plan.Quickship,
		Reservable:              plan.Reservable,
		ExcludeHeldNumbers:      plan.ExcludeHeldNumbers,
		PhoneNumbers:            phoneNumbersList,
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PhoneNumberLookupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PhoneNumberLookupModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No API call is needed, just return the current state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PhoneNumberLookupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Use the same logic as Create, cuz resource is static
	r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, (*resource.CreateResponse)(resp))
}


func (r *PhoneNumberLookupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Remove the resource from the state
	resp.State.RemoveResource(ctx)
}

// Implement ImportState function
func (r *PhoneNumberLookupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

type PhoneNumberLookupDataSource struct {
	client *telnyx.TelnyxClient
}

func NewPhoneNumberLookupDataSource() datasource.DataSource {
	return &PhoneNumberLookupDataSource{}
}

func (d *PhoneNumberLookupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		client, ok := req.ProviderData.(*telnyx.TelnyxClient)
		if !ok {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected *telnyx.TelnyxClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
			return
		}
		d.client = client
	}
}

func (d *PhoneNumberLookupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_phone_number_lookup"
}

func (d *PhoneNumberLookupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
	Limit                   types.Int64  `tfsdk:"limit"`
	BestEffort              types.Bool   `tfsdk:"best_effort"`
	Quickship               types.Bool   `tfsdk:"quickship"`
	Reservable              types.Bool   `tfsdk:"reservable"`
	ExcludeHeldNumbers      types.Bool   `tfsdk:"exclude_held_numbers"`
	PhoneNumbers            types.List   `tfsdk:"phone_numbers"`
}

func (d *PhoneNumberLookupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PhoneNumberLookupModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := telnyx.AvailablePhoneNumbersRequest{
		StartsWith:              config.StartsWith.ValueString(),
		EndsWith:                config.EndsWith.ValueString(),
		Contains:                config.Contains.ValueString(),
		Locality:                config.Locality.ValueString(),
		AdministrativeArea:      config.AdministrativeArea.ValueString(),
		CountryCode:             config.CountryCode.ValueString(),
		NationalDestinationCode: config.NationalDestinationCode.ValueString(),
		RateCenter:              config.RateCenter.ValueString(),
		PhoneNumberType:         config.PhoneNumberType.ValueString(),
		Limit:                   int(config.Limit.ValueInt64()),
		BestEffort:              config.BestEffort.ValueBool(),
		Quickship:               config.Quickship.ValueBool(),
		Reservable:              config.Reservable.ValueBool(),
		ExcludeHeldNumbers:      config.ExcludeHeldNumbers.ValueBool(),
	}

	client := d.client
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
		// Convert features to a slice of strings
		featureNames := make([]string, len(phoneNumber.Features))
		for j, feature := range phoneNumber.Features {
			featureNames[j] = feature.Name
		}

		// Sort the feature names
		sort.Strings(featureNames)

		// Convert sorted feature names to attr.Value
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

	diags = resp.State.Set(ctx, &PhoneNumberLookupModel{
		PhoneNumbers: phoneNumbersList,
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

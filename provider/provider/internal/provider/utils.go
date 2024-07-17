package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func convertListToStrings(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	var strings []string
	diags := list.ElementsAs(ctx, &strings, false)
	return strings, diags
}

func convertStringsToList(strings []string) types.List {
	elements := make([]attr.Value, len(strings))
	for i, str := range strings {
		elements[i] = types.StringValue(str)
	}
	return types.ListValueMust(types.StringType, elements)
}

// ConvertListToPhoneNumbers converts a types.List to a slice of PhoneNumberResourceModel.
func ConvertListToPhoneNumbers(ctx context.Context, list types.List) ([]PhoneNumberResourceModel, diag.Diagnostics) {
	var phoneNumbers []PhoneNumberResourceModel
	diags := list.ElementsAs(ctx, &phoneNumbers, false)
	return phoneNumbers, diags
}

// ConvertPhoneNumbersToList converts a slice of PhoneNumberResourceModel to a types.List.
func ConvertPhoneNumbersToList(ctx context.Context, phoneNumbers []PhoneNumberResourceModel) (types.List, diag.Diagnostics) {
	elements := make([]attr.Value, len(phoneNumbers))
	for i, pn := range phoneNumbers {
		elements[i] = types.ObjectValueMust(
			PhoneNumberResourceModel{}.AttrTypes(),
			map[string]attr.Value{
				"id":                      pn.ID,
				"phone_number":            pn.PhoneNumber,
				"status":                  pn.Status,
				"regulatory_requirements": pn.RegulatoryRequirements,
			},
		)
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: PhoneNumberResourceModel{}.AttrTypes()}, elements), nil
}

// ConvertListToRegulatoryRequirements converts a types.List to a slice of RegulatoryRequirementResourceModel.
func ConvertListToRegulatoryRequirements(ctx context.Context, list types.List) ([]RegulatoryRequirementResourceModel, diag.Diagnostics) {
	var requirements []RegulatoryRequirementResourceModel
	diags := list.ElementsAs(ctx, &requirements, false)
	return requirements, diags
}

// ConvertRegulatoryRequirementsToList converts a slice of RegulatoryRequirementResourceModel to a types.List.
func ConvertRegulatoryRequirementsToList(ctx context.Context, requirements []RegulatoryRequirementResourceModel) (types.List, diag.Diagnostics) {
	elements := make([]attr.Value, len(requirements))
	for i, req := range requirements {
		elements[i] = types.ObjectValueMust(
			RegulatoryRequirementResourceModel{}.AttrTypes(),
			map[string]attr.Value{
				"requirement_id": req.RequirementID,
				"field_value":    req.FieldValue,
				"field_type":     req.FieldType,
			},
		)
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: RegulatoryRequirementResourceModel{}.AttrTypes()}, elements), nil
}

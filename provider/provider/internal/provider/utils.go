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

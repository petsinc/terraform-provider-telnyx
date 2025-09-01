package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
)

var (
	_ resource.Resource = &NumberOrderResource{}
)

func NewNumberOrderResource() resource.Resource {
	return &NumberOrderResource{}
}

type NumberOrderResource struct {
	client *telnyx.TelnyxClient
}

type NumberOrderResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ConnectionID       types.String `tfsdk:"connection_id"`
	MessagingProfileID types.String `tfsdk:"messaging_profile_id"`
	BillingGroupID     types.String `tfsdk:"billing_group_id"`
	CustomerReference  types.String `tfsdk:"customer_reference"`
	Status             types.String `tfsdk:"status"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
	PhoneNumbers       types.List   `tfsdk:"phone_numbers"`
	SubNumberOrderIDs  types.List   `tfsdk:"sub_number_orders_ids"`
}

type PhoneNumberResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	PhoneNumber            types.String `tfsdk:"phone_number"`
	Status                 types.String `tfsdk:"status"`
	RegulatoryRequirements types.List   `tfsdk:"regulatory_requirements"`
}

type RegulatoryRequirementResourceModel struct {
	RequirementID types.String `tfsdk:"requirement_id"`
	FieldValue    types.String `tfsdk:"field_value"`
	FieldType     types.String `tfsdk:"field_type"`
}

func (p PhoneNumberResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"phone_number":            types.StringType,
		"status":                  types.StringType,
		"regulatory_requirements": types.ListType{ElemType: types.ObjectType{AttrTypes: RegulatoryRequirementResourceModel{}.AttrTypes()}},
	}
}

func (r RegulatoryRequirementResourceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"requirement_id": types.StringType,
		"field_value":    types.StringType,
		"field_type":     types.StringType,
	}
}

func (r *NumberOrderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_number_order"
}

func (r *NumberOrderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx Number Orders",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the number order",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_id": schema.StringAttribute{
				Description: "Connection ID associated with the number order",
				Required:    true,
			},
			"messaging_profile_id": schema.StringAttribute{
				Description: "Messaging profile ID associated with the number order",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"billing_group_id": schema.StringAttribute{
				Description: "Billing group ID associated with the number order",
				Required:    true,
			},
			"customer_reference": schema.StringAttribute{
				Description: "Customer reference for the number order",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the number order",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation time of the number order",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update time of the number order",
				Computed:    true,
			},
			"phone_numbers": schema.ListNestedAttribute{
				Description: "List of phone numbers in the order",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the phone number",
							Computed:    true,
						},
						"phone_number": schema.StringAttribute{
							Description: "Phone number in E.164 format",
							Required:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the phone number",
							Computed:    true,
						},
						"regulatory_requirements": schema.ListNestedAttribute{
							Description: "List of regulatory requirements for the phone number",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"requirement_id": schema.StringAttribute{
										Description: "Unique identifier of the regulatory requirement",
										Computed:    true,
									},
									"field_value": schema.StringAttribute{
										Description: "Value of the regulatory requirement field",
										Computed:    true,
									},
									"field_type": schema.StringAttribute{
										Description: "Type of the regulatory requirement field",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"sub_number_orders_ids": schema.ListAttribute{
				Description: "List of sub number order IDs associated with the number order",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *NumberOrderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for NumberOrderResource")
	}
}

func (r *NumberOrderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NumberOrderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	phoneNumbers, diags := ConvertListToPhoneNumbers(ctx, plan.PhoneNumbers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var phoneNumbersRequest []telnyx.PhoneNumberRequest
	for _, phoneNumber := range phoneNumbers {
		phoneNumbersRequest = append(phoneNumbersRequest, telnyx.PhoneNumberRequest{PhoneNumber: phoneNumber.PhoneNumber.ValueString()})
	}

	request := telnyx.CreateNumberOrderRequest{
		PhoneNumbers:       phoneNumbersRequest,
		ConnectionID:       plan.ConnectionID.ValueString(),
		MessagingProfileID: plan.MessagingProfileID.ValueString(),
		BillingGroupID:     plan.BillingGroupID.ValueString(),
		CustomerReference:  plan.CustomerReference.ValueString(),
	}

	order, err := r.client.CreateNumberOrder(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating number order", err.Error())
		return
	}

	// Set state based on response from the API
	setStateFromOrderResponse(&plan, order)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *NumberOrderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NumberOrderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	order, err := r.client.GetNumberOrder(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading number order", err.Error())
		return
	}
	// Update state based on response from the API
	setStateFromOrderResponse(&state, order)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *NumberOrderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NumberOrderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO Prepare the request with an empty array for regulatory requirements until we support this prop
	request := telnyx.UpdateNumberOrderRequest{
		CustomerReference:      plan.CustomerReference.ValueString(),
		RegulatoryRequirements: []telnyx.NumberOrderRegulatoryRequirement{}, // Ensure this is always an empty array
	}

	order, err := r.client.UpdateNumberOrder(plan.ID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error updating number order", err.Error())
		return
	}

	// Update state based on response from the API
	setStateFromOrderResponse(&plan, order)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *NumberOrderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NumberOrderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the SubNumberOrderIDs list to a slice of strings using the utility function
	subNumberOrderIDs, diags := convertListToStrings(ctx, state.SubNumberOrderIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Cancel sub number orders if they exist
	for _, subOrderID := range subNumberOrderIDs {
		if subOrder, err := r.client.GetSubNumberOrder(subOrderID); err == nil {
			if subOrder.Status == "deleted" {
				continue
			}
		}
		fmt.Printf("Cancelling sub number order with ID: %s\n", subOrderID)
		_, err := r.client.CancelSubNumberOrder(subOrderID)
		if err != nil {
			fmt.Printf("Error cancelling sub number order ID: %s, Error: %s\n", subOrderID, err.Error())
			resp.Diagnostics.AddError("Error cancelling sub number order", fmt.Sprintf("ID: %s, Error: %s", subOrderID, err.Error()))
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

func setStateFromOrderResponse(state *NumberOrderResourceModel, order *telnyx.PhoneNumberOrderResponse) {
	state.ID = types.StringValue(order.ID)
	state.Status = types.StringValue(order.Status)
	state.CreatedAt = types.StringValue(order.CreatedAt.String())
	state.UpdatedAt = types.StringValue(order.UpdatedAt.String())

	var phoneNumbersModel []PhoneNumberResourceModel
	for _, pn := range order.PhoneNumbers {
		phoneNumberModel := PhoneNumberResourceModel{
			ID:          types.StringValue(pn.ID),
			PhoneNumber: types.StringValue(pn.PhoneNumber),
			Status:      types.StringValue(pn.Status),
		}

		if pn.RegulatoryRequirements != nil {
			reqs := make([]RegulatoryRequirementResourceModel, len(pn.RegulatoryRequirements))
			for j, req := range pn.RegulatoryRequirements {
				reqs[j] = RegulatoryRequirementResourceModel{
					RequirementID: types.StringValue(req.RequirementID),
					FieldValue:    types.StringValue(req.FieldValue),
					FieldType:     types.StringValue(req.FieldType),
				}
			}
			phoneNumberModel.RegulatoryRequirements, _ = ConvertRegulatoryRequirementsToList(context.Background(), reqs)
		} else {
			phoneNumberModel.RegulatoryRequirements = types.ListValueMust(
				types.ObjectType{AttrTypes: RegulatoryRequirementResourceModel{}.AttrTypes()}, []attr.Value{})
		}

		phoneNumbersModel = append(phoneNumbersModel, phoneNumberModel)
	}

	state.PhoneNumbers, _ = ConvertPhoneNumbersToList(context.Background(), phoneNumbersModel)
	state.SubNumberOrderIDs = convertStringsToList(order.SubNumberOrderIDs)
}

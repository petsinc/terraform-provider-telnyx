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
	_ resource.Resource = &FQDNResource{}
)

func NewFQDNResource() resource.Resource {
	return &FQDNResource{}
}

type FQDNResource struct {
	client *telnyx.TelnyxClient
}

type FQDNResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ConnectionID  types.Int64  `tfsdk:"connection_id"`
	FQDN          types.String `tfsdk:"fqdn"`
	Port          types.Int64  `tfsdk:"port"`
	DNSRecordType types.String `tfsdk:"dns_record_type"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (r *FQDNResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fqdn"
}

func (r *FQDNResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Telnyx FQDNs",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the FQDN",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_id": schema.Int64Attribute{
				Description: "ID of the connection associated with the FQDN",
				Required:    true,
			},
			"fqdn": schema.StringAttribute{
				Description: "Fully Qualified Domain Name",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Port associated with the FQDN",
				Required:    true,
			},
			"dns_record_type": schema.StringAttribute{
				Description: "DNS record type",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "ISO 8601 formatted date indicating when the resource was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "ISO 8601 formatted date indicating when the resource was updated",
				Computed:    true,
			},
		},
	}
}

func (r *FQDNResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
		tflog.Info(ctx, "Configured Telnyx client for FQDNResource")
	}
}

func (r *FQDNResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FQDNResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fqdn := telnyx.FQDN{
		ConnectionID:  int(plan.ConnectionID.ValueInt64()),
		FQDN:          plan.FQDN.ValueString(),
		Port:          int(plan.Port.ValueInt64()),
		DNSRecordType: plan.DNSRecordType.ValueString(),
	}

	createdFQDN, err := r.client.CreateFQDN(fqdn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FQDN",
			"Could not create FQDN, unexpected error: "+err.Error(),
		)
		return
	}

	setFQDNState(ctx, &plan, createdFQDN)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FQDNResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FQDNResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fqdn, err := r.client.GetFQDN(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FQDN",
			"Could not read FQDN, unexpected error: "+err.Error(),
		)
		return
	}

	setFQDNState(ctx, &state, fqdn)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *FQDNResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FQDNResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FQDNResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fqdn := telnyx.FQDN{
		ConnectionID:  int(plan.ConnectionID.ValueInt64()),
		FQDN:          plan.FQDN.ValueString(),
		Port:          int(plan.Port.ValueInt64()),
		DNSRecordType: plan.DNSRecordType.ValueString(),
	}

	updatedFQDN, err := r.client.UpdateFQDN(state.ID.ValueString(), fqdn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating FQDN",
			"Could not update FQDN, unexpected error: "+err.Error(),
		)
		return
	}

	setFQDNState(ctx, &plan, updatedFQDN)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *FQDNResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FQDNResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFQDN(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FQDN",
			"Could not delete FQDN, unexpected error: "+err.Error(),
		)
		return
	}
}

func setFQDNState(ctx context.Context, state *FQDNResourceModel, fqdn *telnyx.FQDN) {
	state.ID = types.StringValue(fqdn.ID)
	state.ConnectionID = types.Int64Value(int64(fqdn.ConnectionID))
	state.FQDN = types.StringValue(fqdn.FQDN)
	state.Port = types.Int64Value(int64(fqdn.Port))
	state.DNSRecordType = types.StringValue(fqdn.DNSRecordType)
	state.CreatedAt = types.StringValue(fqdn.CreatedAt.String())
	state.UpdatedAt = types.StringValue(fqdn.UpdatedAt.String())
}

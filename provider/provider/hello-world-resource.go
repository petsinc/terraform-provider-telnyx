package provider

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var (
    _ resource.Resource = &HelloWorldResource{}
)

func NewHelloWorldResource() resource.Resource {
    return &HelloWorldResource{}
}

type HelloWorldResource struct {
    client *http.Client
}

type HelloWorldResourceModel struct {
    ID      types.String `tfsdk:"id"`
    Message types.String `tfsdk:"message"`
}

func (r *HelloWorldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_request"
}

func (r *HelloWorldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Simple Hello World Resource",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Resource identifier",
                Computed:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "message": schema.StringAttribute{
                Description: "Message to send",
                Required:    true,
            },
        },
    }
}

func (r *HelloWorldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData != nil {
        client, ok := req.ProviderData.(*http.Client)
        if !ok {
            resp.Diagnostics.AddError(
                "Unexpected Resource Configure Type",
                fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
            )
            return
        }
        r.client = client
    }
}

func (r *HelloWorldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan HelloWorldResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    reqBody := fmt.Sprintf(`{"message": "%s"}`, plan.Message.ValueString())
    httpReq, err := http.NewRequest("POST", "http://httpbin.org/post", strings.NewReader(reqBody))
    if err != nil {
        resp.Diagnostics.AddError("Error creating request", err.Error())
        return
    }
    httpReq.Header.Set("Content-Type", "application/json")

    httpResp, err := r.client.Do(httpReq)
    if err != nil {
        resp.Diagnostics.AddError("Error sending request", err.Error())
        return
    }
    defer httpResp.Body.Close()

    plan.ID = types.StringValue("unique-id")

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
}

func (r *HelloWorldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // No-op for this example
}

func (r *HelloWorldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // No-op for this example
}

func (r *HelloWorldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // No-op for this example
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx" // Ensure this path is correct
)

var (
	_ provider.Provider = &TelnyxProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TelnyxProvider{
			version: version,
		}
	}
}

type TelnyxProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

type TelnyxProvider struct {
	version string
}

func (p *TelnyxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "telnyx"
	resp.Version = p.version
}

func (p *TelnyxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Telnyx Provider",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "Endpoint for HTTP requests",
				Required:    true,
			},
		},
	}
}

func (p *TelnyxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config TelnyxProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := telnyx.NewClient()

	tflog.Info(ctx, "Configured Telnyx provider", map[string]interface{}{"endpoint": config.Endpoint})

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TelnyxProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBillingGroupResource,
		NewOutboundVoiceProfileResource,
	}
}

func (p *TelnyxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

package main

import (
	"context"
	"flag"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/petsinc/terraform-provider-telnyx/internal/provider"
)

var version string

func init() {
	// Version can be injected during build time
	// Example: go build -ldflags "-X main.version=$(VERSION)"
	if version == "" {
		version = "dev"
	}
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/petsinc/telnyx",
		Debug:   debug,
	}

	tflog.Info(context.Background(), "Starting telnyx provider", map[string]interface{}{
		"version": version,
	})

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		tflog.Error(context.Background(), "Error starting provider", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}

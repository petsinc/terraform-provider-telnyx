package main

import (
	"github.com/petsinc/telnyx-rest-client/internal/test_runner"
	"github.com/petsinc/telnyx-rest-client/pkg/telnyx"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	client := telnyx.NewClient()
	runner := test_runner.NewTestRunner(client, logger)

	// Perform create operations
	// runner.PerformCreates()

	// // Perform update operations
	// runner.PerformUpdates()

	// // Perform cascading delete operations
	// runner.PerformCascadingDeletes()
}

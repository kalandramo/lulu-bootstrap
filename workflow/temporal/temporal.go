// Package temporal provides a bootstrap workflow builder for Temporal.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/workflow/temporal"
package temporal

import (
	"context"
	"fmt"

	temporalPlugin "github.com/kalandramo/lulu-ext/workflow/temporal"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterWorkflowBuilder(bootstrap.WorkflowTypeTemporal, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Workflow) (any, func(), error) {
	c := cfg.GetTemporal()
	if c == nil {
		return nil, nil, fmt.Errorf("temporal: config is nil")
	}

	var opts []func(*temporalPlugin.ClientOptions)

	if hp := c.GetHostPort(); hp != "" {
		opts = append(opts, temporalPlugin.WithClientHostPort(hp))
	}
	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, temporalPlugin.WithClientNamespace(ns))
	}

	client, err := temporalPlugin.NewClient(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("temporal: create client: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}

// Package nats provides a bootstrap broker builder for NATS.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/broker/nats"
package nats

import (
	"context"
	"fmt"

	natsPlugin "github.com/kalandramo/lulu-ext/transport/nats"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeNATS, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetNats()
	if c == nil {
		return nil, nil, fmt.Errorf("nats: config is nil")
	}

	var opts []natsPlugin.ServerOption

	if url := c.GetUrl(); url != "" {
		opts = append(opts, natsPlugin.WithAddress([]string{url}))
	}

	srv := natsPlugin.NewServer(opts...)
	return srv, func() {}, nil
}

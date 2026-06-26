// Package redis provides a bootstrap broker builder for Redis (Pub/Sub).
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/broker/redis"
package redis

import (
	"context"
	"fmt"

	redisPlugin "github.com/kalandramo/lulu-ext/transport/redis"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeRedis, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetRedis()
	if c == nil {
		return nil, nil, fmt.Errorf("redis: config is nil")
	}

	var opts []redisPlugin.ServerOption

	if addr := c.GetAddress(); addr != "" {
		opts = append(opts, redisPlugin.WithAddress(addr))
	}

	srv := redisPlugin.NewServer(opts...)
	return srv, func() {}, nil
}

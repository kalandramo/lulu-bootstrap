// Package asynq provides a bootstrap server builder for Asynq (Redis-based task queue).
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/transport/asynq"
package asynq

import (
	"fmt"

	asynqPlugin "github.com/kalandramo/lulu-ext/transport/asynq"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/kalandramo/lulu/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeAsynq, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetAsynq()
	if c == nil {
		return nil, fmt.Errorf("asynq: config is nil")
	}

	var opts []asynqPlugin.Option

	if addr := c.GetRedisAddress(); addr != "" {
		opts = append(opts, asynqPlugin.WithRedisAddress(addr))
	}
	if pwd := c.GetRedisPassword(); pwd != "" {
		opts = append(opts, asynqPlugin.WithRedisPassword(pwd))
	}
	if db := c.GetRedisDb(); db > 0 {
		opts = append(opts, asynqPlugin.WithRedisDB(db))
	}

	srv := asynqPlugin.NewServer(opts...)
	return srv, nil
}

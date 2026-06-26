// Package nats provides a bootstrap server builder for NATS transport.
package nats

import (
	"fmt"

	natsPlugin "github.com/kalandramo/lulu-ext/transport/nats"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/kalandramo/lulu/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeNATS, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetNats()
	if c == nil {
		return nil, fmt.Errorf("nats: config is nil")
	}

	var opts []natsPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, natsPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, natsPlugin.WithCodec(codec))
	}

	return natsPlugin.NewServer(opts...), nil
}

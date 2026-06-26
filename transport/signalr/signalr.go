// Package signalr provides a bootstrap server builder for SignalR.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/transport/signalr"
package signalr

import (
	"crypto/tls"
	"fmt"

	signalrPlugin "github.com/kalandramo/lulu-ext/transport/signalr"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/kalandramo/lulu/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeSignalR, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetSignalr()
	if c == nil {
		return nil, fmt.Errorf("signalr: config is nil")
	}

	var opts []signalrPlugin.Option

	if addr := c.GetAddr(); addr != "" {
		opts = append(opts, signalrPlugin.WithAddress(addr))
	}
	if network := c.GetNetwork(); network != "" {
		opts = append(opts, signalrPlugin.WithNetwork(network))
	}
	if cert := c.GetTls(); cert != nil {
		opts = append(opts, signalrPlugin.WithTLSConfig(&tls.Config{}))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, signalrPlugin.WithCodec(codec))
	}

	srv := signalrPlugin.NewServer(opts...)
	return srv, nil
}

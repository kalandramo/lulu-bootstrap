// Package webtransport provides a bootstrap server builder for the WebTransport (HTTP/3) transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/transport/webtransport"
package webtransport

import (
	"crypto/tls"
	"fmt"

	wtPlugin "github.com/kalandramo/lulu-ext/transport/webtransport"
	"github.com/kalandramo/lulu/transport"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeWebTransport, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetWebtransport()
	if c == nil {
		return nil, fmt.Errorf("webtransport: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":443"
	}

	var opts []wtPlugin.ServerOption
	opts = append(opts, wtPlugin.WithAddress(addr))

	if tlsCfg := buildTLS(c.GetTls()); tlsCfg != nil {
		opts = append(opts, wtPlugin.WithTLSConfig(tlsCfg))
	}

	srv := wtPlugin.NewServer(opts...)
	return srv, nil
}

func buildTLS(tlsCfg *v1.Server_TLS) *tls.Config {
	if tlsCfg == nil {
		return nil
	}
	if f := tlsCfg.GetFile(); f != nil {
		cert, err := tls.LoadX509KeyPair(f.GetCertPath(), f.GetKeyPath())
		if err != nil {
			return nil
		}
		return &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: tlsCfg.GetInsecureSkipVerify(),
		}
	}
	if c := tlsCfg.GetConfig(); c != nil {
		cert, err := tls.X509KeyPair(c.GetCertPem(), c.GetKeyPem())
		if err != nil {
			return nil
		}
		return &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: tlsCfg.GetInsecureSkipVerify(),
		}
	}
	return nil
}

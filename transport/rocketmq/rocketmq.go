// Package rocketmq provides a bootstrap server builder for RocketMQ transport.
package rocketmq

import (
	"fmt"

	rocketmqOption "github.com/kalandramo/lulu-ext/broker/rocketmq/option"
	rocketmqPlugin "github.com/kalandramo/lulu-ext/transport/rocketmq"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/kalandramo/lulu/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeRocketMQ, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetRocketmq()
	if c == nil {
		return nil, fmt.Errorf("rocketmq: config is nil")
	}

	// Determine driver type, default to v5.
	driverType := rocketmqOption.DriverTypeV5
	if dt := c.GetDriverType(); dt != "" {
		driverType = rocketmqOption.DriverType(dt)
	}

	var opts []rocketmqPlugin.ServerOption
	if addrs := c.GetNameServerAddrs(); len(addrs) > 0 {
		opts = append(opts, rocketmqPlugin.WithNameServer(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, rocketmqPlugin.WithCodec(codec))
	}

	return rocketmqPlugin.NewServer(driverType, opts...), nil
}

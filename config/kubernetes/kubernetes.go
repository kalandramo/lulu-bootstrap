// Package kubernetes provides a bootstrap config action for Kubernetes ConfigMap config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/config/kubernetes"
package kubernetes

import (
	"context"
	"fmt"

	kubePlugin "github.com/kalandramo/lulu-ext/config/kubernetes"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeKubernetes, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetKubernetes()
	if c == nil {
		return nil, fmt.Errorf("kubernetes: config is nil")
	}

	var opts []kubePlugin.Option

	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, kubePlugin.WithNamespace(ns))
	}
	if label := c.GetLabelSelector(); label != "" {
		opts = append(opts, kubePlugin.WithLabelSelector(label))
	}
	if field := c.GetFieldSelector(); field != "" {
		opts = append(opts, kubePlugin.WithFieldSelector(field))
	}
	if kubeConfig := c.GetKubeConfig(); kubeConfig != "" {
		opts = append(opts, kubePlugin.WithKubeConfig(kubeConfig))
	}
	if master := c.GetMaster(); master != "" {
		opts = append(opts, kubePlugin.WithMaster(master))
	}

	kubePlugin.New(opts...)

	return func() {}, nil
}

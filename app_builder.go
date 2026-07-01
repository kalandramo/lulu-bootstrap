package bootstrap

import (
	"time"

	"github.com/kalandramo/lulu"

	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveApp converts [App] into a slice of [lulu.Option].
func resolveApp(cfg *v1.App) []lulu.Option {
	var opts []lulu.Option

	if cfg.GetId() != "" {
		opts = append(opts, lulu.WithID(cfg.GetId()))
	}
	if cfg.GetName() != "" {
		opts = append(opts, lulu.WithName(cfg.GetName()))
	}
	if cfg.GetVersion() != "" {
		opts = append(opts, lulu.WithVersion(cfg.GetVersion()))
	}
	if dur := cfg.GetStopTimeout(); dur != nil {
		opts = append(opts, lulu.WithStopTimeout(dur.AsDuration()))
	} else if cfg.GetEnv() == "prd" {
		// 生产环境默认 30 秒优雅停机。
		opts = append(opts, lulu.WithStopTimeout(30*time.Second))
	}

	return opts
}

package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveConfig 检查 Config 配置中每个 optional 字段，
// 对已设置的配置源类型分别调用对应 action。
func resolveConfig(ctx context.Context, cfg *v1.Config) (func(), error) {
	type field struct {
		name   string
		action ConfigAction
	}

	var fields []field

	if cfg.GetFile() != nil {
		a, err := getConfigAction(ConfigTypeFile)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeFile, action: a})
	}
	if cfg.GetNacos() != nil {
		a, err := getConfigAction(ConfigTypeNacos)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeNacos, action: a})
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("bootstrap: no config source specified")
	}

	var cleanups []func()
	for _, f := range fields {
		cleanup, err := f.action(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, fmt.Errorf("bootstrap: build config %q: %w", f.name, err)
		}
		if cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
	}

	return cleanup, nil
}

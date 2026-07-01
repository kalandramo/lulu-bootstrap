package bootstrap

import (
	"fmt"

	"github.com/kalandramo/lulu/log"

	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func resolveLog(cfg *v1.Logger) (log.Logger, func(), error) {
	if cfg.GetType() == "" {
		return nil, nil, fmt.Errorf("bootstrap: logger type not specified")
	}

	b, err := getLogBuilder(cfg.GetType())
	if err != nil {
		return nil, nil, err
	}

	logger, cleanup, err := b(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("bootstrap: build logger %q: %w", cfg.GetType(), err)
	}
	return logger, cleanup, nil
}

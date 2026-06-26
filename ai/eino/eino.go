// Package eino provides a bootstrap AI builder for ByteDance Eino ChatModel.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/kalandramo/lulu-bootstrap/ai/eino"
package eino

import (
	"context"
	"fmt"

	einoPlugin "github.com/kalandramo/lulu-ext/ai/eino"

	bootstrap "github.com/kalandramo/lulu-bootstrap"
	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterAiBuilder(bootstrap.AiTypeEino, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Ai) (any, func(), error) {
	c := cfg.GetEino()
	if c == nil {
		return nil, nil, fmt.Errorf("eino: config is nil")
	}

	pluginCfg := &einoPlugin.Config{
		Type:      einoPlugin.ModelType(c.GetModelType()),
		ModelName: c.GetModelName(),
	}

	if ts := c.GetTimeoutSeconds(); ts > 0 {
		pluginCfg.TimeoutSeconds = ts
	}

	if cloud := c.GetCloud(); cloud != nil {
		pluginCfg.Cloud = &einoPlugin.CloudConfig{
			ApiKey:  cloud.GetApiKey(),
			BaseUrl: cloud.GetBaseUrl(),
		}
	}

	if local := c.GetLocal(); local != nil {
		pluginCfg.Local = &einoPlugin.LocalConfig{
			Host: local.GetHost(),
			Port: local.GetPort(),
		}
	}

	chatModel, err := einoPlugin.NewChatModel(ctx, pluginCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("eino: create chat model: %w", err)
	}

	return chatModel, func() {}, nil
}

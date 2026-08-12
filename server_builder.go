package bootstrap

import (
	"fmt"

	"github.com/kalandramo/lulu/transport"

	v1 "github.com/kalandramo/lulu-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveServer 检查 Server 配置中每个 optional 字段，
// 对已设置的传输层类型分别调用对应 builder 创建实例。
func resolveServer(cfg *v1.Server) ([]transport.Server, func(), error) {
	var servers []transport.Server

	// 按字段依次检查，调用对应 builder。
	type fieldBuilder struct {
		name string
		typ  string
		set  bool
	}
	fields := []fieldBuilder{
		{"http", serverTypeHTTP, cfg.GetHttp() != nil},
	}

	for _, f := range fields {
		if !f.set {
			continue
		}
		b, err := getServerBuilder(f.typ)
		if err != nil {
			return nil, nil, err
		}
		srv, err := b(cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: build server %q: %w", f.name, err)
		}
		if srv != nil {
			servers = append(servers, srv)
		}
	}

	return servers, func() {}, nil
}

// 内部常量，用于 server builder 注册的 key。
const (
	serverTypeHTTP         = "http"
)

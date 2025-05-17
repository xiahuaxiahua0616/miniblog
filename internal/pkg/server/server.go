package server

import (
	"context"
	"net/http"
)

type Server interface {
	// RunOrDie启动服务器
	RunOrDie()
	// GracefulStop 优雅关停
	GracefulStop(ctx context.Context)
}

// protocolName 从 http.Server 中获取协议名.
func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}
	return "http"
}

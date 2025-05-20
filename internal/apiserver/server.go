// Copyright 2025 xiahua <xhxiangshuijiao.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is github.com/xiahuaxiahua0616/miniblog. The professional

package apiserver

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"

	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/biz"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/pkg/contextx"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/store"
	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/log"
	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/server"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"
	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	ServerMode   string
	JWTKey       string
	Expiration   time.Duration
	HTTPOptions  *genericoptions.HTTPOptions
	GRPCOptions  *genericoptions.GRPCOptions
	MySQLOptions *genericoptions.MySQLOptions
}

// UnionServer 定义一个联合服务器. 根据 ServerMode 决定要启动的服务器类型.
//
// 联合服务器分为以下 2 大类：
//  1. Gin 服务器：由 Gin 框架创建的标准的 REST 服务器。根据是否开启 TLS，
//     来判断启动 HTTP 或者 HTTPS；
//  2. GRPC 服务器：由 gRPC 框架创建的标准 RPC 服务器
//  3. HTTP 反向代理服务器：由 grpc-gateway 框架创建的 HTTP 反向代理服务器。
//     根据是否开启 TLS，来判断启动 HTTP 或者 HTTPS；
//
// HTTP 反向代理服务器依赖 gRPC 服务器，所以在开启 HTTP 反向代理服务器时，会先启动 gRPC 服务器.
type UnionServer struct {
	cfg *Config
	srv server.Server
	lis net.Listener
}

type ServerConfig struct {
	cfg *Config
	biz biz.IBiz
}

// NewUnionServer 根据配置创建联合服务器.
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 注册租户解析函数，通过上下文获取用户ID
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	serverConfig, err := cfg.NewServerConfig()
	if err != nil {
		return nil, err
	}

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	// 根据服务模式创建对应的服务实例
	// 实际企业开发中，可以根据需要只选择一种服务器模式.
	// 这里为了方便给你展示，通过 cfg.ServerMode 同时支持了 Gin 和 GRPC 2 种服务器模式.
	var srv server.Server
	switch cfg.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}
	if err != nil {
		return nil, err
	}

	return &UnionServer{srv: srv}, nil
}

// Run 运行应用.
func (s *UnionServer) Run() error {
	go s.srv.RunOrDie()

	// 创建一个os.Signal类型的channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Infow("Shutdown down server...")

	ctx, channel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer channel()

	s.srv.GracefulStop(ctx)
	log.Infow("Server stopped")
	return nil
}

// NewServerConfig 创建一个 *ServerConfig 实例.
// 进阶：这里其实可以使用依赖注入的方式，来创建 *ServerConfig.
func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	// 初始化数据库连接
	db, err := cfg.NewDB()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)
	return &ServerConfig{
		cfg: cfg,
		biz: biz.NewBiz(store),
	}, nil
}

// NewDB 创建一个 *gorm.DB 实例.
func (cfg *Config) NewDB() (*gorm.DB, error) {
	return cfg.MySQLOptions.NewDB()
}

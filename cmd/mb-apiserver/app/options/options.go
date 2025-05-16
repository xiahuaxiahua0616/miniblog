// Copyright 2025 xiahua <xhxiangshuijiao.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is github.com/xiahuaxiahua0616/miniblog. The professional

package options

import (
	"errors"
	"fmt"
	"time"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
	stringsutil "github.com/onexstack/onexstack/pkg/util/strings"
	"github.com/spf13/pflag"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver"
	"k8s.io/apimachinery/pkg/util/sets"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// 定义支持的服务器模式集合.
var availableServerModes = sets.New(
	apiserver.GinServerMode,
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
)

// ServerOptions 包含服务器配置选项。
type ServerOptions struct {
	// ServerMode 定义服务器模式：gRPC、Gin HTTP、HTTP Reverse Proxy。
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	// JWTKey 定义 JWT 密钥。
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key"`
	// Expiration 定义 JWT Token 的过期时间。
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`
	// GRPCOptions 包含gRPC配置选项
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
}

func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:  apiserver.GRPCGatewayServerMode,
		JWTKey:      "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:  2 * time.Hour,
		GRPCOptions: genericoptions.NewGRPCOptions(),
	}

	opts.GRPCOptions.Addr = ":6666"
	return opts
}

// AddFlags 将 ServerOptions 的选项绑定到命令行标志。
// 通过使用 pflag 包，可以实现从命令行中解析这些选项的功能。
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")
	// 绑定 JWT Token 的过期时间选项到命令行标志。
	// 参数名称为 `--expiration`，默认值为 o.Expiration
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "The expiration duration of JWT tokens.")
	o.GRPCOptions.AddFlags(fs)
}

// Validate 校验 ServerOptions 中的选项是否合法.
func (o *ServerOptions) Validate() error {
	errs := []error{}

	// 校验 ServerMode 是否有效
	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: must be one of %v", availableServerModes.UnsortedList()))
	}

	// 校验 JWTKey 长度
	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("JWTKey must be at least 6 characters long"))
	}

	// 如果是gRPC或gRPC-Gateway模式，校验gRPC配置
	if stringsutil.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}

	// 合并所有错误并返回
	return utilerrors.NewAggregate(errs)
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:  o.ServerMode,
		JWTKey:      o.JWTKey,
		Expiration:  o.Expiration,
		GRPCOptions: o.GRPCOptions,
	}, nil
}

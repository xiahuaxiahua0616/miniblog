//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"

	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/biz"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/pkg/validation"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/store"
	ginmw "github.com/xiahuaxiahua0616/miniblog/internal/pkg/middleware/gin"
	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/server"
	"github.com/xiahuaxiahua0616/miniblog/pkg/auth"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginmw.UserRetriever), new(*UserRetriever)),
		),
		auth.ProviderSet,
	)
	return nil, nil
}

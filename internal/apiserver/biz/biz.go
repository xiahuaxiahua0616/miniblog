package biz

import (
	postv1 "github.com/xiahuaxiahua0616/miniblog/internal/apiserver/biz/v1/post"
	userv1 "github.com/xiahuaxiahua0616/miniblog/internal/apiserver/biz/v1/user"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/store"
)

type IBiz interface {
	UserV1() userv1.UserBiz

	PostV1() postv1.PostBiz
	// 获取帖子业务接口（V2版本）.
	// PostV2() post.PostBiz
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	store store.IStore
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(store store.IStore) *biz {
	return &biz{store: store}
}

// UserV1 返回一个实现了 UserBiz 接口的实例.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store)
}

// PostV1 返回一个实现了 PostBiz 接口的实例.
func (b *biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}

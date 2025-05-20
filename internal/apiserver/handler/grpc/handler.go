// Copyright 2025 xiahua <xhxiangshuijiao.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is github.com/xiahuaxiahua0616/miniblog. The professional

package grpc

import (
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/biz"
	apiv1 "github.com/xiahuaxiahua0616/miniblog/pkg/api/apiserver/v1"
)

type Handler struct {
	apiv1.UnimplementedMiniBlogServer

	biz biz.IBiz
}

func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}

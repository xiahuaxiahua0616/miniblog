// Copyright 2025 xiahua <xhxiangshuijiao.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is github.com/xiahuaxiahua0616/miniblog. The professional

package handler

import (
	"context"
	"time"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	apiv1 "github.com/xiahuaxiahua0616/miniblog/pkg/api/apiserver/v1"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}

// // Copyright 2025 xiahua <xhxiangshuijiao.com>. All rights reserved.
// // Use of this source code is governed by a MIT style
// // license that can be found in the LICENSE file. The original repo for
// // this file is github.com/xiahuaxiahua0616/miniblog. The professional

package main

import (
	"os"

	"github.com/xiahuaxiahua0616/miniblog/cmd/mb-apiserver/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	command := app.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

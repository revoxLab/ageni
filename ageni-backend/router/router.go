package router

import (
	"context"
	"github.com/readonme/open-studio/common/httpserver"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/router/handler"
)

func Init(c *conf.Config) {
	handler.StartHttp(c)
}

func Stop(ctx context.Context) {
	// stop http server
	httpserver.Stop(ctx)
	// stop grpc server
}

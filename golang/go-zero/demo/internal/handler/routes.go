// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/ntk148v/testing/golang/go-zero-test/demo/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/from/:name",
				Handler: DemoHandler(serverCtx),
			},
		},
	)
}
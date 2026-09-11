package controllers

import (
	"net/http"

	"github.com/sudzekai/web-os-api/server"
)

type healthController struct{}

var HealthController = healthController{}

func (ctr *healthController) Connect(srv *server.Server) {
	srv.AddHandler("GET /health", ctr.Get)
}

func (*healthController) Get(r *http.Request) server.MethodResult {
	return server.MethodResult{
		StatusCode: 200,
		Data:       "web-os-api alive",
	}
}

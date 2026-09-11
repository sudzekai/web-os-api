package controllers

import (
	"net/http"

	"github.com/sudzekai/web-os-api/server/abstractions"
	"github.com/sudzekai/web-os-api/server/types"
)

type healthController struct{}

var HealthController = healthController{}

func (ctr *healthController) Connect(srv abstractions.IServer) {
	srv.AddHandler("GET /health", ctr.Get)
}

func (*healthController) Get(r *http.Request) types.MethodResult {
	return types.MethodResult{
		StatusCode: 200,
		Data:       "web-os-api alive",
	}
}

package controllers

import (
	"net/http"

	"github.com/sudzekai/web-os-api/internal/objects/transfer"
)

func TestHandler(r *http.Request) transfer.MethodResult {
	return transfer.MethodResult{
		Data:       "Тест",
		StatusCode: 200,
	}
}

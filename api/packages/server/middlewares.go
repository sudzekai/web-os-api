package server

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type ResultFilterMiddleware func(http.ResponseWriter, *http.Request, MethodResult) error

type JwtMiddleware func(http.HandlerFunc, []string) http.HandlerFunc

package server

import "net/http"

type HandlerFunc func(*http.Request) MethodResult

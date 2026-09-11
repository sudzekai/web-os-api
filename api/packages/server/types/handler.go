package types

import "net/http"

type HandlerFunc func(*http.Request) MethodResult

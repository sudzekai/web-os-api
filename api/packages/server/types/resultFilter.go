package types

import (
	"net/http"
)

type ResultFilter func(http.ResponseWriter, *http.Request, MethodResult) error

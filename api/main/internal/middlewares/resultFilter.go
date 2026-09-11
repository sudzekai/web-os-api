package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sudzekai/web-os-api/internal/objects/responses"
	"github.com/sudzekai/web-os-api/server"
)

func ResultFilter(w http.ResponseWriter, r *http.Request, result server.MethodResult) error {
	response := responses.ResponseEnvelope{}

	response.IsSuccess = result.Error == nil

	response.Data = result.Data

	if result.Error != nil {
		response.Error = &responses.Error{
			StatusCode: result.StatusCode,
			Message:    result.Error.Error(),
		}
	} else {
		response.Error = nil
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(result.StatusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf(
			"Ошибка сериализации ответа на %s %s: %s",
			r.Method,
			r.URL.Path,
			err,
		)
	}
	return nil
}

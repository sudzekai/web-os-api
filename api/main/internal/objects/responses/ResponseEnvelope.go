package responses

type ResponseEnvelope struct {
	IsSuccess bool
	Data      any
	Error     Error
}

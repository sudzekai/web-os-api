package server

type MethodResult struct {
	Data       any
	Error      error
	StatusCode int
}

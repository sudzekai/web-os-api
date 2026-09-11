package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/sudzekai/web-os-api/logging/abstractions"
	"github.com/sudzekai/web-os-api/server/types"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux

	host string
	port int

	handlers map[string]http.HandlerFunc

	loggingProvider abstractions.ILogger

	jwtMiddleware types.JwtMiddleware
	resultFilter  types.ResultFilter

	middlewares map[string]types.Middleware

	isListening bool
	stats       types.ServerStats
}

func NewServer(host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		handlers:    make(map[string]http.HandlerFunc),
		middlewares: make(map[string]types.Middleware),
	}
}

func (srv *Server) Start() error {
	if srv.isListening {
		return fmt.Errorf("сервер уже запущен")
	}

	srv.mux = http.NewServeMux()

	srv.registerAllHandlers()

	srv.httpServer = &http.Server{
		Addr:    net.JoinHostPort(srv.host, strconv.Itoa(srv.port)),
		Handler: srv.buildMiddlewareChain(),
	}

	srv.isListening = true

	srv.stats.StartTime = time.Now()
	srv.stats.StopTime = time.Time{}

	srv.logInformation("сервер запущен и слушает " + srv.host + ":" + strconv.Itoa(srv.port))

	err := srv.httpServer.ListenAndServe()

	srv.isListening = false

	if err != nil && err != http.ErrServerClosed {
		srv.logError(
			"ошибка HTTP-сервера: %s", err,
		)
		return err
	}

	return nil
}

func (srv *Server) Stop() error {
	if !srv.isListening {
		return fmt.Errorf("сервер не запущен")
	}

	srv.stats.StopTime = time.Now()
	srv.isListening = false

	if err := srv.httpServer.Close(); err != nil {
		return fmt.Errorf(
			"ошибка остановки сервера: %w",
			err,
		)
	}

	srv.logInformation("Сервер остановлен")

	return nil
}

func (srv *Server) WaitForShutdown() error {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		os.Interrupt,
	)

	defer signal.Stop(signalChan)

	<-signalChan

	return srv.Stop()
}

func (srv *Server) IsListening() bool {
	return srv.isListening
}

func (srv *Server) AddMiddleware(name string, middleware types.Middleware) error {
	if srv.isListening {
		return fmt.Errorf("сервер уже запущен")
	}

	if middleware == nil {
		return fmt.Errorf("middleware %q равен nil", name)
	}

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("имя middleware не может быть пустым")
	}

	if _, exists := srv.middlewares[name]; exists {
		return fmt.Errorf("middleware %q уже зарегистрирован", name)
	}

	srv.middlewares[name] = middleware

	return nil
}

func (srv *Server) AddHandler(endpointPattern string, handler types.HandlerFunc) error {
	if err := srv.validateHandlerRegistration(endpointPattern, handler); err != nil {
		return err
	}

	srv.addHandler(endpointPattern, srv.createHandler(handler))

	return nil
}

func (srv *Server) AddProtectedHandler(endpointPattern string, handler types.HandlerFunc, roles []string) error {
	if err := srv.validateHandlerRegistration(endpointPattern, handler); err != nil {
		return err
	}

	if srv.jwtMiddleware == nil {
		return fmt.Errorf("JWT middleware не настроен")
	}

	protectedHandler := srv.jwtMiddleware(
		srv.createHandler(handler),
		roles,
	)

	srv.addHandler(endpointPattern, protectedHandler)

	return nil
}

// setters

func (srv *Server) SetLoggingProvider(provider abstractions.ILogger) {
	srv.loggingProvider = provider
}

func (srv *Server) SetJWTMiddleware(middleware types.JwtMiddleware) {
	srv.jwtMiddleware = middleware
}

func (srv *Server) SetResultFilter(filter types.ResultFilter) {
	srv.resultFilter = filter
}

// getters

func (srv *Server) GetEndpoints() []string {
	var result = make([]string, 0)

	for pattern := range srv.handlers {
		result = append(result, pattern)
	}

	return result
}

func (srv *Server) GetMiddlewares() []string {
	var result = make([]string, 0)

	for name := range srv.middlewares {
		result = append(result, name)
	}

	return result
}

func (srv *Server) GetStats() *types.ServerStats {
	return &srv.stats
}

func (srv *Server) buildMiddlewareChain() http.Handler {
	handler := http.Handler(srv.mux)

	for _, middleware := range srv.middlewares {
		handler = middleware(handler)
	}

	return handler
}

// handlers validation

func (srv *Server) validateHandlerRegistration(endpointPattern string, handler types.HandlerFunc) error {
	if srv.isListening {
		return fmt.Errorf("сервер уже запущен")
	}

	if handler == nil {
		return fmt.Errorf("handler равен nil")
	}

	fields := strings.Fields(endpointPattern)

	if len(fields) != 2 {
		return fmt.Errorf(
			"неверный формат паттерна %q, ожидается: GET /example/{route}",
			endpointPattern,
		)
	}

	return nil
}

// handlers addition + registration

func (srv *Server) registerAllHandlers() {
	srv.logInformation("регистрация обработчиков...")

	for pattern, handler := range srv.handlers {
		srv.mux.HandleFunc(pattern, handler)
		srv.logDebug("зарегистрирован обработчик для %s", pattern)
	}

	srv.logInformation("зарегистрировано %d обработчиков", len(srv.handlers))
}

func (srv *Server) addHandler(pattern string, hnd http.HandlerFunc) {
	srv.handlers[pattern] = hnd
	srv.logDebug("добавлен обработчик для %s", pattern)
}

func (srv *Server) createHandler(handler types.HandlerFunc) http.HandlerFunc {
	return func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		srv.stats.Requests.Add(1)

		result := handler(request)

		if result.Error != nil {
			srv.stats.Errors.Add(1)
		}

		if srv.writeResult(writer, request, result) {
			srv.stats.Responses.Add(1)
		}
	}
}

func (srv *Server) writeResult(writer http.ResponseWriter, request *http.Request, result types.MethodResult) bool {
	if srv.resultFilter != nil {
		err := srv.resultFilter(writer, request, result)

		if err != nil {
			return false
		}

		return true
	}

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	writer.WriteHeader(result.StatusCode)

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		srv.logError(
			"ошибка сериализации ответа на %s %s: %s",
			request.Method,
			request.URL.Path,
			err,
		)

		return false
	}

	return true
}

// LOGGING

func (srv *Server) logDebug(message string, args ...any) {
	if srv.loggingProvider != nil {
		srv.loggingProvider.LogDebug(message, args...)
	}
}

func (srv *Server) logInformation(message string, args ...any) {
	if srv.loggingProvider != nil {
		srv.loggingProvider.LogInformation(message, args...)
	}
}

func (srv *Server) logError(message string, args ...any) {
	if srv.loggingProvider != nil {
		srv.loggingProvider.LogError(message, args...)
	}
}

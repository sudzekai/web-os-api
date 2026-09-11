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
	"sync/atomic"
	"time"
)

type ServerStats struct {
	Requests  atomic.Int64
	Responses atomic.Int64
	Errors    atomic.Int64
	StartTime time.Time
	StopTime  time.Time
}

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux

	host string
	port int

	endpoints []string

	loggingProvider LoggingProvider
	jwtMiddleware   JwtMiddleware
	resultFilter    ResultFilterMiddleware
	middlewares     map[string]Middleware

	isListening bool
	stats       ServerStats
}

func NewServer(host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		mux:         http.NewServeMux(),
		middlewares: make(map[string]Middleware),
	}
}

func (srv *Server) Start() error {
	if srv.isListening {
		return fmt.Errorf("сервер уже запущен")
	}

	srv.httpServer = &http.Server{
		Addr:    net.JoinHostPort(srv.host, strconv.Itoa(srv.port)),
		Handler: srv.buildMiddlewareChain(),
	}

	srv.isListening = true

	srv.stats.StartTime = time.Now()
	srv.stats.StopTime = time.Time{}

	srv.logInformation("Сервер запущен и слушает " + srv.host + ":" + strconv.Itoa(srv.port))

	err := srv.httpServer.ListenAndServe()

	srv.isListening = false

	if err != nil && err != http.ErrServerClosed {
		srv.logError(
			fmt.Sprintf("Ошибка HTTP-сервера: %s", err),
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

func (srv *Server) AddLoggingProvider(provider LoggingProvider) {
	srv.loggingProvider = provider
}

func (srv *Server) SetJWTMiddleware(middleware JwtMiddleware) {
	srv.jwtMiddleware = middleware
}

func (srv *Server) SetResultFilter(filter ResultFilterMiddleware) {
	srv.resultFilter = filter
}

func (srv *Server) AddMiddleware(
	name string,
	middleware Middleware,
) error {
	if middleware == nil {
		return fmt.Errorf(
			"middleware %q равен nil",
			name,
		)
	}

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf(
			"имя middleware не может быть пустым",
		)
	}

	if srv.isListening {
		return fmt.Errorf("сервер уже запущен")
	}

	if _, exists := srv.middlewares[name]; exists {
		return fmt.Errorf(
			"middleware %q уже зарегистрирован",
			name,
		)
	}

	srv.middlewares[name] = middleware

	return nil
}

func (srv *Server) AddHandler(
	endpointPattern string,
	handler HandlerFunc,
) error {
	if err := srv.validateHandlerRegistration(endpointPattern); err != nil {
		return err
	}

	if handler == nil {
		return fmt.Errorf("handler равен nil")
	}

	srv.mux.HandleFunc(
		endpointPattern,
		srv.getHandler(handler),
	)

	srv.endpoints = append(
		srv.endpoints,
		endpointPattern,
	)

	srv.logDebug(
		fmt.Sprintf(
			"Добавлен обработчик для %s",
			endpointPattern,
		),
	)

	return nil
}

func (srv *Server) AddAuthorizedHandler(
	endpointPattern string,
	handler HandlerFunc,
	roles []string,
) error {
	if err := srv.validateHandlerRegistration(endpointPattern); err != nil {
		return err
	}

	if handler == nil {
		return fmt.Errorf("handler равен nil")
	}

	if srv.jwtMiddleware == nil {
		return fmt.Errorf("JWT middleware не настроен")
	}

	protectedHandler := srv.jwtMiddleware(
		srv.getHandler(handler),
		roles,
	)

	srv.mux.HandleFunc(
		endpointPattern,
		protectedHandler,
	)

	info := endpointPattern

	if len(roles) > 0 {
		info += " | JWT roles: " + strings.Join(roles, ", ")
	} else {
		info += " | JWT roles: any"
	}

	srv.endpoints = append(
		srv.endpoints,
		info,
	)

	srv.logDebug(
		fmt.Sprintf(
			"Добавлен авторизованный обработчик для %s",
			endpointPattern,
		),
	)

	return nil
}

func (srv *Server) GetEndpoints() []string {
	return append(
		[]string(nil),
		srv.endpoints...,
	)
}

func (srv *Server) GetStats() *ServerStats {
	return &srv.stats
}

func (srv *Server) buildMiddlewareChain() http.Handler {
	handler := http.Handler(srv.mux)

	for _, middleware := range srv.middlewares {
		handler = middleware(handler)
	}

	return handler
}

func (srv *Server) getHandler(
	handler HandlerFunc,
) http.HandlerFunc {
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

func (srv *Server) writeResult(
	writer http.ResponseWriter,
	request *http.Request,
	result MethodResult,
) bool {
	if srv.resultFilter != nil {
		srv.resultFilter(writer, result)
		return true
	}

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	writer.WriteHeader(result.StatusCode)

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		srv.logError(
			fmt.Sprintf(
				"Ошибка сериализации ответа на %s %s: %s",
				request.Method,
				request.URL.Path,
				err,
			),
		)

		return false
	}

	return true
}

func (srv *Server) validateHandlerRegistration(
	endpointPattern string,
) error {
	if srv.isListening {
		return fmt.Errorf("сервер уже запущен")
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

func (srv *Server) logDebug(message string) {
	if srv.loggingProvider != nil {
		srv.loggingProvider.LogDebug(message)
	}
}

func (srv *Server) logInformation(message string) {
	if srv.loggingProvider != nil {
		srv.loggingProvider.LogInformation(message)
	}
}

func (srv *Server) logError(message string) {
	if srv.loggingProvider != nil {
		srv.loggingProvider.LogError(message)
	}
}

package server

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/sudzekai/web-os-api/internal/utilities/logging"
)

var log logging.Logger = logging.NewLogger("server")

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux

	host string
	port int

	endpoints []string

	isStarted bool
}

func NewServer(host string, port int) *Server {
	return &Server{
		host:      host,
		port:      port,
		mux:       http.NewServeMux(),
		isStarted: false,
	}
}

func (server *Server) AddHandler(endpointPattern string, handler http.HandlerFunc) {
	if server.isStarted {
		log.LogError("Невозможно добавить обработчик для \"%s\": сервер уже запущен", endpointPattern)
		return
	}

	server.mux.HandleFunc(endpointPattern, handler)

	server.endpoints = append(server.endpoints, endpointPattern)

	log.LogDebug("Добавлен обработчик для \"%s\"", endpointPattern)
}

func (server *Server) Start() error {
	server.httpServer = &http.Server{
		Addr:    net.JoinHostPort(server.host, strconv.Itoa(server.port)),
		Handler: server.mux,
	}

	log.LogInformation("Сервер запущен и слушает %s:%d...", server.host, server.port)

	if len(server.endpoints) == 0 {
		log.LogWarning("Эндпоинты сервера отсутствуют")
	} else {
		log.LogInformation("Эндпоинты:\n\t%s", strings.Join(server.endpoints, "\n\t"))
	}

	server.isStarted = true
	return server.httpServer.ListenAndServe()
}

func (server *Server) Stop() error {
	if !server.isStarted {
		log.LogError("Невозможно остановить сервер: сервер не запущен")
		return nil
	}

	return server.httpServer.Close()
}

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/sudzekai/web-os-api/internal/objects/responses"
	"github.com/sudzekai/web-os-api/internal/objects/transfer"
	"github.com/sudzekai/web-os-api/internal/utilities/logging"
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

	jwtMiddleware func(http.HandlerFunc, string) http.HandlerFunc

	IsListening bool

	stats ServerStats
}

func NewServer(host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		mux:         http.NewServeMux(),
		IsListening: false,
	}
}

func (server *Server) AddHandler(endpointPattern string, hnd func(*http.Request) transfer.MethodResult) {
	log := logging.LoggerFactory.NewLogger("server:handler")

	if server.IsListening {
		log.LogError("Невозможно добавить обработчик для \"%s\": сервер уже запущен", endpointPattern)
		return
	}

	if len(strings.Split(endpointPattern, " ")) != 2 {
		log.LogError("Невозможно добавить обработчик для \"%s\": неверный формат паттерна. Шаблон: GET /example/{route}", endpointPattern)
		return
	}

	server.mux.HandleFunc(
		endpointPattern,
		getHandler(hnd, endpointPattern, server),
	)

	server.endpoints = append(server.endpoints, endpointPattern)

	log.LogDebug("Добавлен обработчик для \"%s\"", endpointPattern)
}

func (server *Server) AddAuthorizedHandler(endpointPattern string, hnd func(*http.Request) transfer.MethodResult, role string) {
	log := logging.LoggerFactory.NewLogger("server:handler")

	if server.IsListening {
		log.LogError("Невозможно добавить обработчик для \"%s\": сервер уже запущен", endpointPattern)
		return
	}

	if len(strings.Split(endpointPattern, " ")) != 2 {
		log.LogError("Невозможно добавить обработчик для \"%s\": неверный формат паттерна. Шаблон: GET /example/{route}", endpointPattern)
		return
	}

	server.mux.HandleFunc(endpointPattern, server.jwtMiddleware(getHandler(hnd, endpointPattern, server), role))

	server.endpoints = append(server.endpoints, fmt.Sprintf("%s | JWT role: %s", endpointPattern, role))

	log.LogDebug("Добавлен обработчик для \"%s\"", endpointPattern)
}

func (server *Server) Start() {
	log := logging.LoggerFactory.NewLogger("server:starter")

	if server.IsListening {
		log.LogError("Невозможно запустить сервер: сервер уже запущен")
		return
	}

	server.stats = ServerStats{
		StartTime: time.Now(),
	}

	server.httpServer = &http.Server{
		Addr:    net.JoinHostPort(server.host, strconv.Itoa(server.port)),
		Handler: server.mux,
	}

	if len(server.endpoints) == 0 {
		log.LogWarning("Эндпоинты сервера отсутствуют")
	}

	log.LogInformation(
		"Сервер запущен и слушает %s:%d...",
		server.host,
		server.port,
	)

	server.IsListening = true

	err := server.httpServer.ListenAndServe()

	server.IsListening = false

	if err != nil && err != http.ErrServerClosed {
		log.LogCritical("%s", err.Error())
		os.Exit(-1)
	}
}

func (server *Server) GetEndpoints() []string {
	return server.endpoints
}

func (server *Server) Stop() {
	log := logging.LoggerFactory.NewLogger("server:stopper")

	if !server.IsListening {
		log.LogError("Невозможно остановить сервер: сервер не запущен")
		return
	}

	server.stats.StopTime = time.Now()

	err := server.httpServer.Close()

	if err != nil {
		log.LogCritical("%s", err.Error())
		os.Exit(-1)
	} else {
		log.LogInformation("Сервер остановлен")
	}

	server.IsListening = false
}

func (server *Server) WaitForShutdown() {
	log := logging.LoggerFactory.NewLogger("server:waiter")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	<-ctx.Done()

	log.LogInformation("Остановка сервера...")

	server.Stop()
}

func getHandler(fun func(*http.Request) transfer.MethodResult, endpointPattern string, server *Server) http.HandlerFunc {
	endpoint := strings.Split(endpointPattern, " ")[1]
	hndLog := logging.LoggerFactory.NewLogger("handler:" + endpoint)

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ip := r.Header.Get("X-Forwarded-For")

		if ip == "" {
			ip, _, err = net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				hndLog.LogWarning(
					"Ошибка получения IP-адреса клиента: %s",
					err.Error(),
				)
				ip = r.RemoteAddr
			}
		}

		hndLog.LogInformation(
			"Получен запрос: %s => %s %s?%s",
			ip,
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
		)
		server.stats.Requests.Add(1)

		result := fun(r)

		success := result.Error == nil

		response := responses.ResponseEnvelope{
			IsSuccess: success,
			Data:      result.Data,
		}

		if result.Error != nil {
			hndLog.LogError("%s", result.Error.Error())

			response.Error = responses.Error{
				StatusCode: result.StatusCode,
				Message:    result.Error.Error(),
			}
			server.stats.Errors.Add(1)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(result.StatusCode)
		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			hndLog.LogError("Ошибка сериализации ответа: %s", err.Error())
		} else {
			hndLog.LogInformation("Отправлен ответ на запрос: %d", result.StatusCode)
			server.stats.Responses.Add(1)
		}
	}
}

func (server *Server) GetStats() *ServerStats {
	return &server.stats
}

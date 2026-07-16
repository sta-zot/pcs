package provisioning

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"strconv"
	"time"
)

// UseCase for provisioning configuration files
type provisioning interface {
	Provision(ctx context.Context, fileName, userAgent, ip string) ([]byte, error)
}

type registrator interface {
	Get(pattern string, handler http.HandlerFunc)
}

type Handler struct {
	pService provisioning
	reg      registrator
}

// Получение запроса конфигурации и возврат конфигурации
func (s *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fileName := r.PathValue("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}
	var ip string = ""

	if ip = r.Header.Get("X-Real-IP"); ip != "" {
	} else {
		ip = r.RemoteAddr
	}
	ip, _, _ = net.SplitHostPort(ip)

	config, err := s.pService.Provision(ctx, fileName, userAgent, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
	w.Header().Set("Content-Length", strconv.Itoa(len(config)))
	http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(config))
}

// Регистрация маршрутов для обработки запросов конфигурации
// Регистрирует маршруты для обработчика, маршруты добавляются к базовому пути
// переданному в registrator, например, "/api/v1/config"
func (s *Handler) registerRoutes() {
	//BasePath in route + path in the handler
	// For example, /api/v1/config/{file}
	s.reg.Get("/{file}", s.GetConfig)
}

func NewHandler(pService provisioning, registrator registrator) *Handler {
	handler := &Handler{
		pService: pService,
		reg:      registrator,
	}
	handler.registerRoutes()
	return handler
}

package provisioning

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type provisionService interface {
	GetConfig(mac, ip string, info ...string) ([]byte, error)
}

type registrator interface {
	Get(pattern string, handler http.HandlerFunc)
}

type Service struct {
	pService provisionService
	reg      registrator
}

func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	mac := r.PathValue("mac")
	ua := r.Header.Get("User-Agent")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	query := r.URL.Query()
	var err error
	//err := fmt.Errorf("Config not found")
	type response struct {
		Config []byte `json:"config"`
		Err    error  `json:"error"`
	}
	fileContent := fmt.Sprintf("MAC Address : \t %s\nUser-Agent : \t %s\nIp Address : %s", mac, ua, ip)
	for key, values := range query {
		for _, value := range values {
			fileContent = fmt.Sprintf("%s\n%s : \t %s", fileContent, key, value)
		}
	}
	fileStream := strings.NewReader(fileContent)
	if err != nil {
		res := response{Err: fmt.Errorf("Config not found")}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(res)
		return
	}
	fileName := strings.ReplaceAll(mac, ":", "") + ".cfg"
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	http.ServeContent(w, r, fileName, time.Now(), fileStream)

}

func (s *Service) ListHttpHeaders(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string][]string)
	for key, values := range r.Header {
		headers[key] = values
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		headers["X-Real-IP"] = []string{ip}
	} else {
		headers["X-Real-IP"] = []string{r.RemoteAddr}
	}
	json.NewEncoder(w).Encode(headers)
}

func (s *Service) registerRoutes() {
	s.reg.Get("/{mac}.cfg", s.GetConfig)
	s.reg.Get("/headers", s.ListHttpHeaders)
}

func NewHandler(pService provisionService, reg registrator) *Service {
	svc := &Service{
		pService: pService,
		reg:      reg,
	}
	svc.registerRoutes()
	return svc
}

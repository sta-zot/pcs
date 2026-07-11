package registrator

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// API interface, который мы будем отдавать хэндлерам.
// Они знают только про методы регистрации, но не знают, какой префикс у них СЕЙЧАС.
type IRegistrator interface {
	Get(pattern string, handler http.HandlerFunc)
	Post(pattern string, handler http.HandlerFunc)
	// Метод для создания дочернего регистратора (например, для ухода вглубь: /api -> /api/v1)
	WithPrefix(prefix string) *Registrator
}

var _ = IRegistrator(&Registrator{})

// Конкретная реализация
type Registrator struct {
	root   chi.Router
	router chi.Router
}

// New создает корневой регистратор приложения
func New(baseURL string) *Registrator {
	root := chi.NewRouter()
	if baseURL == "" {
		baseURL = "/"
	}
	if baseURL[0] != '/' {
		baseURL = "/" + baseURL
	}
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	var subRouter chi.Router
	root.Route(baseURL, func(r chi.Router) {
		subRouter = r
	})
	return &Registrator{router: subRouter, root: root}
}

func (reg *Registrator) Get(pattern string, handler http.HandlerFunc) {
	reg.router.Get(pattern, handler)
}

func (reg *Registrator) Post(pattern string, handler http.HandlerFunc) {
	reg.router.Post(pattern, handler)
}

func (reg *Registrator) Put(pattern string, handler http.HandlerFunc) {
	reg.router.Put(pattern, handler)
}

func (reg *Registrator) Delete(pattern string, handler http.HandlerFunc) {
	reg.router.Delete(pattern, handler)
}

func (reg *Registrator) Patch(pattern string, handler http.HandlerFunc) {
	reg.router.Patch(pattern, handler)
}

func (reg *Registrator) Handler() http.Handler {
	return reg.root
}

// WithPrefix — это ключевой метод. Он берет текущий роутер, делает из него
// под-роутер  через Route() и возвращает НОВЫЙ регистратор, замкнутый на эту ветку.
func (reg *Registrator) WithPrefix(prefix string) *Registrator {
	var subRouter chi.Router

	// Используем стандартный chi.Route, чтобы создать изолированную ветку путей
	reg.router.Route(prefix, func(r chi.Router) {
		subRouter = r
	})

	return &Registrator{router: subRouter}
}

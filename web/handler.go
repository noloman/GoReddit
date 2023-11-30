package web

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/noloman/goreddit"
)

type Handler struct {
	*chi.Mux
	goreddit.Store
}

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		Store: store,
	}
	h.Use(middleware.Logger)
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadsList())
	})
	return h
}

var ThreadsHtml = `
<h1>Threads</h1>
<dl>
	{{range .Threads}}
		<dt><strong>{{.Title}}</strong></dt>
		<dd>{{.Description}}</dd>
	{{end}}
</dl>
`

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}
	tmpl := template.Must(template.New("").Parse(ThreadsHtml))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.Store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Threads: tt})
	}
}

package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/noloman/goreddit"
)

type ThreadsHandler struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *ThreadsHandler) List() http.HandlerFunc {
	type data struct {
		SessionData
		Threads []goreddit.Thread
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/threads.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Threads: tt, SessionData: GetSessionData(h.sessions, r.Context())})
	}
}

func (h *ThreadsHandler) Create() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{CSRF: csrf.TemplateField(r), SessionData: GetSessionData(h.sessions, r.Context())})
	}
}

func (h *ThreadsHandler) Show() http.HandlerFunc {
	type data struct {
		SessionData
		Thread goreddit.Thread
		Posts  []goreddit.Post
		CSRF   template.HTML
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pp, err := h.store.PostsByThread(t.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Thread: t, Posts: pp, CSRF: csrf.TemplateField(r), SessionData: GetSessionData(h.sessions, r.Context())})
	}
}

func (h *ThreadsHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := CreateThreadForm{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Errors: 	FormErrors{},
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}
		if err := h.store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       form.Title,
			Description: form.Description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessions.Put(r.Context(), "flash", "Thread created successfully!")
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *ThreadsHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessions.Put(r.Context(), "flash", "Thread deleted successfully!")
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

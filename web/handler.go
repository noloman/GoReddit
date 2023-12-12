package web

import (
	"context"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/noloman/goreddit"
)

type Handler struct {
	*chi.Mux
	store    goreddit.Store
	sessions *scs.SessionManager
}

func NewHandler(store goreddit.Store, sessions *scs.SessionManager, csrfKey []byte) *Handler {
	h := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	threads := ThreadsHandler{store: store, sessions: sessions}
	comments := CommentsHandler{store: store, sessions: sessions}
	posts := PostsHandler{store: store, sessions: sessions}
	users := UserHandler{store: store, sessions: sessions}

	h.Use(middleware.Logger)
	h.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	h.Use(sessions.LoadAndSave)
	h.Use(h.withUser)

	h.Get("/", h.Home())
	h.Route("/threads", func(r chi.Router) {
		// THREADS
		r.Get("/", threads.List())
		r.Get("/new", threads.Create())
		r.Post("/", threads.Store())
		r.Post("/{id}/delete", threads.Delete())
		r.Get("/{id}", threads.Show())

		// POSTS
		r.Post("/{id}", posts.Store())
		r.Get("/{threadID}/new", posts.Create())
		r.Get("/{threadID}/{postID}", posts.Show())

		// COMMENTS
		r.Post("/{threadID}/{postID}", comments.Store())
	})
	h.Get("/comments/{id}/vote", comments.Vote())
	h.Get("/posts/{id}/vote", posts.Vote()) // TODO Should this be here or inside /threads?
	// REGISTRATION
	h.Get("/register", users.Register())
	h.Post("/register", users.RegisterSubmit())
	// LOGIN
	h.Get("/login", users.Login())
	h.Post("/login", users.LoginSubmit())
	// LOGOUT
	h.Get("/logout", users.Logout())
	return h
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		SessionData
		Posts []goreddit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, &data{
			Posts:       pp,
			SessionData: GetSessionData(h.sessions, r.Context()),
		})
	}
}

func (h *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := h.sessions.Get(r.Context(), "user_id").(uuid.UUID)
		user, err := h.store.User(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

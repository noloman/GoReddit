package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/noloman/goreddit"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *UserHandler) Register() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/user_register.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{CSRF: csrf.TemplateField(r), SessionData: GetSessionData(h.sessions, r.Context())})
	}
}

func (h *UserHandler) RegisterSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := RegisterUserForm{
			Username:      r.FormValue("username"),
			Password:      r.FormValue("password"),
			UsernameTaken: false,
		}
		if _, err := h.store.UserByUsername(form.Username); err == nil {
			form.UsernameTaken = true
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.store.CreateUser(&goreddit.User{
			ID:       uuid.New(),
			Username: form.Username,
			Password: string(password),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessions.Put(r.Context(), "flash", "You have successfully registered! Please login.")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

package web

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/noloman/goreddit"
	"net/http"
)

type CommentsHandler struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *CommentsHandler) Vote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c, err := h.store.Comment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			c.Votes++
		} else {
			c.Votes--
		}

		if err := h.store.UpdateComment(&c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (h *CommentsHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := CreateCommentForm{
			Content: r.FormValue("content"),
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}
		idStr := chi.URLParam(r, "postID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.store.CreateComment(
			&goreddit.Comment{
				ID:      uuid.New(),
				PostID:  id,
				Content: form.Content,
			}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.sessions.Put(r.Context(), "flash", "Comment created successfully!")
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

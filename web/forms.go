package web

import "encoding/gob"

func init() {
	gob.Register(CreatePostForm{})
	gob.Register(CreateThreadForm{})
	gob.Register(CreateCommentForm{})
	gob.Register(FormErrors{})
}

type FormErrors map[string]string

type CreatePostForm struct {
	Title   string
	Content string
	Errors  FormErrors
}

type CreateThreadForm struct {
	Title       string
	Description string
	Errors      FormErrors
}

type CreateCommentForm struct {
	Content string
	Errors  FormErrors
}

// Validate validates the post form with Title and Content
func (f *CreatePostForm) Validate() bool {
	if f.Title == "" {
		f.Errors["Title"] = "Title is required."
	}
	if f.Content == "" {
		f.Errors["Content"] = "Content is required."
	}

	return len(f.Errors) == 0
}

// Validate validates the thread form with Title and Description
func (f *CreateThreadForm) Validate() bool {
	if f.Title == "" {
		f.Errors["Title"] = "Title is required."
	}
	if f.Description == "" {
		f.Errors["Description"] = "Description is required."
	}

	return len(f.Errors) == 0
}

// Validate validates the comment form with Content
func (f *CreateCommentForm) Validate() bool {
	if f.Content == "" {
		f.Errors["Content"] = "Content is required."
	}

	return len(f.Errors) == 0
}

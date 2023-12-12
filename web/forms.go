package web

import "encoding/gob"

func init() {
	gob.Register(CreatePostForm{})
	gob.Register(CreateThreadForm{})
	gob.Register(CreateCommentForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
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

type RegisterForm struct {
	Username      string
	Password      string
	UsernameTaken bool
	Errors        FormErrors
}

type LoginForm struct {
	Username             string
	Password             string
	IncorrectCredentials bool
	Errors               FormErrors
}

// Validate validates the post form with Title and Content
func (f *CreatePostForm) Validate() bool {
	f.Errors = FormErrors{}
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
	f.Errors = FormErrors{}
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
	f.Errors = FormErrors{}
	if f.Content == "" {
		f.Errors["Content"] = "Content is required."
	}

	return len(f.Errors) == 0
}

// Validate validates the user registration form with Username and Password
func (f *RegisterForm) Validate() bool {
	f.Errors = FormErrors{}
	if f.Username == "" {
		f.Errors["Username"] = "Username is required."
	} else if f.UsernameTaken {
		f.Errors["Username"] = "Username is already taken."
	} else if len(f.Username) < 3 {
		f.Errors["Username"] = "Username must be at least 3 characters."
	} else if len(f.Username) > 20 {
		f.Errors["Username"] = "Username must be at most 20 characters."
	}

	if f.Password == "" {
		f.Errors["Password"] = "Password is required."
	} else if len(f.Password) < 8 {
		f.Errors["Password"] = "Password must be at least 8 characters."
	} else if len(f.Password) > 50 {
		f.Errors["Password"] = "Password must be at most 50 characters."
	} else if f.Password == f.Username {
		f.Errors["Password"] = "Password must not be the same as the Username."
	} else if f.Password == "password" {
		f.Errors["Password"] = "Password must not be 'password'."
	}

	return len(f.Errors) == 0
}

// Validate validates the user login form with Username and Password
func (f *LoginForm) Validate() bool {
	f.Errors = FormErrors{}
	if f.Username == "" {
		f.Errors["Username"] = "Username is required."
	} else if f.IncorrectCredentials {
		f.Errors["Username"] = "User and/or password are incorrect."
	}

	if f.Password == "" {
		f.Errors["Password"] = "Please enter a password."
	}

	return len(f.Errors) == 0
}

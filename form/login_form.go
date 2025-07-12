package form

import "net/mail"

type LoginForm struct {
	Email    string
	Password string

	Errors map[string]string
}

func (f *LoginForm) Validate() bool {
	f.Errors = map[string]string{}

	if f.Email == "" {
		f.Errors["Email"] = "The email is required"
	}

	if f.Email == "" {
		f.Errors["Email"] = "The email is required"
	} else if _, err := mail.ParseAddress(f.Email); err != nil {
		f.Errors["Email"] = "The email format is invalid"
	}

	if f.Password == "" {
		f.Errors["Password"] = "The password is required"
	}

	return len(f.Errors) == 0
}

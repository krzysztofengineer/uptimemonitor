package form

import "net/mail"

type SetupForm struct {
	Name     string
	Email    string
	Password string

	Errors map[string]string
}

func (f *SetupForm) Validate() bool {
	f.Errors = map[string]string{}

	if f.Name == "" {
		f.Errors["Name"] = "The name field is required"
	}

	if f.Email == "" {
		f.Errors["Email"] = "The email field is required"
	} else if _, err := mail.ParseAddress(f.Email); err != nil {
		f.Errors["Email"] = "The email format is invalid"
	}

	if f.Password == "" {
		f.Errors["Password"] = "The password field is required"
	}

	return len(f.Errors) == 0
}

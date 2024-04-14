package user

import "regexp"

var alphaNumericRegex *regexp.Regexp

type CreateUserOpts struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type GetUserByIdOpts struct {
	Id string `json:"id"`
}

func init() {
	alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9_.]*$")
}

func (c CreateUserOpts) Validate() bool {
	if c.Username == "" {
		return false
	}

	if !alphaNumericRegex.MatchString(c.Username) {
		return false
	}

	if c.Password == "" {
		return false
	}

	return true
}

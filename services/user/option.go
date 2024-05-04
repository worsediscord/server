package user

type CreateUserOpts struct {
	Username string
	Password string
}

type GetUserByIdOpts struct {
	Id string
}

type DeleteUserOpts struct {
	Id string
}

func (c CreateUserOpts) Validate() error {
	if c.Username == "" {
		return ErrInvalidUsername
	}

	if len(c.Password) < 8 {
		return ErrInvalidPassword
	}

	return nil
}

package user

type CreateUserOpts struct {
	Username string
	Password string
}

type GetUserByIdOpts struct {
	Id string
}

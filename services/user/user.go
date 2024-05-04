package user

type User struct {
	Username string `json:"username,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Password string `json:"password,omitempty"`
}

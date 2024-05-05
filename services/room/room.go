package room

type Room struct {
	Id     int64
	Name   string
	Users  []string
	Admins []string
}

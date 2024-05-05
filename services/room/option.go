package room

type CreateRoomOpts struct {
	Name   string
	UserId string
}

type GetRoomByIdOpts struct {
	Id int64
}

type DeleteRoomOpts struct {
	Id     int64
	UserId string
	Force  bool
}

type JoinRoomOpts struct {
	Id     int64
	UserId string
}

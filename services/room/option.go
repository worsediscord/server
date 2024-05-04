package room

type CreateRoomOpts struct {
	Name string
}

type GetRoomByIdOpts struct {
	Id int64
}

type DeleteRoomOpts struct {
	Id int64
}

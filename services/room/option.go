package room

type CreateRoomOpts struct {
	Name string `json:"name"`
}

type GetRoomByIdOpts struct {
	Id int64 `json:"id,omitempty"`
}

func (c CreateRoomOpts) Validate() bool {
	return c.Name != ""
}

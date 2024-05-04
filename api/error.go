package api

// swagger:response Error
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

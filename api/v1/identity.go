package v1

type Identity struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Identifiable interface {
	CommonName() string
	UID() string
}

func (i Identity) CommonName() string {
	return i.Name
}

func (i Identity) UID() string {
	return i.ID
}

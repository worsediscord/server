package v2

type Identity struct {
	Name string
	ID   string
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

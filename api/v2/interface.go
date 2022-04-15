package v2

type Identifiable interface {
	CommonName() string
	UID() string
}

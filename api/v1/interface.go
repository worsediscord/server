package v1

type Identifiable interface {
	CommonName() string
	UID() string
}

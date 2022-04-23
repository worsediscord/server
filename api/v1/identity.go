package v2

type Identity struct {
	Name string
	UID  string
}

type Identifiable interface {
	CommonName() string
	UID() string
}

package storage

type Reader[K comparable, V any] interface {
	Read(K) (V, error)
	ReadAll() ([]V, error)
}

type Writer[K comparable, V any] interface {
	Write(K, V) error
}

type ReadWriter[K comparable, V any] interface {
	Reader[K, V]
	Writer[K, V]
}

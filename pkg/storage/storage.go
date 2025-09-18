package storage

type storage interface {
	Set(key string, value string) error
	SetDocument(key string, value any) error
	Get(key string) ([]byte, error)
}

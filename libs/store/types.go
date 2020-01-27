package store

//StoreInterface -
type StoreInterface interface {
	Load([]byte) ([]byte, error)
	Save([]byte, []byte) error
}

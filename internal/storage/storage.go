package storage

type Storage interface {
	// GetUploadLink will request an upload link for the given key in the storage.
	// If the upload fails, returns an error
	GetUploadLink(key string) (string, string, error)

	// Delete will delete the given key from the storage.
	// If the deletion fails, returns an error
	DeleteObject(key string) error

	Healthcheck() error
}

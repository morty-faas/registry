package storage

type Storage interface {
	// GetUploadLink will request an upload link for the given key in the storage.
	// If the upload fails, returns an error
	GetUploadLink(key string) (string, string, error)

	Healthcheck() error
}

package storage

import "io"

type Storage interface {
	// GetUploadLink will request an upload link for the given key in the storage.
	// If the upload fails, returns an error
	GetUploadLink(key string) (string, string, error)

	// PutFile insert a new file at the given key in the storage
	PutFile(key string, body io.Reader) error

	Healthcheck() error
}

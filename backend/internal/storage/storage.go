package storage

import (
	"io"
)

// Storage abstracts local disk vs future cloud (S3, etc.)
type Storage interface {
	Save(categorySlug, filename string, reader io.Reader) (relativePath string, err error)
	Delete(relativePath string) error
	PublicURL(relativePath string) string
	EnsureCategoryDir(categorySlug string) error
}

package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	baseDir    string
	publicBase string
}

func NewLocalStorage(baseDir, publicBase string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{baseDir: baseDir, publicBase: strings.TrimRight(publicBase, "/")}, nil
}

func (s *LocalStorage) categoryPath(slug string) string {
	return filepath.Join(s.baseDir, slug)
}

func (s *LocalStorage) videoPath(slug string) string {
	return filepath.Join(s.baseDir, slug, "videos")
}

func (s *LocalStorage) EnsureCategoryDir(categorySlug string) error {
	if err := os.MkdirAll(s.categoryPath(categorySlug), 0755); err != nil {
		return err
	}
	return os.MkdirAll(s.videoPath(categorySlug), 0755)
}

func (s *LocalStorage) Save(categorySlug, filename string, reader io.Reader) (string, error) {
	if err := s.EnsureCategoryDir(categorySlug); err != nil {
		return "", err
	}
	dest := filepath.Join(s.categoryPath(categorySlug), filename)
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}
	rel := filepath.ToSlash(filepath.Join(categorySlug, filename))
	return rel, nil
}

func (s *LocalStorage) SaveVideo(categorySlug, filename string, reader io.Reader) (string, error) {
	if err := s.EnsureCategoryDir(categorySlug); err != nil {
		return "", err
	}
	dest := filepath.Join(s.videoPath(categorySlug), filename)
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}
	rel := filepath.ToSlash(filepath.Join(categorySlug, "videos", filename))
	return rel, nil
}

func (s *LocalStorage) SaveRelative(relativePath string, reader io.Reader) error {
	rel := filepath.ToSlash(relativePath)
	dest := filepath.Join(s.baseDir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

func (s *LocalStorage) Delete(relativePath string) error {
	full := filepath.Join(s.baseDir, filepath.FromSlash(relativePath))
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(full)
}

func (s *LocalStorage) PublicURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.publicBase, strings.TrimLeft(filepath.ToSlash(relativePath), "/"))
}

func (s *LocalStorage) BaseDir() string { return s.baseDir }

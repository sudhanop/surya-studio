package repositories

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

type SettingsRepo struct {
	file string
	mu   sync.Mutex
}

func NewSettingsRepo(dataDir string) *SettingsRepo {
	return &SettingsRepo{file: filepath.Join(dataDir, "settings.json")}
}

func (r *SettingsRepo) load() (map[string]string, error) {
	var out map[string]string
	err := readJSONFile(r.file, &out)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	if out == nil {
		out = map[string]string{}
	}
	return out, nil
}

func (r *SettingsRepo) save(m map[string]string) error {
	return writeJSONFile(r.file, m)
}

func (r *SettingsRepo) GetAll(ctx context.Context) (map[string]string, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.load()
}

func (r *SettingsRepo) Set(ctx context.Context, key, value string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	m, err := r.load()
	if err != nil {
		return err
	}
	m[key] = value
	return r.save(m)
}

func (r *SettingsRepo) SetMany(ctx context.Context, pairs map[string]string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	m, err := r.load()
	if err != nil {
		return err
	}
	for k, v := range pairs {
		m[k] = v
	}
	return r.save(m)
}

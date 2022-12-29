package data

import (
	"encoding/json"
	"github.com/timshannon/bolthold"
	"os"
	"path/filepath"
)

const (
	boltDBFileName = ".tasklist.db"
	exportPattern  = "taskstore_*.snapshot"
)

var _ TaskStorage = (*boltStorage)(nil)

type boltStorage struct {
	store *bolthold.Store
}

func newBoltStorage() (*boltStorage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return newBoltStorageAt(filepath.Join(home, boltDBFileName))
}

func newBoltStorageAt(file string) (*boltStorage, error) {
	store, err := bolthold.Open(file, 0700, &bolthold.Options{
		Encoder: func(value any) ([]byte, error) {
			return json.Marshal(value)
		},
		Decoder: func(data []byte, value any) error {
			return json.Unmarshal(data, value)
		},
	})
	if err != nil {
		return nil, err
	}
	return &boltStorage{store: store}, nil
}

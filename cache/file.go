package cache

import (
	"encoding/json"
	"go-dyndns/updater"
	"os"
	"path"
)

type fileCache struct {
	path string
}

func newFileCache(path string) *fileCache {
	return &fileCache{path: path}
}

func (f *fileCache) GetLastRequest() (*updater.UpdateRequest, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0755)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	var req updater.UpdateRequest
	if err = decoder.Decode(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (f *fileCache) SetLastRequest(req *updater.UpdateRequest) error {
	// Create directory if not exists
	dir := path.Dir(f.path)
	if err := os.MkdirAll(dir, 1777); err != nil {
		return err
	}

	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666) // The permissions of the beast
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(req)
}

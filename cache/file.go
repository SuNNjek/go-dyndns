package cache

import (
	"io"
	"net"
	"os"
)

type fileCache struct {
	path string
}

func newFileCache(path string) *fileCache {
	return &fileCache{path: path}
}

func (f *fileCache) GetLastIp() (net.IP, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0755)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var ip net.IP
	if err = ip.UnmarshalText(content); err != nil {
		return nil, err
	}

	return ip, nil
}

func (f *fileCache) SetLastIp(ip net.IP) error {
	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	newContent, err := ip.MarshalText()
	if err != nil {
		return err
	}

	_, err = file.Write(newContent)
	return err
}

package cache

import (
	"go-dyndns/updater"
)

type memoryCache struct {
	value *updater.UpdateRequest
}

func (m *memoryCache) GetLastRequest() (*updater.UpdateRequest, error) {
	return m.value, nil
}

func (m *memoryCache) SetLastRequest(req *updater.UpdateRequest) error {
	m.value = req
	return nil
}

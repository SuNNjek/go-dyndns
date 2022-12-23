package cache

import (
	"github.com/stretchr/testify/mock"
	"go-dyndns/updater"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetLastRequest() (*updater.UpdateRequest, error) {
	args := m.Called()

	if req, ok := args.Get(0).(*updater.UpdateRequest); ok {
		return req, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *MockCache) SetLastRequest(req *updater.UpdateRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

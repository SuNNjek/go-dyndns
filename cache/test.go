package cache

import (
	"github.com/stretchr/testify/mock"
	"net"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetLastIp() (net.IP, error) {
	args := m.Called()

	if ip, ok := args.Get(0).(net.IP); ok {
		return ip, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *MockCache) SetLastIp(ip net.IP) error {
	args := m.Called(ip)
	return args.Error(0)
}

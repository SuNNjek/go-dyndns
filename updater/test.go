package updater

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net"
)

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) UpdateIP(ctx context.Context, addr net.IP) error {
	args := m.Called(ctx, addr)
	return args.Error(0)
}

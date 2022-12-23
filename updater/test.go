package updater

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) UpdateIP(ctx context.Context, req *UpdateRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

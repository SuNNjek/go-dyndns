package addrproviders

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go-dyndns/util"
	"net"
)

var (
	httpClientMock = &util.MockHttpClient{}
)

type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) GetIP(ctx context.Context) (net.IP, error) {
	args := m.Called(ctx)

	if ip, ok := args.Get(0).(net.IP); ok {
		return ip, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

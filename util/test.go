package util

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)

	if resp, ok := args.Get(0).(*http.Response); ok {
		return resp, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	if resp, ok := args.Get(0).(*http.Response); ok {
		return resp, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

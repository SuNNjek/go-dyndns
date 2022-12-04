package log

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Log(level LogLevel, format string, args ...interface{}) {
	m.Called(level, format, args)
}

func (m *MockLogger) Panic(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatal(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Error(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warn(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Info(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Trace(format string, args ...interface{}) {
	m.Called(format, args)
}

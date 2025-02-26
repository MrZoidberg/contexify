package mocks

import (
	"os"
	"sync"
)

// MockFile represents a mock implementation of os.File
type MockFile struct {
	WriteFunc func(b []byte) (n int, err error)
	CloseFunc func() error
}

// Write calls the WriteFunc if defined or returns default values
func (m *MockFile) Write(b []byte) (n int, err error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(b)
	}
	return len(b), nil
}

// Close calls the CloseFunc if defined or returns nil
func (m *MockFile) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

var (
	// DefaultCreateFunc is the default implementation of Create
	DefaultCreateFunc = os.Create
	// MockCreateFunc can be set to mock the Create function
	MockCreateFunc func(name string) (*os.File, error)
	mockMutex      sync.RWMutex
)

// Create is a mockable version of os.Create
func Create(name string) (*os.File, error) {
	mockMutex.RLock()
	mock := MockCreateFunc
	mockMutex.RUnlock()

	if mock != nil {
		return mock(name)
	}
	return DefaultCreateFunc(name)
}

// SetMockCreate sets the mock function for Create
func SetMockCreate(mock func(name string) (*os.File, error)) {
	mockMutex.Lock()
	MockCreateFunc = mock
	mockMutex.Unlock()
}

// ResetMockCreate resets the mock function for Create
func ResetMockCreate() {
	mockMutex.Lock()
	MockCreateFunc = nil
	mockMutex.Unlock()
}

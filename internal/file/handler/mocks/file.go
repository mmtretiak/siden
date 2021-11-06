// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/file/handler/file.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileService is a mock of FileService interface.
type MockFileService struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceMockRecorder
}

// MockFileServiceMockRecorder is the mock recorder for MockFileService.
type MockFileServiceMockRecorder struct {
	mock *MockFileService
}

// NewMockFileService creates a new mock instance.
func NewMockFileService(ctrl *gomock.Controller) *MockFileService {
	mock := &MockFileService{ctrl: ctrl}
	mock.recorder = &MockFileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileService) EXPECT() *MockFileServiceMockRecorder {
	return m.recorder
}

// ReadFile mocks base method.
func (m *MockFileService) ReadFile(fileName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", fileName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockFileServiceMockRecorder) ReadFile(fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockFileService)(nil).ReadFile), fileName)
}

// WriteToFile mocks base method.
func (m *MockFileService) WriteToFile(fileName string, reader io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteToFile", fileName, reader)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteToFile indicates an expected call of WriteToFile.
func (mr *MockFileServiceMockRecorder) WriteToFile(fileName, reader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToFile", reflect.TypeOf((*MockFileService)(nil).WriteToFile), fileName, reader)
}

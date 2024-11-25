// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/service.go -destination=internal/service/mock/mock_service.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CreateGuest mocks base method.
func (m *MockAuthorization) CreateGuest(ctx context.Context, phone, pincode string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGuest", ctx, phone, pincode)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGuest indicates an expected call of CreateGuest.
func (mr *MockAuthorizationMockRecorder) CreateGuest(ctx, phone, pincode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGuest", reflect.TypeOf((*MockAuthorization)(nil).CreateGuest), ctx, phone, pincode)
}

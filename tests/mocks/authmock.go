// Code generated by MockGen. DO NOT EDIT.
// Source: authentication.go

// Package mock_authhandle is a generated GoMock package.
package mock_authhandle

import (
	models "github.com/blockseeker999th/URLShortener/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthUser is a mock of AuthUser interface.
type MockAuthUser struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUserMockRecorder
}

// MockAuthUserMockRecorder is the mock recorder for MockAuthUser.
type MockAuthUserMockRecorder struct {
	mock *MockAuthUser
}

// NewMockAuthUser creates a new mock instance.
func NewMockAuthUser(ctrl *gomock.Controller) *MockAuthUser {
	mock := &MockAuthUser{ctrl: ctrl}
	mock.recorder = &MockAuthUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUser) EXPECT() *MockAuthUserMockRecorder {
	return m.recorder
}

// SignInUser mocks base method.
func (m *MockAuthUser) SignInUser(loginData *models.LoginData) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInUser", loginData)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignInUser indicates an expected call of SignInUser.
func (mr *MockAuthUserMockRecorder) SignInUser(loginData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInUser", reflect.TypeOf((*MockAuthUser)(nil).SignInUser), loginData)
}

// SignUpUser mocks base method.
func (m *MockAuthUser) SignUpUser(arg0 *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUpUser", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUpUser indicates an expected call of SignUpUser.
func (mr *MockAuthUserMockRecorder) SignUpUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUpUser", reflect.TypeOf((*MockAuthUser)(nil).SignUpUser), arg0)
}

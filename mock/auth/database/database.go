// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	entity "github.com/and-period/furumane/internal/auth/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockAdmin is a mock of Admin interface.
type MockAdmin struct {
	ctrl     *gomock.Controller
	recorder *MockAdminMockRecorder
}

// MockAdminMockRecorder is the mock recorder for MockAdmin.
type MockAdminMockRecorder struct {
	mock *MockAdmin
}

// NewMockAdmin creates a new mock instance.
func NewMockAdmin(ctrl *gomock.Controller) *MockAdmin {
	mock := &MockAdmin{ctrl: ctrl}
	mock.recorder = &MockAdminMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdmin) EXPECT() *MockAdminMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAdmin) Create(ctx context.Context, admin *entity.Admin, auth func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, admin, auth)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAdminMockRecorder) Create(ctx, admin, auth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAdmin)(nil).Create), ctx, admin, auth)
}

// Delete mocks base method.
func (m *MockAdmin) Delete(ctx context.Context, adminID string, auth func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, adminID, auth)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAdminMockRecorder) Delete(ctx, adminID, auth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAdmin)(nil).Delete), ctx, adminID, auth)
}

// Get mocks base method.
func (m *MockAdmin) Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, adminID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAdminMockRecorder) Get(ctx, adminID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, adminID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAdmin)(nil).Get), varargs...)
}

// GetByCognitoID mocks base method.
func (m *MockAdmin) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, cognitoID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByCognitoID", varargs...)
	ret0, _ := ret[0].(*entity.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCognitoID indicates an expected call of GetByCognitoID.
func (mr *MockAdminMockRecorder) GetByCognitoID(ctx, cognitoID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, cognitoID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCognitoID", reflect.TypeOf((*MockAdmin)(nil).GetByCognitoID), varargs...)
}

// GetByEmail mocks base method.
func (m *MockAdmin) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Admin, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, email}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByEmail", varargs...)
	ret0, _ := ret[0].(*entity.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockAdminMockRecorder) GetByEmail(ctx, email interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, email}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockAdmin)(nil).GetByEmail), varargs...)
}

// UpdateEmail mocks base method.
func (m *MockAdmin) UpdateEmail(ctx context.Context, adminID, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", ctx, adminID, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockAdminMockRecorder) UpdateEmail(ctx, adminID, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockAdmin)(nil).UpdateEmail), ctx, adminID, email)
}

// UpdateVerifiedAt mocks base method.
func (m *MockAdmin) UpdateVerifiedAt(ctx context.Context, adminID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVerifiedAt", ctx, adminID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVerifiedAt indicates an expected call of UpdateVerifiedAt.
func (mr *MockAdminMockRecorder) UpdateVerifiedAt(ctx, adminID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVerifiedAt", reflect.TypeOf((*MockAdmin)(nil).UpdateVerifiedAt), ctx, adminID)
}

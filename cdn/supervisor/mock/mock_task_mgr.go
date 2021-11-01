// Code generated by MockGen. DO NOT EDIT.
// Source: d7y.io/dragonfly/v2/cdn/supervisor (interfaces: SeedTaskManager)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "d7y.io/dragonfly/v2/cdn/types"
	gomock "github.com/golang/mock/gomock"
)

// MockSeedTaskManager is a mock of SeedTaskManager interface.
type MockSeedTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockSeedTaskManagerMockRecorder
}

// MockSeedTaskManagerMockRecorder is the mock recorder for MockSeedTaskManager.
type MockSeedTaskManagerMockRecorder struct {
	mock *MockSeedTaskManager
}

// NewMockSeedTaskManager creates a new mock instance.
func NewMockSeedTaskManager(ctrl *gomock.Controller) *MockSeedTaskManager {
	mock := &MockSeedTaskManager{ctrl: ctrl}
	mock.recorder = &MockSeedTaskManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSeedTaskManager) EXPECT() *MockSeedTaskManagerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSeedTaskManager) Delete(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", arg0)
}

// Delete indicates an expected call of Delete.
func (mr *MockSeedTaskManagerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSeedTaskManager)(nil).Delete), arg0)
}

// Exist mocks base method.
func (m *MockSeedTaskManager) Exist(arg0 string) (*types.SeedTask, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exist", arg0)
	ret0, _ := ret[0].(*types.SeedTask)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Exist indicates an expected call of Exist.
func (mr *MockSeedTaskManagerMockRecorder) Exist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exist", reflect.TypeOf((*MockSeedTaskManager)(nil).Exist), arg0)
}

// Get mocks base method.
func (m *MockSeedTaskManager) Get(arg0 string) (*types.SeedTask, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*types.SeedTask)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSeedTaskManagerMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSeedTaskManager)(nil).Get), arg0)
}

// GetPieces mocks base method.
func (m *MockSeedTaskManager) GetPieces(arg0 context.Context, arg1 string) ([]*types.SeedPiece, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieces", arg0, arg1)
	ret0, _ := ret[0].([]*types.SeedPiece)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPieces indicates an expected call of GetPieces.
func (mr *MockSeedTaskManagerMockRecorder) GetPieces(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieces", reflect.TypeOf((*MockSeedTaskManager)(nil).GetPieces), arg0, arg1)
}

// Register mocks base method.
func (m *MockSeedTaskManager) Register(arg0 context.Context, arg1 *types.SeedTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockSeedTaskManagerMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockSeedTaskManager)(nil).Register), arg0, arg1)
}

// Update mocks base method.
func (m *MockSeedTaskManager) Update(arg0 string, arg1 *types.SeedTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSeedTaskManagerMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSeedTaskManager)(nil).Update), arg0, arg1)
}

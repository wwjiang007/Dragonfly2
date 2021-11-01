// Code generated by MockGen. DO NOT EDIT.
// Source: d7y.io/dragonfly/v2/cdn/supervisor (interfaces: SeedProgressManager)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "d7y.io/dragonfly/v2/cdn/types"
	gomock "github.com/golang/mock/gomock"
)

// MockSeedProgressManager is a mock of SeedProgressManager interface.
type MockSeedProgressManager struct {
	ctrl     *gomock.Controller
	recorder *MockSeedProgressManagerMockRecorder
}

// MockSeedProgressManagerMockRecorder is the mock recorder for MockSeedProgressManager.
type MockSeedProgressManagerMockRecorder struct {
	mock *MockSeedProgressManager
}

// NewMockSeedProgressManager creates a new mock instance.
func NewMockSeedProgressManager(ctrl *gomock.Controller) *MockSeedProgressManager {
	mock := &MockSeedProgressManager{ctrl: ctrl}
	mock.recorder = &MockSeedProgressManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSeedProgressManager) EXPECT() *MockSeedProgressManagerMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockSeedProgressManager) Clear(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear", arg0)
}

// Clear indicates an expected call of Clear.
func (mr *MockSeedProgressManagerMockRecorder) Clear(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockSeedProgressManager)(nil).Clear), arg0)
}

// GetPieces mocks base method.
func (m *MockSeedProgressManager) GetPieces(arg0 context.Context, arg1 string) ([]*types.SeedPiece, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPieces", arg0, arg1)
	ret0, _ := ret[0].([]*types.SeedPiece)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetPieces indicates an expected call of GetPieces.
func (mr *MockSeedProgressManagerMockRecorder) GetPieces(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPieces", reflect.TypeOf((*MockSeedProgressManager)(nil).GetPieces), arg0, arg1)
}

// InitSeedProgress mocks base method.
func (m *MockSeedProgressManager) InitSeedProgress(arg0 context.Context, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InitSeedProgress", arg0, arg1)
}

// InitSeedProgress indicates an expected call of InitSeedProgress.
func (mr *MockSeedProgressManagerMockRecorder) InitSeedProgress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitSeedProgress", reflect.TypeOf((*MockSeedProgressManager)(nil).InitSeedProgress), arg0, arg1)
}

// PublishPiece mocks base method.
func (m *MockSeedProgressManager) PublishPiece(arg0 context.Context, arg1 string, arg2 *types.SeedPiece) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishPiece", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishPiece indicates an expected call of PublishPiece.
func (mr *MockSeedProgressManagerMockRecorder) PublishPiece(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishPiece", reflect.TypeOf((*MockSeedProgressManager)(nil).PublishPiece), arg0, arg1, arg2)
}

// PublishTask mocks base method.
func (m *MockSeedProgressManager) PublishTask(arg0 context.Context, arg1 string, arg2 *types.SeedTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishTask indicates an expected call of PublishTask.
func (mr *MockSeedProgressManagerMockRecorder) PublishTask(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishTask", reflect.TypeOf((*MockSeedProgressManager)(nil).PublishTask), arg0, arg1, arg2)
}

// WatchSeedProgress mocks base method.
func (m *MockSeedProgressManager) WatchSeedProgress(arg0 context.Context, arg1 *types.SeedTask) (<-chan *types.SeedPiece, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchSeedProgress", arg0, arg1)
	ret0, _ := ret[0].(<-chan *types.SeedPiece)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchSeedProgress indicates an expected call of WatchSeedProgress.
func (mr *MockSeedProgressManagerMockRecorder) WatchSeedProgress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchSeedProgress", reflect.TypeOf((*MockSeedProgressManager)(nil).WatchSeedProgress), arg0, arg1)
}

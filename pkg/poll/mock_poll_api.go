// Code generated by MockGen. DO NOT EDIT.
// Source: poll.go

// Package poll is a generated GoMock package.
package poll

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// MockPollActions is a mock of PollActions interface.
type MockPollActions struct {
	ctrl     *gomock.Controller
	recorder *MockPollActionsMockRecorder
}

// MockPollActionsMockRecorder is the mock recorder for MockPollActions.
type MockPollActionsMockRecorder struct {
	mock *MockPollActions
}

// NewMockPollActions creates a new mock instance.
func NewMockPollActions(ctrl *gomock.Controller) *MockPollActions {
	mock := &MockPollActions{ctrl: ctrl}
	mock.recorder = &MockPollActionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPollActions) EXPECT() *MockPollActionsMockRecorder {
	return m.recorder
}

// ForDaemonSet mocks base method.
func (m *MockPollActions) ForDaemonSet(arg0 context.Context, arg1 *unstructured.Unstructured) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForDaemonSet", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForDaemonSet indicates an expected call of ForDaemonSet.
func (mr *MockPollActionsMockRecorder) ForDaemonSet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForDaemonSet", reflect.TypeOf((*MockPollActions)(nil).ForDaemonSet), arg0, arg1)
}

// ForDaemonSetLogs mocks base method.
func (m *MockPollActions) ForDaemonSetLogs(arg0 context.Context, arg1 *unstructured.Unstructured, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForDaemonSetLogs", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForDaemonSetLogs indicates an expected call of ForDaemonSetLogs.
func (mr *MockPollActionsMockRecorder) ForDaemonSetLogs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForDaemonSetLogs", reflect.TypeOf((*MockPollActions)(nil).ForDaemonSetLogs), arg0, arg1, arg2)
}

// ForResource mocks base method.
func (m *MockPollActions) ForResource(arg0 context.Context, arg1 *unstructured.Unstructured) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForResource", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForResource indicates an expected call of ForResource.
func (mr *MockPollActionsMockRecorder) ForResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForResource", reflect.TypeOf((*MockPollActions)(nil).ForResource), arg0, arg1)
}

// ForResourceUnavailability mocks base method.
func (m *MockPollActions) ForResourceUnavailability(arg0 context.Context, arg1 *unstructured.Unstructured) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForResourceUnavailability", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForResourceUnavailability indicates an expected call of ForResourceUnavailability.
func (mr *MockPollActionsMockRecorder) ForResourceUnavailability(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForResourceUnavailability", reflect.TypeOf((*MockPollActions)(nil).ForResourceUnavailability), arg0, arg1)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/anyproto/any-sync-coordinator/eventlog (interfaces: EventLog)
//
// Generated by this command:
//
//	mockgen -destination mock_eventlog/mock_eventlog.go github.com/anyproto/any-sync-coordinator/eventlog EventLog
//

// Package mock_eventlog is a generated GoMock package.
package mock_eventlog

import (
	context "context"
	reflect "reflect"

	eventlog "github.com/anyproto/any-sync-coordinator/eventlog"
	app "github.com/anyproto/any-sync/app"
	gomock "go.uber.org/mock/gomock"
)

// MockEventLog is a mock of EventLog interface.
type MockEventLog struct {
	ctrl     *gomock.Controller
	recorder *MockEventLogMockRecorder
}

// MockEventLogMockRecorder is the mock recorder for MockEventLog.
type MockEventLogMockRecorder struct {
	mock *MockEventLog
}

// NewMockEventLog creates a new mock instance.
func NewMockEventLog(ctrl *gomock.Controller) *MockEventLog {
	mock := &MockEventLog{ctrl: ctrl}
	mock.recorder = &MockEventLogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventLog) EXPECT() *MockEventLogMockRecorder {
	return m.recorder
}

// AddLog mocks base method.
func (m *MockEventLog) AddLog(arg0 context.Context, arg1 eventlog.EventLogEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLog", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLog indicates an expected call of AddLog.
func (mr *MockEventLogMockRecorder) AddLog(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLog", reflect.TypeOf((*MockEventLog)(nil).AddLog), arg0, arg1)
}

// GetAfter mocks base method.
func (m *MockEventLog) GetAfter(arg0 context.Context, arg1, arg2 string, arg3 uint32) ([]eventlog.EventLogEntry, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAfter", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]eventlog.EventLogEntry)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAfter indicates an expected call of GetAfter.
func (mr *MockEventLogMockRecorder) GetAfter(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAfter", reflect.TypeOf((*MockEventLog)(nil).GetAfter), arg0, arg1, arg2, arg3)
}

// Init mocks base method.
func (m *MockEventLog) Init(arg0 *app.App) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockEventLogMockRecorder) Init(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockEventLog)(nil).Init), arg0)
}

// Name mocks base method.
func (m *MockEventLog) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockEventLogMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockEventLog)(nil).Name))
}

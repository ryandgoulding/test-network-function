// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ryandgoulding/goexpect (interfaces: Expecter)

// Package mock_expect is a generated GoMock package.
package mock_interactive

import (
	reflect "reflect"
	regexp "regexp"
	time "time"

	gomock "github.com/golang/mock/gomock"
	goexpect "github.com/ryandgoulding/goexpect"
)

// MockExpecter is a mocks of Expecter interface
type MockExpecter struct {
	ctrl     *gomock.Controller
	recorder *MockExpecterMockRecorder
}

// MockExpecterMockRecorder is the mocks recorder for MockExpecter
type MockExpecterMockRecorder struct {
	mock *MockExpecter
}

// NewMockExpecter creates a new mocks instance
func NewMockExpecter(ctrl *gomock.Controller) *MockExpecter {
	mock := &MockExpecter{ctrl: ctrl}
	mock.recorder = &MockExpecterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExpecter) EXPECT() *MockExpecterMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockExpecter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockExpecterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockExpecter)(nil).Close))
}

// Expect mocks base method
func (m *MockExpecter) Expect(arg0 *regexp.Regexp, arg1 time.Duration) (string, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expect", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Expect indicates an expected call of Expect
func (mr *MockExpecterMockRecorder) Expect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expect", reflect.TypeOf((*MockExpecter)(nil).Expect), arg0, arg1)
}

// ExpectBatch mocks base method
func (m *MockExpecter) ExpectBatch(arg0 []goexpect.Batcher, arg1 time.Duration) ([]goexpect.BatchRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpectBatch", arg0, arg1)
	ret0, _ := ret[0].([]goexpect.BatchRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExpectBatch indicates an expected call of ExpectBatch
func (mr *MockExpecterMockRecorder) ExpectBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpectBatch", reflect.TypeOf((*MockExpecter)(nil).ExpectBatch), arg0, arg1)
}

// ExpectSwitchCase mocks base method
func (m *MockExpecter) ExpectSwitchCase(arg0 []goexpect.Caser, arg1 time.Duration) (string, []string, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpectSwitchCase", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ExpectSwitchCase indicates an expected call of ExpectSwitchCase
func (mr *MockExpecterMockRecorder) ExpectSwitchCase(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpectSwitchCase", reflect.TypeOf((*MockExpecter)(nil).ExpectSwitchCase), arg0, arg1)
}

// Send mocks base method
func (m *MockExpecter) Send(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockExpecterMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockExpecter)(nil).Send), arg0)
}

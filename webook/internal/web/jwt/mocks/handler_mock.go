// Code generated by MockGen. DO NOT EDIT.
// Source: ./webook/internal/web/jwt/types.go
//
// Generated by this command:
//
//	mockgen -source=./webook/internal/web/jwt/types.go -package=jwtmocks -destination=./webook/internal/web/jwt/mocks/handler_mock.go
//

// Package jwtmocks is a generated GoMock package.
package jwtmocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// CheckSession mocks base method.
func (m *MockHandler) CheckSession(ctx *gin.Context, ssid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", ctx, ssid)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockHandlerMockRecorder) CheckSession(ctx, ssid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockHandler)(nil).CheckSession), ctx, ssid)
}

// ClearToken mocks base method.
func (m *MockHandler) ClearToken(ctx *gin.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearToken", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearToken indicates an expected call of ClearToken.
func (mr *MockHandlerMockRecorder) ClearToken(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearToken", reflect.TypeOf((*MockHandler)(nil).ClearToken), ctx)
}

// ExtractToken mocks base method.
func (m *MockHandler) ExtractToken(ctx *gin.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtractToken", ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// ExtractToken indicates an expected call of ExtractToken.
func (mr *MockHandlerMockRecorder) ExtractToken(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractToken", reflect.TypeOf((*MockHandler)(nil).ExtractToken), ctx)
}

// SetJWTToken mocks base method.
func (m *MockHandler) SetJWTToken(ctx *gin.Context, uid int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetJWTToken", ctx, uid)
}

// SetJWTToken indicates an expected call of SetJWTToken.
func (mr *MockHandlerMockRecorder) SetJWTToken(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetJWTToken", reflect.TypeOf((*MockHandler)(nil).SetJWTToken), ctx, uid)
}

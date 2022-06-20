// Code generated by MockGen. DO NOT EDIT.
// Source: food-roulette-api/internal/facade (interfaces: ServiceI)

// Package facade is a generated GoMock package.
package facade

import (
	context "context"
	models "food-roulette-api/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceI is a mock of ServiceI interface.
type MockServiceI struct {
	ctrl     *gomock.Controller
	recorder *MockServiceIMockRecorder
}

// MockServiceIMockRecorder is the mock recorder for MockServiceI.
type MockServiceIMockRecorder struct {
	mock *MockServiceI
}

// NewMockServiceI creates a new mock instance.
func NewMockServiceI(ctrl *gomock.Controller) *MockServiceI {
	mock := &MockServiceI{ctrl: ctrl}
	mock.recorder = &MockServiceIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceI) EXPECT() *MockServiceIMockRecorder {
	return m.recorder
}

// AddCuisine mocks base method.
func (m *MockServiceI) AddCuisine(arg0 context.Context, arg1 models.AddCuisineRequest) models.CuisineResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCuisine", arg0, arg1)
	ret0, _ := ret[0].(models.CuisineResponse)
	return ret0
}

// AddCuisine indicates an expected call of AddCuisine.
func (mr *MockServiceIMockRecorder) AddCuisine(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCuisine", reflect.TypeOf((*MockServiceI)(nil).AddCuisine), arg0, arg1)
}

// AllCuisines mocks base method.
func (m *MockServiceI) AllCuisines(arg0 context.Context) models.AllCuisinesResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllCuisines", arg0)
	ret0, _ := ret[0].(models.AllCuisinesResponse)
	return ret0
}

// AllCuisines indicates an expected call of AllCuisines.
func (mr *MockServiceIMockRecorder) AllCuisines(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllCuisines", reflect.TypeOf((*MockServiceI)(nil).AllCuisines), arg0)
}

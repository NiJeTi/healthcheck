// Code generated by mockery. DO NOT EDIT.

package healthcheck

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockProbeFunc is an autogenerated mock type for the ProbeFunc type
type MockProbeFunc struct {
	mock.Mock
}

type MockProbeFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProbeFunc) EXPECT() *MockProbeFunc_Expecter {
	return &MockProbeFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx
func (_m *MockProbeFunc) Execute(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProbeFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockProbeFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockProbeFunc_Expecter) Execute(ctx interface{}) *MockProbeFunc_Execute_Call {
	return &MockProbeFunc_Execute_Call{Call: _e.mock.On("Execute", ctx)}
}

func (_c *MockProbeFunc_Execute_Call) Run(run func(ctx context.Context)) *MockProbeFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockProbeFunc_Execute_Call) Return(_a0 error) *MockProbeFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProbeFunc_Execute_Call) RunAndReturn(run func(context.Context) error) *MockProbeFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProbeFunc creates a new instance of MockProbeFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProbeFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProbeFunc {
	mock := &MockProbeFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

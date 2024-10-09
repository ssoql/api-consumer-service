// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockRetryer is an autogenerated mock type for the Retryer type
type MockRetryer struct {
	mock.Mock
}

type MockRetryer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRetryer) EXPECT() *MockRetryer_Expecter {
	return &MockRetryer_Expecter{mock: &_m.Mock}
}

// Retry provides a mock function with given fields: operation
func (_m *MockRetryer) Retry(operation func() error) error {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for Retry")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(func() error) error); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRetryer_Retry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Retry'
type MockRetryer_Retry_Call struct {
	*mock.Call
}

// Retry is a helper method to define mock.On call
//   - operation func() error
func (_e *MockRetryer_Expecter) Retry(operation interface{}) *MockRetryer_Retry_Call {
	return &MockRetryer_Retry_Call{Call: _e.mock.On("Retry", operation)}
}

func (_c *MockRetryer_Retry_Call) Run(run func(operation func() error)) *MockRetryer_Retry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func() error))
	})
	return _c
}

func (_c *MockRetryer_Retry_Call) Return(_a0 error) *MockRetryer_Retry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRetryer_Retry_Call) RunAndReturn(run func(func() error) error) *MockRetryer_Retry_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRetryer creates a new instance of MockRetryer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRetryer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRetryer {
	mock := &MockRetryer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

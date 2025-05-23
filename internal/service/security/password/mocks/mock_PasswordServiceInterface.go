// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package password

import (
	mock "github.com/stretchr/testify/mock"
)

// NewMockPasswordServiceInterface creates a new instance of MockPasswordServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPasswordServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPasswordServiceInterface {
	mock := &MockPasswordServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockPasswordServiceInterface is an autogenerated mock type for the PasswordServiceInterface type
type MockPasswordServiceInterface struct {
	mock.Mock
}

type MockPasswordServiceInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPasswordServiceInterface) EXPECT() *MockPasswordServiceInterface_Expecter {
	return &MockPasswordServiceInterface_Expecter{mock: &_m.Mock}
}

// ComparePassword provides a mock function for the type MockPasswordServiceInterface
func (_mock *MockPasswordServiceInterface) ComparePassword(hashedPassword string, password string) error {
	ret := _mock.Called(hashedPassword, password)

	if len(ret) == 0 {
		panic("no return value specified for ComparePassword")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = returnFunc(hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockPasswordServiceInterface_ComparePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ComparePassword'
type MockPasswordServiceInterface_ComparePassword_Call struct {
	*mock.Call
}

// ComparePassword is a helper method to define mock.On call
//   - hashedPassword
//   - password
func (_e *MockPasswordServiceInterface_Expecter) ComparePassword(hashedPassword interface{}, password interface{}) *MockPasswordServiceInterface_ComparePassword_Call {
	return &MockPasswordServiceInterface_ComparePassword_Call{Call: _e.mock.On("ComparePassword", hashedPassword, password)}
}

func (_c *MockPasswordServiceInterface_ComparePassword_Call) Run(run func(hashedPassword string, password string)) *MockPasswordServiceInterface_ComparePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockPasswordServiceInterface_ComparePassword_Call) Return(err error) *MockPasswordServiceInterface_ComparePassword_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockPasswordServiceInterface_ComparePassword_Call) RunAndReturn(run func(hashedPassword string, password string) error) *MockPasswordServiceInterface_ComparePassword_Call {
	_c.Call.Return(run)
	return _c
}

// HashAndSaltPassword provides a mock function for the type MockPasswordServiceInterface
func (_mock *MockPasswordServiceInterface) HashAndSaltPassword(password string) (string, error) {
	ret := _mock.Called(password)

	if len(ret) == 0 {
		panic("no return value specified for HashAndSaltPassword")
	}

	var r0 string
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (string, error)); ok {
		return returnFunc(password)
	}
	if returnFunc, ok := ret.Get(0).(func(string) string); ok {
		r0 = returnFunc(password)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(password)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockPasswordServiceInterface_HashAndSaltPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HashAndSaltPassword'
type MockPasswordServiceInterface_HashAndSaltPassword_Call struct {
	*mock.Call
}

// HashAndSaltPassword is a helper method to define mock.On call
//   - password
func (_e *MockPasswordServiceInterface_Expecter) HashAndSaltPassword(password interface{}) *MockPasswordServiceInterface_HashAndSaltPassword_Call {
	return &MockPasswordServiceInterface_HashAndSaltPassword_Call{Call: _e.mock.On("HashAndSaltPassword", password)}
}

func (_c *MockPasswordServiceInterface_HashAndSaltPassword_Call) Run(run func(password string)) *MockPasswordServiceInterface_HashAndSaltPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPasswordServiceInterface_HashAndSaltPassword_Call) Return(s string, err error) *MockPasswordServiceInterface_HashAndSaltPassword_Call {
	_c.Call.Return(s, err)
	return _c
}

func (_c *MockPasswordServiceInterface_HashAndSaltPassword_Call) RunAndReturn(run func(password string) (string, error)) *MockPasswordServiceInterface_HashAndSaltPassword_Call {
	_c.Call.Return(run)
	return _c
}

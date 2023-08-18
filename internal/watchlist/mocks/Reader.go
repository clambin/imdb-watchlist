// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	imdb "github.com/clambin/imdb-watchlist/pkg/imdb"
	mock "github.com/stretchr/testify/mock"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

type Reader_Expecter struct {
	mock *mock.Mock
}

func (_m *Reader) EXPECT() *Reader_Expecter {
	return &Reader_Expecter{mock: &_m.Mock}
}

// ReadByTypes provides a mock function with given fields: validTypes
func (_m *Reader) ReadByTypes(validTypes ...imdb.EntryType) ([]imdb.Entry, error) {
	_va := make([]interface{}, len(validTypes))
	for _i := range validTypes {
		_va[_i] = validTypes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []imdb.Entry
	var r1 error
	if rf, ok := ret.Get(0).(func(...imdb.EntryType) ([]imdb.Entry, error)); ok {
		return rf(validTypes...)
	}
	if rf, ok := ret.Get(0).(func(...imdb.EntryType) []imdb.Entry); ok {
		r0 = rf(validTypes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]imdb.Entry)
		}
	}

	if rf, ok := ret.Get(1).(func(...imdb.EntryType) error); ok {
		r1 = rf(validTypes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reader_ReadByTypes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadByTypes'
type Reader_ReadByTypes_Call struct {
	*mock.Call
}

// ReadByTypes is a helper method to define mock.On call
//   - validTypes ...imdb.EntryType
func (_e *Reader_Expecter) ReadByTypes(validTypes ...interface{}) *Reader_ReadByTypes_Call {
	return &Reader_ReadByTypes_Call{Call: _e.mock.On("ReadByTypes",
		append([]interface{}{}, validTypes...)...)}
}

func (_c *Reader_ReadByTypes_Call) Run(run func(validTypes ...imdb.EntryType)) *Reader_ReadByTypes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]imdb.EntryType, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(imdb.EntryType)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Reader_ReadByTypes_Call) Return(entries []imdb.Entry, err error) *Reader_ReadByTypes_Call {
	_c.Call.Return(entries, err)
	return _c
}

func (_c *Reader_ReadByTypes_Call) RunAndReturn(run func(...imdb.EntryType) ([]imdb.Entry, error)) *Reader_ReadByTypes_Call {
	_c.Call.Return(run)
	return _c
}

// NewReader creates a new instance of Reader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reader {
	mock := &Reader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
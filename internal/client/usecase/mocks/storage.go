// Code generated by mockery v2.26.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/Spear5030/yagophkeeper/internal/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// storage is an autogenerated mock type for the storage type
type storage struct {
	mock.Mock
}

// AddBinaryData provides a mock function with given fields: _a0
func (_m *storage) AddBinaryData(_a0 domain.BinaryData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.BinaryData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddCardData provides a mock function with given fields: _a0
func (_m *storage) AddCardData(_a0 domain.CardData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.CardData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddLoginPassword provides a mock function with given fields: _a0
func (_m *storage) AddLoginPassword(_a0 domain.LoginPassword) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.LoginPassword) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddTextData provides a mock function with given fields: _a0
func (_m *storage) AddTextData(_a0 domain.TextData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.TextData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBinaryData provides a mock function with given fields:
func (_m *storage) GetBinaryData() []domain.BinaryData {
	ret := _m.Called()

	var r0 []domain.BinaryData
	if rf, ok := ret.Get(0).(func() []domain.BinaryData); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BinaryData)
		}
	}

	return r0
}

// GetCardsData provides a mock function with given fields:
func (_m *storage) GetCardsData() []domain.CardData {
	ret := _m.Called()

	var r0 []domain.CardData
	if rf, ok := ret.Get(0).(func() []domain.CardData); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.CardData)
		}
	}

	return r0
}

// GetData provides a mock function with given fields:
func (_m *storage) GetData() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLocalSyncTime provides a mock function with given fields:
func (_m *storage) GetLocalSyncTime() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetLogins provides a mock function with given fields:
func (_m *storage) GetLogins() []domain.LoginPassword {
	ret := _m.Called()

	var r0 []domain.LoginPassword
	if rf, ok := ret.Get(0).(func() []domain.LoginPassword); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.LoginPassword)
		}
	}

	return r0
}

// GetTextData provides a mock function with given fields:
func (_m *storage) GetTextData() []domain.TextData {
	ret := _m.Called()

	var r0 []domain.TextData
	if rf, ok := ret.Get(0).(func() []domain.TextData); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TextData)
		}
	}

	return r0
}

// SaveUserData provides a mock function with given fields: user, token
func (_m *storage) SaveUserData(user domain.User, token string) error {
	ret := _m.Called(user, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.User, string) error); ok {
		r0 = rf(user, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetData provides a mock function with given fields: data
func (_m *storage) SetData(data []byte) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTime provides a mock function with given fields:
func (_m *storage) UpdateTime() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewStorage interface {
	mock.TestingT
	Cleanup(func())
}

// newStorage creates a new instance of storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newStorage(t mockConstructorTestingTnewStorage) *storage {
	mock := &storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

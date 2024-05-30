// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// CacheRepository is an autogenerated mock type for the CacheRepository type
type CacheRepository struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *CacheRepository) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, key
func (_m *CacheRepository) Delete(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByPrefix provides a mock function with given fields: ctx, prefix
func (_m *CacheRepository) DeleteByPrefix(ctx context.Context, prefix string) error {
	ret := _m.Called(ctx, prefix)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, prefix)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, key
func (_m *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	ret := _m.Called(ctx, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, key
func (_m *CacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	ret := _m.Called(ctx, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, key, value, ttl
func (_m *CacheRepository) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	ret := _m.Called(ctx, key, value, ttl)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, time.Duration) error); ok {
		r0 = rf(ctx, key, value, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCacheRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCacheRepository creates a new instance of CacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCacheRepository(t mockConstructorTestingTNewCacheRepository) *CacheRepository {
	mock := &CacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	dao "github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// CollectionRepository is an autogenerated mock type for the CollectionRepository type
type CollectionRepository struct {
	mock.Mock
}

// CreateCollection provides a mock function with given fields: ctx, userID, collection
func (_m *CollectionRepository) CreateCollection(ctx context.Context, userID uuid.UUID, collection *dao.CollectionDAO) (*dao.CollectionDAO, error) {
	ret := _m.Called(ctx, userID, collection)

	var r0 *dao.CollectionDAO
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *dao.CollectionDAO) *dao.CollectionDAO); ok {
		r0 = rf(ctx, userID, collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.CollectionDAO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, *dao.CollectionDAO) error); ok {
		r1 = rf(ctx, userID, collection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCollection provides a mock function with given fields: ctx, id
func (_m *CollectionRepository) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCollectionByID provides a mock function with given fields: ctx, id
func (_m *CollectionRepository) GetCollectionByID(ctx context.Context, id uuid.UUID) (*dao.CollectionDAO, error) {
	ret := _m.Called(ctx, id)

	var r0 *dao.CollectionDAO
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *dao.CollectionDAO); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.CollectionDAO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserPartOfCollection provides a mock function with given fields: ctx, userID, collectionID
func (_m *CollectionRepository) IsUserPartOfCollection(ctx context.Context, userID uuid.UUID, collectionID uuid.UUID) (bool, error) {
	ret := _m.Called(ctx, userID, collectionID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) bool); ok {
		r0 = rf(ctx, userID, collectionID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userID, collectionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListCollectionsByUserID provides a mock function with given fields: ctx, userID, skip, limit
func (_m *CollectionRepository) ListCollectionsByUserID(ctx context.Context, userID uuid.UUID, skip uint64, limit uint64) ([]dao.CollectionDAO, error) {
	ret := _m.Called(ctx, userID, skip, limit)

	var r0 []dao.CollectionDAO
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uint64, uint64) []dao.CollectionDAO); ok {
		r0 = rf(ctx, userID, skip, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dao.CollectionDAO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uint64, uint64) error); ok {
		r1 = rf(ctx, userID, skip, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCollection provides a mock function with given fields: ctx, collection
func (_m *CollectionRepository) UpdateCollection(ctx context.Context, collection *dao.CollectionDAO) (*dao.CollectionDAO, error) {
	ret := _m.Called(ctx, collection)

	var r0 *dao.CollectionDAO
	if rf, ok := ret.Get(0).(func(context.Context, *dao.CollectionDAO) *dao.CollectionDAO); ok {
		r0 = rf(ctx, collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.CollectionDAO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dao.CollectionDAO) error); ok {
		r1 = rf(ctx, collection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCollectionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCollectionRepository creates a new instance of CollectionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCollectionRepository(t mockConstructorTestingTNewCollectionRepository) *CollectionRepository {
	mock := &CollectionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

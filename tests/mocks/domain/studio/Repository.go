// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/rl404/akatsuki/internal/domain/studio/entity"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// BatchUpdate provides a mock function with given fields: ctx, data
func (_m *Repository) BatchUpdate(ctx context.Context, data []entity.Studio) (int, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for BatchUpdate")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []entity.Studio) (int, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []entity.Studio) int); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []entity.Studio) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, data
func (_m *Repository) Get(ctx context.Context, data entity.GetRequest) ([]*entity.Studio, int, int, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []*entity.Studio
	var r1 int
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.GetRequest) ([]*entity.Studio, int, int, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.GetRequest) []*entity.Studio); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Studio)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.GetRequest) int); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, entity.GetRequest) int); ok {
		r2 = rf(ctx, data)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(context.Context, entity.GetRequest) error); ok {
		r3 = rf(ctx, data)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Repository) GetByID(ctx context.Context, id int64) (*entity.Studio, int, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *entity.Studio
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entity.Studio, int, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Studio); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Studio)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) int); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int64) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByIDs provides a mock function with given fields: ctx, ids
func (_m *Repository) GetByIDs(ctx context.Context, ids []int64) ([]*entity.Studio, int, error) {
	ret := _m.Called(ctx, ids)

	if len(ret) == 0 {
		panic("no return value specified for GetByIDs")
	}

	var r0 []*entity.Studio
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, []int64) ([]*entity.Studio, int, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []int64) []*entity.Studio); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Studio)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []int64) int); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, []int64) error); ok {
		r2 = rf(ctx, ids)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetHistories provides a mock function with given fields: ctx, data
func (_m *Repository) GetHistories(ctx context.Context, data entity.GetHistoriesRequest) ([]entity.History, int, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for GetHistories")
	}

	var r0 []entity.History
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.GetHistoriesRequest) ([]entity.History, int, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.GetHistoriesRequest) []entity.History); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.History)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.GetHistoriesRequest) int); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, entity.GetHistoriesRequest) error); ok {
		r2 = rf(ctx, data)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

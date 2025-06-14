// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	models "github.com/susek555/BD2/car-dealer-api/internal/models"
)

// BidRepositoryInterface is an autogenerated mock type for the BidRepositoryInterface type
type BidRepositoryInterface struct {
	mock.Mock
}

type BidRepositoryInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *BidRepositoryInterface) EXPECT() *BidRepositoryInterface_Expecter {
	return &BidRepositoryInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0
func (_m *BidRepositoryInterface) Create(_a0 *models.Bid) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Bid) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BidRepositoryInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type BidRepositoryInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 *models.Bid
func (_e *BidRepositoryInterface_Expecter) Create(_a0 interface{}) *BidRepositoryInterface_Create_Call {
	return &BidRepositoryInterface_Create_Call{Call: _e.mock.On("Create", _a0)}
}

func (_c *BidRepositoryInterface_Create_Call) Run(run func(_a0 *models.Bid)) *BidRepositoryInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Bid))
	})
	return _c
}

func (_c *BidRepositoryInterface_Create_Call) Return(_a0 error) *BidRepositoryInterface_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BidRepositoryInterface_Create_Call) RunAndReturn(run func(*models.Bid) error) *BidRepositoryInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with no fields
func (_m *BidRepositoryInterface) GetAll() ([]models.Bid, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []models.Bid
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Bid, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Bid); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Bid)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BidRepositoryInterface_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type BidRepositoryInterface_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *BidRepositoryInterface_Expecter) GetAll() *BidRepositoryInterface_GetAll_Call {
	return &BidRepositoryInterface_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *BidRepositoryInterface_GetAll_Call) Run(run func()) *BidRepositoryInterface_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *BidRepositoryInterface_GetAll_Call) Return(_a0 []models.Bid, _a1 error) *BidRepositoryInterface_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BidRepositoryInterface_GetAll_Call) RunAndReturn(run func() ([]models.Bid, error)) *BidRepositoryInterface_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByAuctionID provides a mock function with given fields: auctionID
func (_m *BidRepositoryInterface) GetByAuctionID(auctionID uint) ([]models.Bid, error) {
	ret := _m.Called(auctionID)

	if len(ret) == 0 {
		panic("no return value specified for GetByAuctionID")
	}

	var r0 []models.Bid
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) ([]models.Bid, error)); ok {
		return rf(auctionID)
	}
	if rf, ok := ret.Get(0).(func(uint) []models.Bid); ok {
		r0 = rf(auctionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Bid)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(auctionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BidRepositoryInterface_GetByAuctionID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByAuctionID'
type BidRepositoryInterface_GetByAuctionID_Call struct {
	*mock.Call
}

// GetByAuctionID is a helper method to define mock.On call
//   - auctionID uint
func (_e *BidRepositoryInterface_Expecter) GetByAuctionID(auctionID interface{}) *BidRepositoryInterface_GetByAuctionID_Call {
	return &BidRepositoryInterface_GetByAuctionID_Call{Call: _e.mock.On("GetByAuctionID", auctionID)}
}

func (_c *BidRepositoryInterface_GetByAuctionID_Call) Run(run func(auctionID uint)) *BidRepositoryInterface_GetByAuctionID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *BidRepositoryInterface_GetByAuctionID_Call) Return(_a0 []models.Bid, _a1 error) *BidRepositoryInterface_GetByAuctionID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BidRepositoryInterface_GetByAuctionID_Call) RunAndReturn(run func(uint) ([]models.Bid, error)) *BidRepositoryInterface_GetByAuctionID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByBidderID provides a mock function with given fields: bidderID
func (_m *BidRepositoryInterface) GetByBidderID(bidderID uint) ([]models.Bid, error) {
	ret := _m.Called(bidderID)

	if len(ret) == 0 {
		panic("no return value specified for GetByBidderID")
	}

	var r0 []models.Bid
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) ([]models.Bid, error)); ok {
		return rf(bidderID)
	}
	if rf, ok := ret.Get(0).(func(uint) []models.Bid); ok {
		r0 = rf(bidderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Bid)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(bidderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BidRepositoryInterface_GetByBidderID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByBidderID'
type BidRepositoryInterface_GetByBidderID_Call struct {
	*mock.Call
}

// GetByBidderID is a helper method to define mock.On call
//   - bidderID uint
func (_e *BidRepositoryInterface_Expecter) GetByBidderID(bidderID interface{}) *BidRepositoryInterface_GetByBidderID_Call {
	return &BidRepositoryInterface_GetByBidderID_Call{Call: _e.mock.On("GetByBidderID", bidderID)}
}

func (_c *BidRepositoryInterface_GetByBidderID_Call) Run(run func(bidderID uint)) *BidRepositoryInterface_GetByBidderID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *BidRepositoryInterface_GetByBidderID_Call) Return(_a0 []models.Bid, _a1 error) *BidRepositoryInterface_GetByBidderID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BidRepositoryInterface_GetByBidderID_Call) RunAndReturn(run func(uint) ([]models.Bid, error)) *BidRepositoryInterface_GetByBidderID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: id
func (_m *BidRepositoryInterface) GetByID(id uint) (*models.Bid, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *models.Bid
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*models.Bid, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *models.Bid); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Bid)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BidRepositoryInterface_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type BidRepositoryInterface_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - id uint
func (_e *BidRepositoryInterface_Expecter) GetByID(id interface{}) *BidRepositoryInterface_GetByID_Call {
	return &BidRepositoryInterface_GetByID_Call{Call: _e.mock.On("GetByID", id)}
}

func (_c *BidRepositoryInterface_GetByID_Call) Run(run func(id uint)) *BidRepositoryInterface_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *BidRepositoryInterface_GetByID_Call) Return(_a0 *models.Bid, _a1 error) *BidRepositoryInterface_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BidRepositoryInterface_GetByID_Call) RunAndReturn(run func(uint) (*models.Bid, error)) *BidRepositoryInterface_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetHighestBid provides a mock function with given fields: auctionID
func (_m *BidRepositoryInterface) GetHighestBid(auctionID uint) (*models.Bid, error) {
	ret := _m.Called(auctionID)

	if len(ret) == 0 {
		panic("no return value specified for GetHighestBid")
	}

	var r0 *models.Bid
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*models.Bid, error)); ok {
		return rf(auctionID)
	}
	if rf, ok := ret.Get(0).(func(uint) *models.Bid); ok {
		r0 = rf(auctionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Bid)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(auctionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BidRepositoryInterface_GetHighestBid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHighestBid'
type BidRepositoryInterface_GetHighestBid_Call struct {
	*mock.Call
}

// GetHighestBid is a helper method to define mock.On call
//   - auctionID uint
func (_e *BidRepositoryInterface_Expecter) GetHighestBid(auctionID interface{}) *BidRepositoryInterface_GetHighestBid_Call {
	return &BidRepositoryInterface_GetHighestBid_Call{Call: _e.mock.On("GetHighestBid", auctionID)}
}

func (_c *BidRepositoryInterface_GetHighestBid_Call) Run(run func(auctionID uint)) *BidRepositoryInterface_GetHighestBid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *BidRepositoryInterface_GetHighestBid_Call) Return(_a0 *models.Bid, _a1 error) *BidRepositoryInterface_GetHighestBid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BidRepositoryInterface_GetHighestBid_Call) RunAndReturn(run func(uint) (*models.Bid, error)) *BidRepositoryInterface_GetHighestBid_Call {
	_c.Call.Return(run)
	return _c
}

// GetHighestBidByUserID provides a mock function with given fields: auctionID, userID
func (_m *BidRepositoryInterface) GetHighestBidByUserID(auctionID uint, userID uint) (*models.Bid, error) {
	ret := _m.Called(auctionID, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetHighestBidByUserID")
	}

	var r0 *models.Bid
	var r1 error
	if rf, ok := ret.Get(0).(func(uint, uint) (*models.Bid, error)); ok {
		return rf(auctionID, userID)
	}
	if rf, ok := ret.Get(0).(func(uint, uint) *models.Bid); ok {
		r0 = rf(auctionID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Bid)
		}
	}

	if rf, ok := ret.Get(1).(func(uint, uint) error); ok {
		r1 = rf(auctionID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BidRepositoryInterface_GetHighestBidByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHighestBidByUserID'
type BidRepositoryInterface_GetHighestBidByUserID_Call struct {
	*mock.Call
}

// GetHighestBidByUserID is a helper method to define mock.On call
//   - auctionID uint
//   - userID uint
func (_e *BidRepositoryInterface_Expecter) GetHighestBidByUserID(auctionID interface{}, userID interface{}) *BidRepositoryInterface_GetHighestBidByUserID_Call {
	return &BidRepositoryInterface_GetHighestBidByUserID_Call{Call: _e.mock.On("GetHighestBidByUserID", auctionID, userID)}
}

func (_c *BidRepositoryInterface_GetHighestBidByUserID_Call) Run(run func(auctionID uint, userID uint)) *BidRepositoryInterface_GetHighestBidByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint), args[1].(uint))
	})
	return _c
}

func (_c *BidRepositoryInterface_GetHighestBidByUserID_Call) Return(_a0 *models.Bid, _a1 error) *BidRepositoryInterface_GetHighestBidByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BidRepositoryInterface_GetHighestBidByUserID_Call) RunAndReturn(run func(uint, uint) (*models.Bid, error)) *BidRepositoryInterface_GetHighestBidByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// NewBidRepositoryInterface creates a new instance of BidRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBidRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *BidRepositoryInterface {
	mock := &BidRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	enums "github.com/susek555/BD2/car-dealer-api/internal/enums"

	models "github.com/susek555/BD2/car-dealer-api/internal/models"

	pagination "github.com/susek555/BD2/car-dealer-api/pkg/pagination"

	sale_offer "github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"

	views "github.com/susek555/BD2/car-dealer-api/internal/views"
)

// SaleOfferRepositoryInterface is an autogenerated mock type for the SaleOfferRepositoryInterface type
type SaleOfferRepositoryInterface struct {
	mock.Mock
}

type SaleOfferRepositoryInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *SaleOfferRepositoryInterface) EXPECT() *SaleOfferRepositoryInterface_Expecter {
	return &SaleOfferRepositoryInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: offer
func (_m *SaleOfferRepositoryInterface) Create(offer *models.SaleOffer) error {
	ret := _m.Called(offer)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SaleOffer) error); ok {
		r0 = rf(offer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaleOfferRepositoryInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type SaleOfferRepositoryInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - offer *models.SaleOffer
func (_e *SaleOfferRepositoryInterface_Expecter) Create(offer interface{}) *SaleOfferRepositoryInterface_Create_Call {
	return &SaleOfferRepositoryInterface_Create_Call{Call: _e.mock.On("Create", offer)}
}

func (_c *SaleOfferRepositoryInterface_Create_Call) Run(run func(offer *models.SaleOffer)) *SaleOfferRepositoryInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.SaleOffer))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_Create_Call) Return(_a0 error) *SaleOfferRepositoryInterface_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SaleOfferRepositoryInterface_Create_Call) RunAndReturn(run func(*models.SaleOffer) error) *SaleOfferRepositoryInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: id
func (_m *SaleOfferRepositoryInterface) Delete(id uint) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaleOfferRepositoryInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type SaleOfferRepositoryInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id uint
func (_e *SaleOfferRepositoryInterface_Expecter) Delete(id interface{}) *SaleOfferRepositoryInterface_Delete_Call {
	return &SaleOfferRepositoryInterface_Delete_Call{Call: _e.mock.On("Delete", id)}
}

func (_c *SaleOfferRepositoryInterface_Delete_Call) Run(run func(id uint)) *SaleOfferRepositoryInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_Delete_Call) Return(_a0 error) *SaleOfferRepositoryInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SaleOfferRepositoryInterface_Delete_Call) RunAndReturn(run func(uint) error) *SaleOfferRepositoryInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllActiveAuctions provides a mock function with no fields
func (_m *SaleOfferRepositoryInterface) GetAllActiveAuctions() ([]views.SaleOfferView, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAllActiveAuctions")
	}

	var r0 []views.SaleOfferView
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]views.SaleOfferView, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []views.SaleOfferView); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]views.SaleOfferView)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaleOfferRepositoryInterface_GetAllActiveAuctions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllActiveAuctions'
type SaleOfferRepositoryInterface_GetAllActiveAuctions_Call struct {
	*mock.Call
}

// GetAllActiveAuctions is a helper method to define mock.On call
func (_e *SaleOfferRepositoryInterface_Expecter) GetAllActiveAuctions() *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call {
	return &SaleOfferRepositoryInterface_GetAllActiveAuctions_Call{Call: _e.mock.On("GetAllActiveAuctions")}
}

func (_c *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call) Run(run func()) *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call) Return(_a0 []views.SaleOfferView, _a1 error) *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call) RunAndReturn(run func() ([]views.SaleOfferView, error)) *SaleOfferRepositoryInterface_GetAllActiveAuctions_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: id
func (_m *SaleOfferRepositoryInterface) GetByID(id uint) (*models.SaleOffer, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *models.SaleOffer
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*models.SaleOffer, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *models.SaleOffer); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SaleOffer)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaleOfferRepositoryInterface_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type SaleOfferRepositoryInterface_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - id uint
func (_e *SaleOfferRepositoryInterface_Expecter) GetByID(id interface{}) *SaleOfferRepositoryInterface_GetByID_Call {
	return &SaleOfferRepositoryInterface_GetByID_Call{Call: _e.mock.On("GetByID", id)}
}

func (_c *SaleOfferRepositoryInterface_GetByID_Call) Run(run func(id uint)) *SaleOfferRepositoryInterface_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetByID_Call) Return(_a0 *models.SaleOffer, _a1 error) *SaleOfferRepositoryInterface_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetByID_Call) RunAndReturn(run func(uint) (*models.SaleOffer, error)) *SaleOfferRepositoryInterface_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetFiltered provides a mock function with given fields: filter, pagRequest
func (_m *SaleOfferRepositoryInterface) GetFiltered(filter sale_offer.OfferFilterInterface, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
	ret := _m.Called(filter, pagRequest)

	if len(ret) == 0 {
		panic("no return value specified for GetFiltered")
	}

	var r0 []views.SaleOfferView
	var r1 *pagination.PaginationResponse
	var r2 error
	if rf, ok := ret.Get(0).(func(sale_offer.OfferFilterInterface, *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)); ok {
		return rf(filter, pagRequest)
	}
	if rf, ok := ret.Get(0).(func(sale_offer.OfferFilterInterface, *pagination.PaginationRequest) []views.SaleOfferView); ok {
		r0 = rf(filter, pagRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]views.SaleOfferView)
		}
	}

	if rf, ok := ret.Get(1).(func(sale_offer.OfferFilterInterface, *pagination.PaginationRequest) *pagination.PaginationResponse); ok {
		r1 = rf(filter, pagRequest)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*pagination.PaginationResponse)
		}
	}

	if rf, ok := ret.Get(2).(func(sale_offer.OfferFilterInterface, *pagination.PaginationRequest) error); ok {
		r2 = rf(filter, pagRequest)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SaleOfferRepositoryInterface_GetFiltered_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFiltered'
type SaleOfferRepositoryInterface_GetFiltered_Call struct {
	*mock.Call
}

// GetFiltered is a helper method to define mock.On call
//   - filter sale_offer.OfferFilterIntreface
//   - pagRequest *pagination.PaginationRequest
func (_e *SaleOfferRepositoryInterface_Expecter) GetFiltered(filter interface{}, pagRequest interface{}) *SaleOfferRepositoryInterface_GetFiltered_Call {
	return &SaleOfferRepositoryInterface_GetFiltered_Call{Call: _e.mock.On("GetFiltered", filter, pagRequest)}
}

func (_c *SaleOfferRepositoryInterface_GetFiltered_Call) Run(run func(filter sale_offer.OfferFilterInterface, pagRequest *pagination.PaginationRequest)) *SaleOfferRepositoryInterface_GetFiltered_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(sale_offer.OfferFilterInterface), args[1].(*pagination.PaginationRequest))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetFiltered_Call) Return(_a0 []views.SaleOfferView, _a1 *pagination.PaginationResponse, _a2 error) *SaleOfferRepositoryInterface_GetFiltered_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetFiltered_Call) RunAndReturn(run func(sale_offer.OfferFilterInterface, *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)) *SaleOfferRepositoryInterface_GetFiltered_Call {
	_c.Call.Return(run)
	return _c
}

// GetViewByID provides a mock function with given fields: id
func (_m *SaleOfferRepositoryInterface) GetViewByID(id uint) (*views.SaleOfferView, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetViewByID")
	}

	var r0 *views.SaleOfferView
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*views.SaleOfferView, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *views.SaleOfferView); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*views.SaleOfferView)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaleOfferRepositoryInterface_GetViewByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetViewByID'
type SaleOfferRepositoryInterface_GetViewByID_Call struct {
	*mock.Call
}

// GetViewByID is a helper method to define mock.On call
//   - id uint
func (_e *SaleOfferRepositoryInterface_Expecter) GetViewByID(id interface{}) *SaleOfferRepositoryInterface_GetViewByID_Call {
	return &SaleOfferRepositoryInterface_GetViewByID_Call{Call: _e.mock.On("GetViewByID", id)}
}

func (_c *SaleOfferRepositoryInterface_GetViewByID_Call) Run(run func(id uint)) *SaleOfferRepositoryInterface_GetViewByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetViewByID_Call) Return(_a0 *views.SaleOfferView, _a1 error) *SaleOfferRepositoryInterface_GetViewByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SaleOfferRepositoryInterface_GetViewByID_Call) RunAndReturn(run func(uint) (*views.SaleOfferView, error)) *SaleOfferRepositoryInterface_GetViewByID_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: offer
func (_m *SaleOfferRepositoryInterface) Update(offer *models.SaleOffer) error {
	ret := _m.Called(offer)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SaleOffer) error); ok {
		r0 = rf(offer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaleOfferRepositoryInterface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type SaleOfferRepositoryInterface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - offer *models.SaleOffer
func (_e *SaleOfferRepositoryInterface_Expecter) Update(offer interface{}) *SaleOfferRepositoryInterface_Update_Call {
	return &SaleOfferRepositoryInterface_Update_Call{Call: _e.mock.On("Update", offer)}
}

func (_c *SaleOfferRepositoryInterface_Update_Call) Run(run func(offer *models.SaleOffer)) *SaleOfferRepositoryInterface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.SaleOffer))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_Update_Call) Return(_a0 error) *SaleOfferRepositoryInterface_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SaleOfferRepositoryInterface_Update_Call) RunAndReturn(run func(*models.SaleOffer) error) *SaleOfferRepositoryInterface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: offer, status
func (_m *SaleOfferRepositoryInterface) UpdateStatus(offer *models.SaleOffer, status enums.Status) error {
	ret := _m.Called(offer, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SaleOffer, enums.Status) error); ok {
		r0 = rf(offer, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaleOfferRepositoryInterface_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type SaleOfferRepositoryInterface_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - offer *models.SaleOffer
//   - status enums.Status
func (_e *SaleOfferRepositoryInterface_Expecter) UpdateStatus(offer interface{}, status interface{}) *SaleOfferRepositoryInterface_UpdateStatus_Call {
	return &SaleOfferRepositoryInterface_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", offer, status)}
}

func (_c *SaleOfferRepositoryInterface_UpdateStatus_Call) Run(run func(offer *models.SaleOffer, status enums.Status)) *SaleOfferRepositoryInterface_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.SaleOffer), args[1].(enums.Status))
	})
	return _c
}

func (_c *SaleOfferRepositoryInterface_UpdateStatus_Call) Return(_a0 error) *SaleOfferRepositoryInterface_UpdateStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SaleOfferRepositoryInterface_UpdateStatus_Call) RunAndReturn(run func(*models.SaleOffer, enums.Status) error) *SaleOfferRepositoryInterface_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewSaleOfferRepositoryInterface creates a new instance of SaleOfferRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSaleOfferRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *SaleOfferRepositoryInterface {
	mock := &SaleOfferRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/agamstawn/inventory-service/internal/models"
	"github.com/agamstawn/inventory-service/internal/service"
	"github.com/agamstawn/inventory-service/internal/worker"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(p *models.Product) error {
	args := m.Called(p)
	return args.Error(0)
}
func (m *MockProductRepo) FindAll(limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}
func (m *MockProductRepo) FindByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockProductRepo) Update(p *models.Product) error {
	args := m.Called(p)
	return args.Error(0)
}
func (m *MockProductRepo) AdjustStock(id uint, delta int, t models.MovementType, reason string) (*models.Product, error) {
	args := m.Called(id, delta, t, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockProductRepo) GetStockHistory(productID uint) ([]models.StockMovement, error) {
	args := m.Called(productID)
	return args.Get(0).([]models.StockMovement), args.Error(1)
}

func newTestPool() *worker.WorkerPool {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	return worker.NewWorkerPool(db, 1, 10)
}

func TestCreateProduct_Success(t *testing.T) {
	repo := new(MockProductRepo)
	svc  := service.NewInventoryService(repo, newTestPool())

	repo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

	req := models.CreateProductRequest{
		Name:  "Widget A",
		SKU:   "SKU-001",
		Price: 29.99,
		Stock: 100,
	}

	product, err := svc.CreateProduct(req)
	assert.NoError(t, err)
	assert.Equal(t, "Widget A", product.Name)
	assert.Equal(t, 10, product.LowStockAt)
	repo.AssertExpectations(t)
}

func TestGetProduct_NotFound(t *testing.T) {
	repo := new(MockProductRepo)
	svc  := service.NewInventoryService(repo, newTestPool())

	repo.On("FindByID", uint(999)).Return(nil, nil)

	_, err := svc.GetProduct(999)
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

func TestDeductStock_InsufficientStock(t *testing.T) {
	repo := new(MockProductRepo)
	svc  := service.NewInventoryService(repo, newTestPool())

	repo.On("AdjustStock", uint(1), -999, models.MovementDeduct, "sale").
		Return(nil, errors.New("insufficient stock"))

	_, err := svc.DeductStock(1, 999, "sale")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
}

func TestDeductStock_TriggersAlert(t *testing.T) {
	repo := new(MockProductRepo)
	svc  := service.NewInventoryService(repo, newTestPool())

	lowStockProduct := &models.Product{ID: 1, Name: "Widget", Stock: 5, LowStockAt: 10}
	repo.On("AdjustStock", uint(1), -95, models.MovementDeduct, "bulk sale").
		Return(lowStockProduct, nil)

	product, err := svc.DeductStock(1, 95, "bulk sale")
	assert.NoError(t, err)
	assert.Equal(t, 5, product.Stock)

}
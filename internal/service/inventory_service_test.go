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


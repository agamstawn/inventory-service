package repository

import (
	"errors"
	"github.com/agamstawn/inventory-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(p *models.Product) error
	FindAll(limit, offset int) ([]models.Product, int64, error)
	FindByID(id uint) (*models.Product, error)
	Update(p *models.Product) error
}



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

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(p *models.Product) error {
	return r.db.Create(p).Error
}

func (r *productRepo) FindAll(limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	r.db.Model(&models.Product{}).Count(&total)
	result := r.db.Limit(limit).Offset(offset).Find(&products)
	return products, total, result.Error
}

func (r *productRepo) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	result := r.db.First(&product, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, result.Error
}

func (r *productRepo) Update(p *models.Product) error {
	return r.db.Save(p).Error
}

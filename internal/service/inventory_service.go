package service

import (
	"errors"

	"github.com/agamstawn/inventory-service/internal/models"
	"github.com/agamstawn/inventory-service/internal/repository"
	"github.com/agamstawn/inventory-service/internal/worker"
)

type InventoryService struct {
	repo repository.ProductRepository
	pool *worker.WorkerPool
}

func NewInventoryService(repo repository.ProductRepository, pool *worker.WorkerPool) *InventoryService {
	return &InventoryService{repo: repo, pool: pool}
}

func (s *InventoryService) CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	if req.LowStockAt == 0 {
		req.LowStockAt = 10 // sensible default
	}
	p := &models.Product{
		Name:        req.Name,
		SKU:         req.SKU,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		LowStockAt:  req.LowStockAt,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *InventoryService) ListProducts(page, pageSize int) ([]models.Product, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindAll(pageSize, offset)
}

func (s *InventoryService) GetProduct(id uint) (*models.Product, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (s *InventoryService) AddStock(id uint, qty int, reason string) (*models.Product, error) {
	p, err := s.repo.AdjustStock(id, qty, models.MovementAdd, reason)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *InventoryService) DeductStock(id uint, qty int, reason string) (*models.Product, error) {
	p, err := s.repo.AdjustStock(id, -qty, models.MovementDeduct, reason)
	if err != nil {
		return nil, err
	}

	if p.Stock <= p.LowStockAt {
		s.pool.Enqueue(worker.AlertJob{
			Product:   *p,
			Stock:     p.Stock,
			Threshold: p.LowStockAt,
		})
	}
	return p, nil
}

func (s *InventoryService) GetStockHistory(productID uint) ([]models.StockMovement, error) {
	return s.repo.GetStockHistory(productID)
}
package models

import "time"

type MovementType string

const (
	MovementAdd    MovementType = "add"
	MovementDeduct MovementType = "deduct"
)

// Product is the core inventory entity
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	SKU         string    `gorm:"size:100;uniqueIndex;not null" json:"sku"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	LowStockAt  int       `gorm:"default:10" json:"low_stock_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockMovement tracks every stock change — audit trail
type StockMovement struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	ProductID  uint         `gorm:"not null;index" json:"product_id"`
	Type       MovementType `gorm:"size:10;not null" json:"type"`
	Quantity   int          `gorm:"not null" json:"quantity"`
	StockAfter int          `gorm:"not null" json:"stock_after"`
	Reason     string       `gorm:"size:255" json:"reason"`
	CreatedAt  time.Time    `json:"created_at"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// StockAlert represents a low-stock notification
type StockAlert struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	Stock     int       `json:"stock"`
	Threshold int       `json:"threshold"`
	Resolved  bool      `gorm:"default:false" json:"resolved"`
	CreatedAt time.Time `json:"created_at"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	LowStockAt  int     `json:"low_stock_at" binding:"min=1"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	LowStockAt  int     `json:"low_stock_at" binding:"omitempty,min=1"`
}

type StockAdjustRequest struct {
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Reason   string `json:"reason" binding:"required"`
}

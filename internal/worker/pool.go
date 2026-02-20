package worker

import (
	"context"
	"log"
	"time"

	"github.com/agamstawn/inventory-service/internal/models"
	"gorm.io/gorm"
)

type AlertJob struct {
	Product   models.Product
	Stock     int
	Threshold int
}

type WorkerPool struct {
	jobs    chan AlertJob
	db      *gorm.DB
	workers int
}

func NewWorkerPool(db *gorm.DB, workers int, bufferSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan AlertJob, bufferSize),
		db:      db,
		workers: workers,
	}
}
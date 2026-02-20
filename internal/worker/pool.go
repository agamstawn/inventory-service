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

func (wp *WorkerPool) Start(ctx context.Context) {
	log.Printf("[WorkerPool] Starting %d workers", wp.workers)

	for i := 0; i < wp.workers; i++ {
		go wp.runWorker(ctx, i)
	}

	<-ctx.Done() 
	log.Println("[WorkerPool] Shutting down")
}

func (wp *WorkerPool) runWorker(ctx context.Context, id int) {
	log.Printf("[Worker-%d] Started", id)
	for {
		select {
		case job := <-wp.jobs:
			wp.processAlert(job)
		case <-ctx.Done():
			log.Printf("[Worker-%d] Stopped", id)
			return
		}
	}
}

func (wp *WorkerPool) processAlert(job AlertJob) {
	log.Printf("[Alert] Product '%s' (ID:%d) low stock: %d (threshold: %d)",
		job.Product.Name, job.Product.ID, job.Stock, job.Threshold)

	alert := models.StockAlert{
		ProductID: job.Product.ID,
		Stock:     job.Stock,
		Threshold: job.Threshold,
	}
	if err := wp.db.Create(&alert).Error; err != nil {
		log.Printf("[Alert] Failed to save alert: %v", err)
		return
	}

	log.Printf("[Alert] Saved alert ID %d for product %d", alert.ID, job.Product.ID)
}
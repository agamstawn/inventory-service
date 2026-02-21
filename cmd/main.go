package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/agamstawn/inventory-service/internal/api"
	"github.com/agamstawn/inventory-service/internal/models"
	"github.com/agamstawn/inventory-service/internal/repository"
	"github.com/agamstawn/inventory-service/internal/service"
	"github.com/agamstawn/inventory-service/internal/worker"
)

func main() {
	_ = godotenv.Load()

	dsn := getEnv("DATABASE_URL",
		"host=localhost user=postgres password=postgres dbname=inventory_db port=5432 sslmode=disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Product{},
		&models.StockMovement{},
		&models.StockAlert{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool    := worker.NewWorkerPool(db, 3, 100)
	repo    := repository.NewProductRepository(db)
	svc     := service.NewInventoryService(repo, pool)
	handler := api.NewHandler(svc)
	router  := api.SetupRouter(handler)

	go pool.Start(ctx)

	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}
	log.Println("Server stopped cleanly")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

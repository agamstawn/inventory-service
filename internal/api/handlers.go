package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/agamstawn/inventory-service/internal/models"
	"github.com/agamstawn/inventory-service/internal/service"
)

type Handler struct {
	svc *service.InventoryService
}

func NewHandler(svc *service.InventoryService) *Handler {
	return &Handler{svc: svc}
}

func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", h.HealthCheck)

		products := v1.Group("/products")
		{
			products.POST("",      h.CreateProduct)
			products.GET("",       h.ListProducts)
			products.GET("/:id",   h.GetProduct)
			products.PUT("/:id",   h.UpdateProduct)
		}

		stock := v1.Group("/stock")
		{
			stock.POST("/:id/add",      h.AddStock)
			stock.POST("/:id/deduct",   h.DeductStock)
			stock.GET("/:id/history",   h.GetStockHistory)
		}
	}

	return r
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "inventory-service"})
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.svc.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *Handler) ListProducts(c *gin.Context) {
	page, _     := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	products, total, err := h.svc.ListProducts(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     products,
	})
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	product, err := h.svc.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "implement now"})
}

func (h *Handler) AddStock(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	var req models.StockAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.svc.AddStock(id, req.Quantity, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) DeductStock(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	var req models.StockAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.svc.DeductStock(id, req.Quantity, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) GetStockHistory(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	history, err := h.svc.GetStockHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, err
	}
	return uint(id), nil
}

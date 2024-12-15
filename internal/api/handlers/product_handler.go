package handlers

import (
    "net/http"
    "strconv"
    "encoding/json"
    "log"
    "github.com/gin-gonic/gin"
    "product-management/internal/db"
    "product-management/internal/services"
)

func CreateProduct(c *gin.Context, messageService services.MessageService) {
    var product db.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    log.Printf("Product Data: %+v", product)
    // Save product to DB
    if err := db.CreateProductInDB(&product); err != nil {
        log.Printf("Error saving to DB: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
        return
    }

     // Publish message to RabbitMQ (if ImageURL exists)
    if product.ImageURL != "" {
        messageService.PublishMessage("image-processing", product.ImageURL)
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

func GetProductByID(c *gin.Context, cacheService services.CacheService) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Try to get product from cache
    cachedProduct, err := cacheService.GetFromCache(c.Param("id"))
    if err == nil {
        c.JSON(http.StatusOK, gin.H{"product": cachedProduct})
        return
    }

    // Get product from DB
    product, err := db.GetProductByIDFromDB(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    // Cache the entire product data (not just the name)
    productJSON, err := json.Marshal(product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize product"})
        return
    }

    // Save serialized product to the cache
    cacheService.SetToCache(c.Param("id"), string(productJSON)) // You can serialize the product object here

    c.JSON(http.StatusOK, gin.H{"product": product})
}

func GetProductsByUser(c *gin.Context, cacheService services.CacheService) {
    userID, err := strconv.Atoi(c.Query("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    products, err := db.GetProductsByUserFromDB(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"products": products})
}

package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadImage handles image uploads and provides a direct access URL
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file") // Ensure the key matches what's sent from Postman
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image"})
		return
	}

	id := uuid.New().String()
	filename := id + filepath.Ext(file.Filename)
	savePath := filepath.Join("uploads", filename)

	// Save file temporarily
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Generate direct link (adjust if serving from a real CDN)
	baseURL := fmt.Sprintf("%s://%s/uploads/%s", c.Request.URL.Scheme, c.Request.Host, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Upload successful",
		"id":        id,
		"path":      savePath,
		"directURL": baseURL,
	})
}

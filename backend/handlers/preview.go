package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"backend/services"
)

func PreviewHandler(c *gin.Context) {
	id := c.Param("id")
	imagePath := filepath.Join("uploads", id+".png") // Append .png

	numColors := 5 
	if c.Query("color") != "" {
		numColors, _ = strconv.Atoi(c.Query("color"))
	}

	previewData, err := services.GeneratePreview(imagePath, numColors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate preview"})
		return
	}

	c.Data(http.StatusOK, "image/png", previewData)
}

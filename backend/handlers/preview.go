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
	numColors, err := strconv.Atoi(c.DefaultQuery("colors", "16"))
	if err != nil || numColors < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid color number"})
		return
	}

	imagePath := filepath.Join("uploads", id+".png")

	previewData, err := services.GeneratePreview(imagePath, numColors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(previewData)
}

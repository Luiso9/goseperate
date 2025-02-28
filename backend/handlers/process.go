package handlers

import (
	"backend/services"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ProcessImage(c *gin.Context) {
	id := c.Param("id")
	numColors, err := strconv.Atoi(c.DefaultQuery("colors", "5"))
	if err != nil || numColors < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid color number"})
		return
	}

	imagePath := filepath.Join("uploads", id+".png")
	outputDir := filepath.Join("extracted", id)

	extractedFiles, err := services.ExtractColors(imagePath, outputDir, numColors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"message":   "Extraction complete",
		"layers":    len(extractedFiles),
		"extracted": extractedFiles,
	})
}

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

	d, err := strconv.Atoi(c.DefaultQuery("d", "9"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid d value"})
		return
	}

	sigmaColor, err := strconv.Atoi(c.DefaultQuery("sigmaColor", "75"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sigmaColor value"})
		return
	}

	sigmaSpace, err := strconv.Atoi(c.DefaultQuery("sigmaSpace", "75"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sigmaSpace value"})
		return
	}

	imagePath := filepath.Join("uploads", id+".png")
	outputDir := filepath.Join("extracted", id)

	extractedFiles, err := services.ExtractColors(imagePath, outputDir, numColors, d, sigmaColor, sigmaSpace)
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


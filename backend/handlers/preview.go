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
	numColors, err := strconv.Atoi(c.DefaultQuery("colors", "8"))
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

	previewData, err := services.GeneratePreview(imagePath, numColors, d, sigmaColor, sigmaSpace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(previewData)
}

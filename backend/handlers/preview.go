package handlers

import (
	"net/http"
	"strconv"

  	"github.com/gin-gonic/gin"
	"backend/services"
	"path/filepath"
)

func PreviewHandler(c *gin.Context) {
	id := c.Param("filename") 
	imagePath := filepath.Join("uploads", id+".png") /

	numColors, err := strconv.Atoi(c.DefaultQuery("color", "5"))
	if err != nil || numColors <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid color parameter"})
		return
	}

	previewData, err := services.GeneratePreview(imagePath, numColors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating preview: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", previewData)
}


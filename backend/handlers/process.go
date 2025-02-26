package handlers

import (
	"net/http"
	"fmt"
	"path/filepath"
	"strconv"

	"backend/services"
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

	// Debugging logs
	fmt.Println("[DEBUG] Processing image:", imagePath)
	fmt.Println("[DEBUG] Extracting to directory:", outputDir)
	fmt.Println("[DEBUG] Number of colors:", numColors)

	err = services.ExtractColors(imagePath, outputDir, numColors)
	if err != nil {
		fmt.Println("[ERROR] Extraction failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image", "details": err.Error()})
		return
	}

	fmt.Println("[DEBUG] Extraction complete for:", id)
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Extraction complete", "layers": numColors})
}

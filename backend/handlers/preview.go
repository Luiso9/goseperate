package handlers

import (
	"net/http"
	"strconv"

  "github.com/gin-gonic/gin"
	"backend/services"
  "backend/utils"
)

func PreviewHandler(c *gin.Context) {
  filename := c.Param("filename")

  imagePath, err := utils.FindExistingFile("backend/uploads", filename)
  if imagePath == "" {
    c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
    return
  }

  numColors, err := strconv.Atoi(c.Query("color"))
  if err != nil || numColors <= 0 {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
    return
  }

  previewData, err := services.GeneratePreview(imagePath, numColors)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating preview" + err.Error()})
    return
  }

  c.Data(http.StatusOK, "image/png", previewData)
} 

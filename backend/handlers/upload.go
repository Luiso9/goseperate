package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err == nil {
		id := uuid.New().String()
		filename := id + filepath.Ext(file.Filename)
		savePath := filepath.Join("uploads", filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Upload successful",
			"id":      id,
			"path":    savePath,
		})
		return
	}

	// If no file was uploaded, check for image_url
	var json struct {
		ImageURL string `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&json); err != nil || json.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Provide either 'image' file or 'image_url'."})
		return
	}

	// Fetch image from URL
	resp, err := http.Get(json.ImageURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to download image from URL"})
		return
	}
	defer resp.Body.Close()

	id := uuid.New().String()
	ext := filepath.Ext(json.ImageURL)
	if ext == "" {
		ext = ".jpg" 
	}
	filename := id + ext
	savePath := filepath.Join("uploads", filename)

	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image to file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image downloaded successfully",
		"id":      id,
		"path":    savePath,
	})
}

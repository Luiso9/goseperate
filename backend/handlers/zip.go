package handlers

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DownloadZip(c *gin.Context) {
	id := c.Param("id")
	sourceDir := filepath.Join("extracted", id)
	zipPath := sourceDir + ".zip"

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "No extracted files found"})
		return
	}

	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ZIP file"})
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		zipWriter.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read extracted files"})
		return
	}

	for _, file := range files {
		filePath := filepath.Join(sourceDir, file.Name())
		srcFile, err := os.Open(filePath)
		if err != nil {
			zipWriter.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open extracted file"})
			return
		}

		zipEntry, err := zipWriter.Create(file.Name())
		if err != nil {
			srcFile.Close()
			zipWriter.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file to ZIP"})
			return
		}

		_, err = io.Copy(zipEntry, srcFile)
		srcFile.Close()
		if err != nil {
			zipWriter.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file to ZIP"})
			return
		}
	}

	err = zipWriter.Close() // Explicitly check for errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize ZIP file"})
		return
	}

	// Set proper headers
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+id+".zip")
	c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(zipPath)

	// Cleanup after serving the file
	go func() {
		os.RemoveAll(sourceDir)
		os.Remove(zipPath)
	}()
}

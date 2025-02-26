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
	zipPath := filepath.Join("extracted", id+".zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ZIP file"})
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	files, _ := os.ReadDir(filepath.Join("extracted", id))
	for _, file := range files {
		f, _ := os.Open(filepath.Join("extracted", id, file.Name()))
		defer f.Close()

		w, _ := zipWriter.Create(file.Name())
		io.Copy(w, f) // Copy file data to ZIP
	}

	c.File(zipPath)

	// Cleanup
	os.RemoveAll(filepath.Join("extracted", id))
	os.Remove(zipPath)
}

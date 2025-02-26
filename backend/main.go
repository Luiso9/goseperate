package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"backend/handlers"
)

func main() {
	os.MkdirAll("uploads", os.ModePerm)
	os.MkdirAll("extracted", os.ModePerm)

	router := gin.Default()

	// Routes
	router.Static("/uploads", "./uploads")
	router.POST("/upload", handlers.UploadImage)
	router.POST("/process/:id", handlers.ProcessImage)
	router.GET("/download/:id", handlers.DownloadZip)

	log.Println("Server running on port 9330")
	router.Run(":9330")
}

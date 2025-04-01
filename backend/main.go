package main

import (
	"log"
	"os"

	"backend/handlers"
	"github.com/gin-gonic/gin"
	"time"
	"backend/utils"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	os.MkdirAll("uploads", os.ModePerm)
	os.MkdirAll("extracted", os.ModePerm)

	router := gin.Default()

	utils.StartCleanupRoutine(1*time.Minute, 10*time.Minute)
	// Routes
	router.Static("/uploads", "./uploads")
  	router.GET("/preview/:id", handlers.PreviewHandler)
	router.POST("/upload", handlers.UploadImage)
	router.POST("/process/:id", handlers.ProcessImage)
	router.GET("/download/:id", handlers.DownloadZip)

	log.Println("Server running on port 9330")
	router.Run(":9330")
}

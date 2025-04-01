package main

import (
	"log"
	"os"
	"time"
	"backend/handlers"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	os.MkdirAll("uploads", os.ModePerm)
	os.MkdirAll("extracted", os.ModePerm)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	utils.StartCleanupRoutine(1*time.Minute, 10*time.Minute)

	router.Static("/uploads", "./uploads")
	router.GET("/preview/:id", handlers.PreviewHandler)
	router.POST("/upload", handlers.UploadImage)
	router.POST("/process/:id", handlers.ProcessImage)
	router.GET("/download/:id", handlers.DownloadZip)

	log.Println("Server running on port 9330")
	router.Run(":9330")
}

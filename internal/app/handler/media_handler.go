package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UploadMedia handles file uploads to Google Cloud Storage and returns a signed URL.
func UploadMedia(c *gin.Context) {
	// Parse uploaded file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing file",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	// Read bucket name from environment
	bucket := os.Getenv("GCS_BUCKET_NAME")
	if bucket == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "GCS_BUCKET_NAME is not set in environment",
		})
		return
	}

	// Debug logging
	fmt.Println("📤 Uploading file:", fileHeader.Filename)
	fmt.Println("🧾 Content-Type:", fileHeader.Header.Get("Content-Type"))
	fmt.Println("🎯 Target bucket:", bucket)

	// Upload to GCS and get signed URL
	url, err := utils.UploadToGCS(file, fileHeader, bucket)
	if err != nil {
		fmt.Println("❌ Upload error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Upload failed",
			"details": err.Error(),
		})
		return
	}

	// Confirm success
	fmt.Println("✅ Upload successful:", url)

	// Respond with signed URL
	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"message": "Upload successful",
	})
}


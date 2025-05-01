package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

	// Optional metadata (can be blank)
	courseID := c.PostForm("course_id")
	lessonID := c.PostForm("lesson_id")
	questionID := c.PostForm("question_id")

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

	// Upload to GCS and get object path
	objectPath, err := utils.UploadToGCS(file, fileHeader, bucket, courseID, lessonID, questionID)
	if err != nil {
		fmt.Println("❌ Upload error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Upload failed",
			"details": err.Error(),
		})
		return
	}

	fmt.Println("✅ Upload successful:", objectPath)

	c.JSON(http.StatusOK, gin.H{
		"url":     objectPath,
		"message": "Upload successful",
	})
}

// GetSignedURL dynamically generates a signed URL for a GCS object.
func GetSignedURL(c *gin.Context) {
	object := c.Query("object")
	if object == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing object query parameter",
		})
		return
	}

	bucket := os.Getenv("GCS_BUCKET_NAME")
	if bucket == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "GCS_BUCKET_NAME is not set in environment",
		})
		return
	}

	url, err := utils.GenerateSignedURL(bucket, object, 15*time.Minute)
	if err != nil {
		fmt.Println("❌ Failed to generate signed URL:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate signed URL",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}


package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ✅ Still used: dynamically generate signed URL for viewing a GCS object
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

// ✅ Used for direct-to-GCS uploads from frontend
func GetUploadURL(c *gin.Context) {
	filename := c.Query("filename")
	contentType := c.Query("type")
	courseID := c.Query("course_id")
	lessonID := c.Query("lesson_id")
	questionID := c.Query("question_id")

	if filename == "" || contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filename or type"})
		return
	}

	// Build structured path
	objectName := "uploads"
	if courseID != "" {
		objectName += fmt.Sprintf("/course_%s", courseID)
	}
	if lessonID != "" {
		objectName += fmt.Sprintf("/lesson_%s", lessonID)
	}
	if questionID != "" {
		objectName += fmt.Sprintf("/question_%s", questionID)
	}

	objectName += fmt.Sprintf("/%s_%s", utils.GenerateUUID(), filename)

	bucket := os.Getenv("GCS_BUCKET_NAME")

	signedURL, err := utils.GenerateV4UploadURL(bucket, objectName, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate signed URL",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":        signedURL,
		"objectName": objectName,
	})
}

package handler

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type MediaHandler struct {
	bucketName string
	client     *storage.Client
}

func NewMediaHandler(bucketName string, credentialsPath string) (*MediaHandler, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}
	return &MediaHandler{
		bucketName: bucketName,
		client:     client,
	}, nil
}

func (h *MediaHandler) GenerateSignedUploadURL(c *gin.Context) {
	fileName := c.Query("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileName query param is required"})
		return
	}

	url, err := h.generateSignedURL(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *MediaHandler) generateSignedURL(fileName string) (string, error) {
	opts := &storage.SignedURLOptions{
		Method:         "PUT",
		Expires:        time.Now().Add(15 * time.Minute),
		ContentType:    "application/octet-stream",
		GoogleAccessID: "your-service-account@gcp-project.iam.gserviceaccount.com",
		PrivateKey:     []byte("..."), // Better to load from service account file
	}
	return storage.SignedURL(h.bucketName, fileName, opts)
}

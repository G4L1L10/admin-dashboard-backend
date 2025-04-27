package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

// UploadToGCS uploads a file to GCS and returns the GCS object name (not signed URL).
func UploadToGCS(file multipart.File, fileHeader *multipart.FileHeader, bucketName string) (string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	// Generate a unique object name
	objectName := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	// Prepare GCS writer
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	// Use the content type if provided, or fall back
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	wc.ContentType = contentType

	// Upload file data
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write file to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	// ✅ Return the object path (not a signed URL)
	return objectName, nil
}


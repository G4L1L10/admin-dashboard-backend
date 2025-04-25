package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

// UploadToGCS uploads a file to GCS and returns a signed URL valid for 15 minutes.
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
	_, copyErr := io.Copy(wc, file)
	if copyErr != nil {
		return "", fmt.Errorf("failed to write file to GCS: %w", copyErr)
	}
	closeErr := wc.Close()
	if closeErr != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", closeErr)
	}

	// Generate a signed URL
	signedURL, signErr := generateSignedURL(bucketName, objectName, 15*time.Minute)
	if signErr != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", signErr)
	}

	return signedURL, nil
}


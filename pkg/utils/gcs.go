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
func UploadToGCS(
	file multipart.File,
	fileHeader *multipart.FileHeader,
	bucketName string,
	courseID string,
	lessonID string,
	questionID string,
) (string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	// Build custom folder structure
	folder := "uploads"
	if courseID != "" {
		folder += fmt.Sprintf("/course_%s", courseID)
	}
	if lessonID != "" {
		folder += fmt.Sprintf("/lesson_%s", lessonID)
	}
	if questionID != "" {
		folder += fmt.Sprintf("/question_%s", questionID)
	}

	// Generate unique file name
	objectName := fmt.Sprintf("%s/%d_%s", folder, time.Now().UnixNano(), fileHeader.Filename)

	// Prepare GCS writer
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	// Content-Type fallback
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	wc.ContentType = contentType

	// Upload the file
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write file to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return objectName, nil
}

// DeleteFromGCS deletes an object from GCS by its object name (not URL).
func DeleteFromGCS(ctx context.Context, bucketName, objectPath string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	obj := client.Bucket(bucketName).Object(objectPath)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectPath, err)
	}
	return nil
}

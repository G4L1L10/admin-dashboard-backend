package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
)

// GenerateV4UploadURL generates a V4 signed URL for uploading a file to GCS.
func GenerateV4UploadURL(bucketName, objectName, contentType string) (string, error) {
	signerEmail := os.Getenv("GCS_SIGNER_EMAIL")
	if signerEmail == "" {
		return "", fmt.Errorf("GCS_SIGNER_EMAIL is not set in environment")
	}

	ctx := context.Background()

	iamClient, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create IAM credentials client: %w", err)
	}
	defer iamClient.Close()

	signBytes := func(b []byte) ([]byte, error) {
		resp, signErr := iamClient.SignBlob(ctx, &credentialspb.SignBlobRequest{
			Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", signerEmail),
			Payload: b,
		})
		if signErr != nil {
			return nil, fmt.Errorf("failed to sign blob: %w", signErr)
		}
		return resp.SignedBlob, nil
	}

	opts := &storage.SignedURLOptions{
		Method:         http.MethodPut,
		Expires:        time.Now().Add(15 * time.Minute),
		Scheme:         storage.SigningSchemeV4,
		ContentType:    contentType,
		GoogleAccessID: signerEmail,
		SignBytes:      signBytes,
	}

	signedURL, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed upload URL: %w", err)
	}

	return signedURL, nil
}


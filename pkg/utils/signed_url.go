package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
)

func GenerateSignedURL(bucketName, objectName string, expiresIn time.Duration) (string, error) {
	ctx := context.Background()

	serviceAccountEmail := os.Getenv("GCS_SIGNER_EMAIL")
	if serviceAccountEmail == "" {
		return "", fmt.Errorf("GCS_SIGNER_EMAIL not set in env")
	}

	// Create IAM credentials client to sign bytes
	iamClient, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create IAM credentials client: %w", err)
	}
	defer iamClient.Close()

	// Set up a custom SignBytes function using IAM credentials API
	signBytes := func(b []byte) ([]byte, error) {
		var resp *credentialspb.SignBlobResponse // Declare resp first

		resp, err = iamClient.SignBlob(ctx, &credentialspb.SignBlobRequest{
			Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail),
			Payload: b,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to sign blob: %w", err)
		}
		return resp.SignedBlob, nil
	}

	opts := &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(expiresIn),
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: serviceAccountEmail,
		SignBytes:      signBytes,
	}

	url, err := storage.SignedURL(bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create signed URL: %w", err)
	}

	return url, nil
}


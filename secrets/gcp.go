package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func GetVersion(name, version string) (string, error) {
	if version == "" {
		version = "latest"
	}

	projectID := os.Getenv("PROJECT_ID")

	if projectID == "" {
		return "", errors.New("PROJECT_ID environment variable must be specified")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, name, version),
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}

func Get(name string) (string, error) {
	return GetVersion(name, "")
}

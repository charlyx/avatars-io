package secrets

import (
	"context"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type GCPClient struct {
	secretManager *secretmanager.Client
	projectID     string
	context       context.Context
}

func NewGCPClient(projectID string) (*GCPClient, error) {
	if projectID == "" {
		return nil, errors.New("projectID must not be empty.")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	return &GCPClient{
		secretManager: client,
		projectID:     projectID,
		context:       ctx,
	}, nil
}

func (c *GCPClient) GetVersion(name, version string) (string, error) {
	if version == "" {
		version = "latest"
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", c.projectID, name, version),
	}

	result, err := c.secretManager.AccessSecretVersion(c.context, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}

func (c *GCPClient) Get(name string) (string, error) {
	return c.GetVersion(name, "")
}

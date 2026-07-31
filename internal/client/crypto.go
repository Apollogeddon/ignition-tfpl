package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EncryptSecret encrypts a plaintext string using the gateway's encryption endpoint
func (c *Client) EncryptSecret(ctx context.Context, plaintext string) (*IgnitionSecret, error) {
	bodyBytes, err := c.doRequest(ctx, http.MethodPost, "/data/api/v1/encryption/encrypt", []byte(plaintext), "text/plain")
	if err != nil {
		return nil, err
	}

	var rawJwe map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawJwe); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted response: %w", err)
	}

	return &IgnitionSecret{
		Type: "Embedded",
		Data: rawJwe,
	}, nil
}

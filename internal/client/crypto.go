package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// EncryptSecret encrypts a plaintext string using the gateway's encryption endpoint
func (c *Client) EncryptSecret(ctx context.Context, plaintext string) (*IgnitionSecret, error) {
	path := "/data/api/v1/encryption/encrypt"

	// Use a custom request to handle text/plain content type for the plaintext body
	fullURL := c.HostURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Set body manually to handle plain text string
	req.Body = io.NopCloser(strings.NewReader(plaintext))
	req.ContentLength = int64(len(plaintext))

	req.Header.Set("X-Ignition-API-Token", c.Token)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("encryption failed with status: %d, body: %s", res.StatusCode, bodyBytes)
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

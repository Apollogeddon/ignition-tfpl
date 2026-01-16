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
		
		// Send plaintext directly as body with text/plain content type
		// Note: doRequest uses application/json by default, so we might need to override headers or use a custom request
		// but checking doRequest implementation... it sets Content-Type to application/json. 
		// We need to handle this.
		// For now, let's assume we can modify doRequest or just make a new request here since it's a specific case.
		
		// Let's implement a specific request here to ensure correct headers
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
		defer res.Body.Close()

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
		},
		nil
	}

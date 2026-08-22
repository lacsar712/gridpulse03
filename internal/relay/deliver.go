package relay

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func PostOutbound(ctx context.Context, endpoint string, payload []byte) error {
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-gridpulse-Entity", "Trip")
	resp, err := client.Do(req)
	if err != nil {
		return WrapRetryable("outbound", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 500 {
		return WrapRetryable("outbound", fmt.Errorf("status %d", resp.StatusCode))
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("outbound client error: %d", resp.StatusCode)
	}
	return nil
}

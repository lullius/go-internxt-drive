package buckets

import (
	"fmt"
	"io"
	"net/http"
)

func Transfer(part UploadPart, r io.Reader, size int64) error {
	req, err := http.NewRequest("PUT", part.URL, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = size
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("transfer failed: status %d, %s", resp.StatusCode, string(body))
	}
	return nil
}

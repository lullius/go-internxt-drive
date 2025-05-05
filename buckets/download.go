package buckets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/StarHack/go-internxt-drive/config"
)

const bucketAPIBase = "https://api.internxt.com"

// ShardInfo mirrors the per‑shard info returned by /files/{fileID}/info
type ShardInfo struct {
	Index int    `json:"index"`
	Hash  string `json:"hash"`
	URL   string `json:"url"`
}

// BucketFileInfo is the metadata returned by GET /buckets/{bucketID}/files/{fileID}/info
type BucketFileInfo struct {
	Bucket   string      `json:"bucket"`
	Index    string      `json:"index"`
	Size     int64       `json:"size"`
	Version  int         `json:"version"`
	Created  string      `json:"created"`
	Renewal  string      `json:"renewal"`
	Mimetype string      `json:"mimetype"`
	Filename string      `json:"filename"`
	ID       string      `json:"id"`
	Shards   []ShardInfo `json:"shards"`
}

// GetBucketFileInfo calls the correct /info endpoint and parses its JSON.
func GetBucketFileInfo(cfg *config.Config, bucketID, fileID string) (*BucketFileInfo, error) {
	url := fmt.Sprintf("%s/buckets/%s/files/%s/info", bucketAPIBase, bucketID, fileID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", cfg.BasicAuthHeader)
	req.Header.Set("internxt-version", "1.0")
	req.Header.Set("internxt-client", "internxt-go-sdk")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file info fetch failed: %d %s", resp.StatusCode, string(body))
	}

	var info BucketFileInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DownloadFile downloads and decrypts the first shard of the given file.
func DownloadFile(cfg *config.Config, fileID, destPath string) error {
	// 1) fetch file info from the bucket API
	info, err := GetBucketFileInfo(cfg, cfg.Bucket, fileID)
	if err != nil {
		return err
	}
	if len(info.Shards) == 0 {
		return fmt.Errorf("no shards found for file %s", fileID)
	}
	shard := info.Shards[0]

	// 2) derive fileKey+iv using the stored index (hex of random index)
	key, iv, err := GenerateFileKey(cfg.Mnemonic, cfg.Bucket, info.Index)
	if err != nil {
		return err
	}

	// 3) GET the encrypted shard directly from its presigned URL
	resp, err := http.Get(shard.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("shard download failed: %d %s", resp.StatusCode, string(body))
	}

	// 4) wrap in AES‑CTR decryptor
	decReader, err := DecryptReader(resp.Body, key, iv)
	if err != nil {
		return err
	}

	// 5) write plaintext to file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, decReader); err != nil {
		return err
	}
	return nil
}

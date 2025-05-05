package buckets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/StarHack/go-internxt-drive/config"
)

const API = "https://api.internxt.com"

// UploadPartSpec defines each part’s index and size for the start call
type UploadPartSpec struct {
	Index int   `json:"index"`
	Size  int64 `json:"size"`
}

type startUploadReq struct {
	Uploads []UploadPartSpec `json:"uploads"`
}

type UploadPart struct {
	Index int    `json:"index"`
	UUID  string `json:"uuid"`
	URL   string `json:"url"`
}

type StartUploadResp struct {
	Uploads []UploadPart `json:"uploads"`
}

// StartUpload reserves all parts at once
func StartUpload(cfg *config.Config, bucketID string, parts []UploadPartSpec) (*StartUploadResp, error) {
	url := fmt.Sprintf("%s/v2/buckets/%s/files/start?multiparts=1", API, bucketID)
	reqBody := startUploadReq{Uploads: parts}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", cfg.BasicAuthHeader)
	req.Header.Set("internxt-version", "1.0")
	req.Header.Set("internxt-client", "drive-web")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result StartUploadResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

/*
package buckets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/StarHack/go-internxt-drive/config"
)

type startUploadReq struct {
	Uploads []struct {
		Index int   `json:"index"`
		Size  int64 `json:"size"`
	} `json:"uploads"`
}

type UploadPart struct {
	Index int    `json:"index"`
	UUID  string `json:"uuid"`
	URL   string `json:"url"`
}

type StartUploadResp struct {
	Uploads []UploadPart `json:"uploads"`
}

func StartUpload(cfg *config.Config, bucketID string, index int, size int64) (*StartUploadResp, error) {
	url := fmt.Sprintf("%s/v2/buckets/%s/files/start?multiparts=1", "https://api.internxt.com", bucketID)
	reqBody := startUploadReq{Uploads: []struct {
		Index int   `json:"index"`
		Size  int64 `json:"size"`
	}{{Index: index, Size: size}}}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader(cfg))
	req.Header.Set("internxt-version", "1.0")
	req.Header.Set("internxt-client", "drive-web")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result StartUploadResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func authHeader(cfg *config.Config) string {
	return "Basic d29sZmdhbmdAbWFpbC5jaDphY2RhMTAzY2IxNDg3NTc5NDA2OTQ1ZWJlNTNhOWJkNzA3YzVmOTM0YjY5MDJhMDc4YWRhZjc4ZDUxODQyYzVk"
}
*/

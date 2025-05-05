package files

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/StarHack/go-internxt-drive/config"
	"github.com/StarHack/go-internxt-drive/folders"
)

const filesPath = "/files"

// GetFileMeta retrieves the metadata for the file with the given UUID.
func GetFileMeta(cfg *config.Config, fileUUID string) (*folders.File, error) {
	endpoint := cfg.DriveAPIURL + filesPath + "/" + fileUUID + "/meta"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GetFileMeta failed: %d %s", resp.StatusCode, string(body))
	}

	var file folders.File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, err
	}
	return &file, nil
}

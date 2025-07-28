package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

// DeleteFile deletes a file by UUID
func DeleteFile(cfg *config.Config, uuid string) error {
	u, err := url.Parse(cfg.DriveAPIURL + filesPath + "/" + uuid)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DeleteFile failed: %d %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdateFileMeta updates the metadata of a file with the given UUID.
func UpdateFileMeta(cfg *config.Config, fileUUID string, updated *folders.File) (*folders.File, error) {
	endpoint := cfg.DriveAPIURL + filesPath + "/" + fileUUID + "/meta"

	body, err := json.Marshal(updated)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("UpdateFileMeta failed: %d %s", resp.StatusCode, string(respBody))
	}

	var updatedFile folders.File
	if err := json.NewDecoder(resp.Body).Decode(&updatedFile); err != nil {
		return nil, err
	}

	return &updatedFile, nil
}

// MoveFile moves the file with the given UUID to the destination folder.
func MoveFile(cfg *config.Config, fileUUID, destinationFolderUUID string) (*folders.File, error) {
	endpoint := cfg.DriveAPIURL + filesPath + "/" + fileUUID

	body, err := json.Marshal(map[string]string{
		"destinationFolder": destinationFolderUUID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("MoveFile failed: %d %s", resp.StatusCode, string(respBody))
	}

	var movedFile folders.File
	if err := json.NewDecoder(resp.Body).Decode(&movedFile); err != nil {
		return nil, err
	}

	return &movedFile, nil
}

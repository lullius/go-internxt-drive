package folders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/StarHack/go-internxt-drive/config"
)

const foldersPath = "/folders"

// FolderStatus represents the status filter for folder operations
// Possible values: EXISTS, TRASHED, DELETED, ALL
type FolderStatus string

const (
	StatusExists  FolderStatus = "EXISTS"
	StatusTrashed FolderStatus = "TRASHED"
	StatusDeleted FolderStatus = "DELETED"
	StatusAll     FolderStatus = "ALL"
)

// CreateFolderRequest is the payload for POST /drive/folders
type CreateFolderRequest struct {
	PlainName        string `json:"plainName"`
	ParentFolderUUID string `json:"parentFolderUuid"`
	ModificationTime string `json:"modificationTime"`
	CreationTime     string `json:"creationTime"`
}

// Folder represents the response from POST/GET /drive/folders
type Folder struct {
	Type             string      `json:"type"`
	ID               int64       `json:"id"`
	ParentID         int64       `json:"parentId"`
	ParentUUID       string      `json:"parentUuid"`
	Name             string      `json:"name"`
	Parent           interface{} `json:"parent"`
	Bucket           interface{} `json:"bucket"`
	UserID           int64       `json:"userId"`
	User             interface{} `json:"user"`
	EncryptVersion   string      `json:"encryptVersion"`
	Deleted          bool        `json:"deleted"`
	DeletedAt        *time.Time  `json:"deletedAt"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
	UUID             string      `json:"uuid"`
	PlainName        string      `json:"plainName"`
	Size             int64       `json:"size"`
	Removed          bool        `json:"removed"`
	RemovedAt        *time.Time  `json:"removedAt"`
	CreationTime     time.Time   `json:"creationTime"`
	ModificationTime time.Time   `json:"modificationTime"`
	Status           string      `json:"status"`
}

// CreateFolder calls POST {DriveAPIURL}/folders with authorization.
// It auto‑fills CreationTime/ModificationTime if empty, checks status,
// and returns the newly created Folder.
func CreateFolder(cfg *config.Config, reqBody CreateFolderRequest) (*Folder, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if reqBody.CreationTime == "" {
		reqBody.CreationTime = now
	}
	if reqBody.ModificationTime == "" {
		reqBody.ModificationTime = now
	}

	endpoint := cfg.DriveAPIURL + foldersPath
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If not 200 OK, read body and return error
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CreateFolder failed: %d %s", resp.StatusCode, string(body))
	}

	// Decode into Folder
	var folder Folder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// DeleteFolders calls DELETE {DriveAPIURL}/folders/{uuid} with authorization and deletes it
func DeleteFolder(cfg *config.Config, uuid string) error {
	u, err := url.Parse(cfg.DriveAPIURL + foldersPath + "/" + uuid)
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

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DeleteFolder failed: %d %s", resp.StatusCode, string(body))
	}

	fmt.Println("Status:", resp.Status)
	return nil
}

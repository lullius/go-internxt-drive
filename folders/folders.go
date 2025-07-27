package folders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// File represents the response object for files in a folder
// under GET /drive/folders/content/{uuid}/files
type File struct {
	ID               int64         `json:"id"`
	FileID           string        `json:"fileId"`
	UUID             string        `json:"uuid"`
	Name             string        `json:"name"`
	PlainName        string        `json:"plainName"`
	Type             string        `json:"type"`
	FolderID         json.Number   `json:"folderId"`
	FolderUUID       string        `json:"folderUuid"`
	Folder           interface{}   `json:"folder"`
	Bucket           string        `json:"bucket"`
	UserID           json.Number   `json:"userId"`
	User             interface{}   `json:"user"`
	EncryptVersion   string        `json:"encryptVersion"`
	Size             json.Number   `json:"size"`
	Deleted          bool          `json:"deleted"`
	DeletedAt        *time.Time    `json:"deletedAt"`
	Removed          bool          `json:"removed"`
	RemovedAt        *time.Time    `json:"removedAt"`
	Shares           []interface{} `json:"shares"`
	Sharings         []interface{} `json:"sharings"`
	Thumbnails       []interface{} `json:"thumbnails"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	CreationTime     time.Time     `json:"creationTime"`
	ModificationTime time.Time     `json:"modificationTime"`
	Status           string        `json:"status"`
}

// ListOptions defines common pagination and sorting parameters
// for list endpoints.
type ListOptions struct {
	Limit  int
	Offset int
	Sort   string
	Order  string
}

// CreateFolder calls {DriveAPIURL}/folders with authorization.
// It auto‑fills CreationTime/ModificationTime if empty, checks status, and returns the newly created Folder.
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CreateFolder failed: %d %s", resp.StatusCode, string(body))
	}

	var folder Folder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

// DeleteFolders deletes a folder by UUID
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

	//Server returns 204 on success
	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DeleteFolder failed: %d %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListFolders lists child folders under the given parent UUID.
// Returns a slice of folders or error otherwise
func ListFolders(cfg *config.Config, parentUUID string, opts ListOptions) ([]Folder, error) {
	base := fmt.Sprintf("%s%s/content/%s/folders", cfg.DriveAPIURL, foldersPath, parentUUID)
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()

	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}
	sortField := opts.Sort
	if sortField == "" {
		sortField = "plainName"
	}
	order := opts.Order
	if order == "" {
		order = "ASC"
	}
	q.Set("offset", strconv.Itoa(offset))
	q.Set("limit", strconv.Itoa(limit))
	q.Set("sort", sortField)
	q.Set("order", order)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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
		return nil, fmt.Errorf("ListFolders failed: %d %s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		Folders []Folder `json:"folders"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return wrapper.Folders, nil
}

// ListFiles lists files under the given parent folder UUID.
// Returns a slice of files or error otherwise
func ListFiles(cfg *config.Config, parentUUID string, opts ListOptions) ([]File, error) {
	base := fmt.Sprintf("%s%s/content/%s/files", cfg.DriveAPIURL, foldersPath, parentUUID)
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()

	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}
	sortField := opts.Sort
	if sortField == "" {
		sortField = "plainName"
	}
	order := opts.Order
	if order == "" {
		order = "ASC"
	}
	q.Set("offset", strconv.Itoa(offset))
	q.Set("limit", strconv.Itoa(limit))
	q.Set("sort", sortField)
	q.Set("order", order)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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
		return nil, fmt.Errorf("ListFiles failed: %d %s", resp.StatusCode, string(body))
	}

	var wrapper struct {
		Files []File `json:"files"`
	}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&wrapper); err != nil {
		return nil, err
	}
	return wrapper.Files, nil
}

// This function will get all of the files in a folder, getting 50 at a time until completed
func ListAllFiles(cfg *config.Config, parentUUID string) ([]File, error) {
	var outFiles []File
	offset := 0
	loops := 0
	maxLoops := 10000 //Find sane number...
	for {
		files, err := ListFiles(cfg, parentUUID, ListOptions{Offset: offset})
		if err != nil {
			return nil, err
		}
		outFiles = append(outFiles, files...)
		offset += 50
		loops += 1
		if len(files) != 50 || loops >= maxLoops {
			break
		}
	}
	return outFiles, nil
}

// This function will get all of the folders in a folder, getting 50 at a time until completed
func ListAllFolders(cfg *config.Config, parentUUID string) ([]Folder, error) {
	var outFolders []Folder
	offset := 0
	loops := 0
	maxLoops := 10000 //Find sane number...
	for {
		files, err := ListFolders(cfg, parentUUID, ListOptions{Offset: offset})
		if err != nil {
			return nil, err
		}
		outFolders = append(outFolders, files...)
		offset += 50
		loops += 1
		if len(files) != 50 || loops >= maxLoops {
			break
		}
	}
	return outFolders, nil
}

// RenameFolder updates the plainName of an existing folder.
// Returns nil on HTTP 200, or an error otherwise.
func RenameFolder(cfg *config.Config, uuid, newName string) error {
	endpoint := fmt.Sprintf("%s%s/%s/meta", cfg.DriveAPIURL, foldersPath, uuid)

	payload := struct {
		PlainName string `json:"plainName"`
	}{PlainName: newName}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("RenameFolder: marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("RenameFolder: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("RenameFolder: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("RenameFolder failed: %d %s", resp.StatusCode, string(body))
	}
	return nil
}

// MoveFolder moves a folder into a new parent.
// Returns nil on HTTP 200, or an error otherwise.
func MoveFolder(cfg *config.Config, uuid, destUUID string) error {
	endpoint := fmt.Sprintf("%s%s/%s", cfg.DriveAPIURL, foldersPath, uuid)

	payload := struct {
		DestinationFolder string `json:"destinationFolder"`
	}{DestinationFolder: destUUID}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("MoveFolder: marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("MoveFolder: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("MoveFolder: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("MoveFolder failed: %d %s", resp.StatusCode, string(body))
	}
	return nil
}

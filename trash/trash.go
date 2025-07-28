package trash

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/StarHack/go-internxt-drive/config"
	"github.com/StarHack/go-internxt-drive/folders"
)

type TrashItemsRequest struct {
	Items []TrashItem `json:"items"`
}

type PaginatedTrashResponse struct {
	Items []any `json:"items"`
	Page  int   `json:"page"`
	Total int   `json:"total"`
}

type TrashItem struct {
	File   *folders.File
	Folder *folders.Folder
}

type paginatedTrashRaw struct {
	Result []json.RawMessage `json:"result"`
}

type SortField string
type Order string
type ItemType string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"

	SortByPlainName SortField = "plainName"
	SortByUpdatedAt SortField = "updatedAt"
	SortBySize      SortField = "size"

	ItemTypeFiles   ItemType = "files"
	ItemTypeFolders ItemType = "folders"
)

func parseTrashItem(data []byte) (*TrashItem, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}

	switch probe.Type {
	case "folder":
		var f folders.Folder
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, err
		}
		return &TrashItem{Folder: &f}, nil
	case "file":
		var f folders.File
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, err
		}
		return &TrashItem{File: &f}, nil
	default:
		return nil, fmt.Errorf("unknown item type: %s", probe.Type)
	}
}

func GetPaginatedTrash(cfg *config.Config, limit, offset int, itemType ItemType, sort SortField, order Order, root bool) ([]*TrashItem, error) {
	url := fmt.Sprintf("%s/storage/trash/paginated?limit=%d&offset=%d&type=%s&root=%t&sort=%s&order=%s",
		cfg.DriveAPIURL, limit, offset, itemType, root, sort, order)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
		return nil, errors.New("failed to fetch paginated trash: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw paginatedTrashRaw
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var items []*TrashItem
	for _, r := range raw.Result {
		item, err := parseTrashItem(r)
		if err == nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func AddToTrash(cfg *config.Config, items []TrashItem) error {
	body, _ := json.Marshal(TrashItemsRequest{Items: items})

	req, err := http.NewRequest(http.MethodPost, cfg.DriveAPIURL+"/storage/trash/add", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add items to trash: %s", resp.Status)
	}
	return nil
}

func DeleteAllTrash(cfg *config.Config) error {
	req, err := http.NewRequest(http.MethodDelete, cfg.DriveAPIURL+"/storage/trash/all", nil)
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
		return fmt.Errorf("failed to delete all trash: %s", resp.Status)
	}
	return nil
}

func RequestDeleteAllTrash(cfg *config.Config) error {
	req, err := http.NewRequest(http.MethodDelete, cfg.DriveAPIURL+"/storage/trash/all/request", nil)
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
		return fmt.Errorf("failed to request delete all trash: %s", resp.Status)
	}
	return nil
}

func DeleteSpecifiedTrashItems(cfg *config.Config, items []TrashItem) error {
	body, _ := json.Marshal(TrashItemsRequest{Items: items})

	req, err := http.NewRequest(http.MethodDelete, cfg.DriveAPIURL+"/storage/trash", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete specified trash items: %s - %s", resp.Status, string(respBody))
	}
	return nil
}

func DeleteTrashFile(cfg *config.Config, fileID int64) error {
	url := fmt.Sprintf("%s/storage/trash/file/%d", cfg.DriveAPIURL, fileID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete trash file: %s", resp.Status)
	}
	return nil
}

func DeleteTrashFolder(cfg *config.Config, folderID int64) error {
	url := fmt.Sprintf("%s/storage/trash/folder/%d", cfg.DriveAPIURL, folderID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete trash folder: %s", resp.Status)
	}
	return nil
}

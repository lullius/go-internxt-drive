package buckets

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/StarHack/go-internxt-drive/config"
)

func UploadFile(cfg *config.Config, filePath string) (string, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	plainSize := int64(len(raw))
	ph := sha256.Sum256(raw)

	ph = [32]byte{}
	if _, err := rand.Read(ph[:]); err != nil {
		return "", fmt.Errorf("cannot generate random index: %w", err)
	}

	plainIndex := hex.EncodeToString(ph[:])
	fileKey, iv, err := GenerateFileKey(cfg.Mnemonic, cfg.Bucket, plainIndex)
	if err != nil {
		return "", err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	encReader, err := EncryptReader(f, fileKey, iv)
	if err != nil {
		return "", err
	}
	sha256Hasher := sha256.New()
	sha1Hasher := sha1.New()
	r := io.TeeReader(encReader, sha256Hasher)
	r = io.TeeReader(r, sha1Hasher)
	specs := []UploadPartSpec{{Index: 0, Size: plainSize}}
	startResp, err := StartUpload(cfg, cfg.Bucket, specs)
	if err != nil {
		return "", err
	}
	part := startResp.Uploads[0]
	if err := Transfer(part, r, plainSize); err != nil {
		return "", err
	}
	encIndex := hex.EncodeToString(ph[:])
	partHash := hex.EncodeToString(sha1Hasher.Sum(nil))

	finishResp, err := FinishUpload(cfg, cfg.Bucket, encIndex, []Shard{{Hash: partHash, UUID: part.UUID}})
	if err != nil {
		return "", err
	}
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	ext := strings.TrimPrefix(filepath.Ext(base), ".")
	meta, err := CreateMetaFile(cfg, name, cfg.Bucket, finishResp.ID, "03-aes", cfg.RootFolderID, name, ext, plainSize)
	if err != nil {
		return "", err
	}
	return meta.UUID, nil
}

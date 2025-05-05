package buckets

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/StarHack/go-internxt-drive/config"
	bip39 "github.com/tyler-smith/go-bip39"
)

const chunkSize = 16 * 1024 * 1024 // 16 MiB

// GetFileDeterministicKey returns SHA512(key||data)
func GetFileDeterministicKey(key, data []byte) []byte {
	h := sha512.New()
	h.Write(key)
	h.Write(data)
	return h.Sum(nil)
}

// GenerateFileBucketKey derives a bucket-level key from mnemonic and bucketID
func GenerateFileBucketKey(mnemonic, bucketID string) ([]byte, error) {
	seed := bip39.NewSeed(mnemonic, "")
	bucketBytes, err := hex.DecodeString(bucketID)
	if err != nil {
		return nil, err
	}
	return GetFileDeterministicKey(seed, bucketBytes), nil
}

// GenerateFileKey derives the per-file key and IV from mnemonic, bucketID, and plaintext index
func GenerateFileKey(mnemonic, bucketID, indexHex string) (key, iv []byte, err error) {
	bucketKey, err := GenerateFileBucketKey(mnemonic, bucketID)
	if err != nil {
		return nil, nil, err
	}
	indexBytes, err := hex.DecodeString(indexHex)
	if err != nil {
		return nil, nil, err
	}
	detKey := GetFileDeterministicKey(bucketKey[:32], indexBytes)
	key = detKey[:32]

	iv = indexBytes[0:16]

	// debug log
	fmt.Printf(
		"Encrypting file using AES256CTR (key %s, iv %s)...\n",
		hex.EncodeToString(key),
		hex.EncodeToString(iv),
	)
	return key, iv, nil
}

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

package buckets

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
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
	return detKey[:32], detKey[32:48], nil // first 32 bytes == key, next 16 bytes == IV
}

// encryptChunkGCM reads at most chunkSize bytes from r, encrypts with AES‑GCM,
// and returns the serialized segment: nonce || ciphertext || authTag.
func encryptChunkGCM(r io.Reader, fileKey []byte) (segment []byte, err error) {
	block, err := aes.NewCipher(fileKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// read up to chunkSize
	buf := make([]byte, chunkSize)
	n, err := io.ReadFull(r, buf)
	if err == io.ErrUnexpectedEOF || err == io.EOF {
		if n == 0 {
			return nil, io.EOF
		}
		buf = buf[:n]
	} else if err != nil {
		return nil, err
	}
	// generate random 12‑byte nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	// Seal: ciphertext||tag
	ct := gcm.Seal(nil, nonce, buf, nil)
	// prepend nonce
	segment = append(nonce, ct...)
	return segment, nil
}

// encryptChunksCTR reads the file in chunkSize blocks, and for each block
// streams it through AES‑CTR with the given key and iv. It returns
// slices of ciphertext segments.
func encryptChunksCTR(f *os.File, fileKey, iv []byte) (segments [][]byte, specs []UploadPartSpec, hashes []Shard, err error) {
	block, err := aes.NewCipher(fileKey)
	if err != nil {
		return
	}
	stream := cipher.NewCTR(block, iv)

	idx := 0
	for {
		buf := make([]byte, chunkSize)
		n, readErr := io.ReadFull(f, buf)
		if readErr == io.ErrUnexpectedEOF || readErr == io.EOF {
			if n == 0 {
				break
			}
			buf = buf[:n]
		} else if readErr != nil {
			err = readErr
			return
		}

		// encrypt in‑place
		ct := make([]byte, len(buf))
		stream.XORKeyStream(ct, buf)

		// record
		segments = append(segments, ct)
		specs = append(specs, UploadPartSpec{Index: idx, Size: int64(len(ct))})
		h := sha1.Sum(ct)
		hashes = append(hashes, Shard{Hash: hex.EncodeToString(h[:]), UUID: ""})

		idx++
	}
	return
}

// ----------------------------------------------------------------
// updated UploadFile
// ----------------------------------------------------------------

func UploadFile(cfg *config.Config, filePath string) (string, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	plainSize := int64(len(raw))
	bucketID := cfg.Bucket

	// 1) plaintext index for key derivation
	ph := sha256.Sum256(raw)
	plainIndex := hex.EncodeToString(ph[:])

	// 2) derive fileKey + iv
	fileKey, iv, err := GenerateFileKey(cfg.Mnemonic, bucketID, plainIndex)
	if err != nil {
		return "", err
	}

	// 3) open file and produce encrypted chunks + metadata
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	segments, specs, shards, err := encryptChunksCTR(f, fileKey, iv)
	if err != nil {
		return "", err
	}

	// 4) overall encrypted index = SHA256(concat all ciphertext)
	h256 := sha256.New()
	for _, seg := range segments {
		h256.Write(seg)
	}
	encIndex := hex.EncodeToString(h256.Sum(nil))

	// 5) start → transfer → finish
	startResp, err := StartUpload(cfg, bucketID, specs)
	if err != nil {
		return "", err
	}

	for i, seg := range segments {
		shards[i].UUID = startResp.Uploads[i].UUID
		if err := Transfer(startResp.Uploads[i], bytes.NewReader(seg), specs[i].Size); err != nil {
			return "", err
		}
	}

	finishResp, err := FinishUpload(cfg, bucketID, encIndex, shards)
	if err != nil {
		return "", err
	}

	// 6) create the meta file with plaintext size
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	ext := strings.TrimPrefix(filepath.Ext(base), ".")

	meta, err := CreateMetaFile(
		cfg, name, bucketID, finishResp.ID,
		"03-aes", cfg.RootFolderID,
		name, ext, plainSize,
	)
	if err != nil {
		return "", err
	}
	return meta.UUID, nil
}

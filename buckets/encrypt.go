package buckets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/tyler-smith/go-bip39"
)

// NewAES256CTRCipher returns a cipher.Stream that performs AES‑256‑CTR encryption
// with the given 32‑byte key and 16‑byte IV, exactly like Node.js’s
// createCipheriv('aes-256-ctr', key, iv).
func NewAES256CTRCipher(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCTR(block, iv), nil
}

// EncryptReader wraps the provided src reader in a StreamReader that
// encrypts all data through AES‑256‑CTR (no padding):
//
//	source -> cipher -> …
func EncryptReader(src io.Reader, key, iv []byte) (io.Reader, error) {
	stream, err := NewAES256CTRCipher(key, iv)
	if err != nil {
		return nil, err
	}
	return cipher.StreamReader{S: stream, R: src}, nil
}

// DecryptReader wraps the provided src reader in a StreamReader that
// decrypts data encrypted with AES‑256‑CTR (no padding):
//
//	encryptedSrc -> source -> …
func DecryptReader(src io.Reader, key, iv []byte) (io.Reader, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	return cipher.StreamReader{S: stream, R: src}, nil
}

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

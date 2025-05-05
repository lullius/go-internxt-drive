package buckets

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
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

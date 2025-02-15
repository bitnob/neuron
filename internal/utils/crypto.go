package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Package utils provides utility functions for cryptographic operations,
// validation, and general helpers.

// Crypto provides cryptographic utilities for hashing, encryption,
// and password management.
type Crypto struct {
	key []byte
}

// NewCrypto creates a new crypto utility instance with the provided key.
// The key is used for encryption operations.
func NewCrypto(key []byte) *Crypto {
	return &Crypto{key: key}
}

// Hash generates a SHA-256 hash of the input string and returns it as a
// hexadecimal string.
func (c *Crypto) Hash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// HashPassword hashes a password using bcrypt with the default cost.
// Returns the hashed password as a string or an error if hashing fails.
//
// Example:
//
//	crypto := NewCrypto([]byte("your-key"))
//	hashedPassword, err := crypto.HashPassword("user-password")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Crypto) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword checks if a password matches its hash
func (c *Crypto) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Encrypt encrypts data using AES-256
func (c *Crypto) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

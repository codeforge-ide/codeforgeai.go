package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"golang.org/x/term"
)

// getSecretFile returns the path to the encrypted token file.
func getSecretFile() string {
	return filepath.Join(config.DataDir(), "github_token.enc")
}

// promptPassword prompts the user for a password (no echo).
func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(pw), err
}

// deriveKey derives a 32-byte key from the password using SHA-256.
func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// Encrypts and saves the token to disk.
func StoreGithubToken(token string, password string) error {
	key := deriveKey(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, []byte(token), nil)
	enc := base64.StdEncoding.EncodeToString(ciphertext)
	secretFile := getSecretFile()
	if err := os.MkdirAll(filepath.Dir(secretFile), 0700); err != nil {
		return err
	}
	return os.WriteFile(secretFile, []byte(enc), 0600)
}

// Loads and decrypts the token from disk.
func LoadGithubToken(password string) (string, error) {
	key := deriveKey(password)
	secretFile := getSecretFile()
	data, err := os.ReadFile(secretFile)
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aesgcm.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}
	nonce := ciphertext[:aesgcm.NonceSize()]
	enc := ciphertext[aesgcm.NonceSize():]
	plaintext, err := aesgcm.Open(nil, nonce, enc, nil)
	if err != nil {
		return "", errors.New("invalid password or corrupted token")
	}
	return string(plaintext), nil
}

// Interactive helper: prompt for password and token, then store.
func InteractiveStoreGithubToken() error {
	pw, err := promptPassword("Set a password for your GitHub token: ")
	if err != nil {
		return err
	}
	fmt.Print("Enter your GitHub token: ")
	var token string
	fmt.Scanln(&token)
	return StoreGithubToken(token, pw)
}

// Interactive helper: prompt for password, load token, and set env var.
func InteractiveLoadGithubToken() (string, error) {
	pw, err := promptPassword("Enter your GitHub token password: ")
	if err != nil {
		return "", err
	}
	token, err := LoadGithubToken(pw)
	if err != nil {
		return "", err
	}
	os.Setenv("GITHUB_TOKEN", token)
	return token, nil
}

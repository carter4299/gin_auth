package sessionmanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func generated_session_ID() ([]byte, error) {
	log.Debug("Entering GenerateSessionID() in sessionmanager / crypt.go ...")

	sessionID := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, sessionID)
	if err != nil {
		log.WithError(err).Error("Failed to generate random session ID.")
		return nil, err
	}

	log.Debug("Returning generated session ID ")
	return sessionID, err
}

func encrypt_session_ID(sessionID []byte) (string, error) {
	log.Debug("Entering EncryptSessionID() in sessionmanager / crypt.go ...")

	rawKey, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_ENC_KEY"))
	if err != nil {
		log.WithError(err).Error("Failed to decode the session key.")
		return "", err
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		log.WithError(err).Error("Failed to create a new cipher block.")
		return "", err
	}

	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.WithError(err).Error("Failed to generate a random nonce.")
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithError(err).Error("Failed to create a new GCM.")
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, sessionID, nil)

	encodedEncryptedData := base64.RawURLEncoding.EncodeToString(append(nonce, ciphertext...))

	clear_bytes(rawKey)
	clear_bytes(sessionID)
	clear_bytes(nonce)
	clear_bytes(ciphertext)

	log.Debug("Cleared Locals ...\tReturning encrypted session ID ... ")
	return encodedEncryptedData, nil
}

func decrypt_session_ID(encryptedSessionID string) ([]byte, error) {
	log.Debug("Entering DecryptSessionID() in sessionmanager / crypt.go ... ")

	ciphertext, err := base64.RawURLEncoding.DecodeString(encryptedSessionID)
	if err != nil {
		log.WithError(err).Error("Failed to decode the encrypted session ID from base64.")
		return nil, err
	}

	if len(ciphertext) < 12 {
		log.WithError(err).Error("Failed to decode the session key.")
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:12], ciphertext[12:]

	rawKey, err := base64.StdEncoding.DecodeString(os.Getenv("AUTH_ENC_KEY"))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		log.WithError(err).Error("Failed to create a new cipher block.")
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithError(err).Error("Failed to create a new GCM.")
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.WithError(err).Error("Failed to decrypt the session ID.")
		return nil, err
	}

	encryptedSessionID = ""
	clear_bytes(rawKey)
	clear_bytes(nonce)
	clear_bytes(ciphertext)

	log.Debug("Cleared Locals ...\tReturning decrypted session ID ... ")
	return plaintext, nil
}

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

func encryptFile(sourcePath, destPath string, key []byte) error {
	// Read the source file
	plaintext, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Create the cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a random IV (nonce) of the correct size (12 bytes for GCM)
	iv := make([]byte, 12) // GCM standard nonce size
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	// Create the GCM cipher mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Encrypt the data
	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)

	// Write the IV and ciphertext to the destination file
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(iv); err != nil {
		return err
	}
	if _, err := f.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func decryptFile(sourcePath, destPath string, key []byte) error {
	// Read the encrypted file
	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Extract the IV (first 12 bytes)
	if len(data) < 12 {
		return errors.New("ciphertext too short")
	}
	iv := data[:12]
	ciphertext := data[12:]

	// Create the cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create the GCM cipher mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Decrypt the data
	plaintext, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return err
	}

	// Write the decrypted data to the destination file
	return ioutil.WriteFile(destPath, plaintext, 0644)
}

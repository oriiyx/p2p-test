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

// encryptFile encrypts the source file with the given key and writes it to the destination
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

	// Create a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	// Create the GCM cipher mode
	aesgcm, err := cipher.NewGCMWithNonceSize(block)
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

// decryptFile decrypts the source file with the given key and writes it to the destination
func decryptFile(sourcePath, destPath string, key []byte) error {
	// Read the encrypted file
	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Extract the IV (first 16 bytes)
	if len(data) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

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

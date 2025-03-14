package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

// createEncryptedTorrent creates a torrent file and encrypts it for authorized users
func createEncryptedTorrent(config Config) error {
	// 1. Create a torrent file
	log.Printf("Creating torrent for file: %s", config.FilePath)
	torrentPath, err := createTorrentFile(config.FilePath, config.TrackerURL)
	if err != nil {
		return fmt.Errorf("failed to create torrent file: %w", err)
	}
	log.Printf("Created torrent file: %s", torrentPath)

	// 2. Generate a key for encryption if not provided
	key := make([]byte, 32)
	if config.KeyHex != "" {
		// Use provided key
		key, err = hex.DecodeString(config.KeyHex)
		if err != nil {
			return fmt.Errorf("invalid key hex: %w", err)
		}
	} else {
		// Generate a random key
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}
	}

	// Store the key for our "server"
	userID := "user1"
	if config.UserID != "" {
		userID = config.UserID
	}
	AuthorizedUsers[userID] = UserKey{
		UserID: userID,
		Key:    key,
	}

	// 3. Encrypt the torrent file
	encryptedPath := torrentPath + ".encrypted"
	if err := encryptFile(torrentPath, encryptedPath, key); err != nil {
		return fmt.Errorf("failed to encrypt torrent: %w", err)
	}
	log.Printf("Encrypted torrent file: %s", encryptedPath)
	log.Printf("User: %s", userID)
	log.Printf("Key: %s", hex.EncodeToString(key))

	return nil
}

// createTorrentFile creates a .torrent file from the given file path
func createTorrentFile(filePath, trackerURL string) (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	// Create a MetaInfo builder
	mi := metainfo.MetaInfo{
		CreatedBy:    "P2P Private File Sharing",
		CreationDate: time.Now().Unix(),
	}

	var private = true
	// Add the file to the metainfo
	info := metainfo.Info{
		PieceLength: 256 * 1024, // 256 KiB pieces
		Private:     &private,
	}

	// Handle single file vs directory differently
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}

	if fileInfo.IsDir() {
		// Adding a directory
		info.Name = filepath.Base(absPath)
		err = filepath.Walk(absPath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(absPath, path)
			if err != nil {
				return err
			}
			info.Files = append(info.Files, metainfo.FileInfo{
				Path:   []string{relPath},
				Length: fi.Size(),
			})
			return nil
		})
		if err != nil {
			return "", err
		}
	} else {
		// Adding a single file
		info.Name = filepath.Base(absPath)
		info.Length = fileInfo.Size()
	}

	// Generate pieces
	err = info.GeneratePieces(func(fi metainfo.FileInfo) (io.ReadCloser, error) {
		if len(fi.Path) == 0 {
			// Single file mode
			return os.Open(absPath)
		}
		// Multi file mode
		return os.Open(filepath.Join(absPath, filepath.Join(fi.Path...)))
	})
	if err != nil {
		return "", err
	}

	// Set the info in the metainfo
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		return "", err
	}

	// Add tracker
	if trackerURL != "" {
		mi.Announce = trackerURL
	}

	// Write the torrent file
	torrentPath := absPath + ".torrent"
	torrentFile, err := os.Create(torrentPath)
	if err != nil {
		return "", err
	}
	defer torrentFile.Close()

	if err := mi.Write(torrentFile); err != nil {
		return "", err
	}

	return torrentPath, nil
}

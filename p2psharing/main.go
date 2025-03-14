package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

// UserKey represents an access key for a specific user
type UserKey struct {
	UserID string
	Key    []byte // Encryption key for this user
}

// Config holds application configuration
type Config struct {
	Mode        string // "create", "seed", or "download"
	FilePath    string // Path to the file to share or download
	TorrentPath string // Path to the torrent file
	OutputDir   string // Directory to download files to
	UserID      string // User ID for decryption
	KeyHex      string // Hex-encoded encryption/decryption key
	TrackerURL  string // URL of the tracker
	ListenAddr  string // Address to listen on
}

// AuthorizedUsers Global map of authorized users and their keys
var AuthorizedUsers = map[string]UserKey{}

func main() {
	// Parse command-line arguments
	config := parseArgs()

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch config.Mode {
	case "create":
		// Create torrent and encrypt for authorized users
		if err := createEncryptedTorrent(config); err != nil {
			log.Fatalf("Failed to create encrypted torrent: %v", err)
		}
	case "seed":
		// Seed a file using its torrent
		if err := seedFile(config); err != nil {
			log.Fatalf("Failed to seed file: %v", err)
		}
	case "download":
		// Download a file using an encrypted torrent
		if err := downloadFile(config); err != nil {
			log.Fatalf("Failed to download file: %v", err)
		}
	default:
		log.Fatalf("Unknown mode: %s", config.Mode)
	}
}

// parseArgs parses command-line arguments and returns a Config
func parseArgs() Config {
	mode := flag.String("mode", "", "Mode: create, seed, or download")
	filePath := flag.String("file", "", "Path to the file to share or download")
	torrentPath := flag.String("torrent", "", "Path to the torrent file")
	outputDir := flag.String("output", ".", "Directory to download files to")
	userID := flag.String("user", "", "User ID for decryption")
	keyHex := flag.String("key", "", "Hex-encoded encryption/decryption key")
	trackerURL := flag.String("tracker", "udp://tracker.opentrackr.org:1337/announce", "URL of the tracker")
	listenAddr := flag.String("listen", ":50007", "Address to listen on")

	flag.Parse()

	if *mode == "" {
		log.Fatal("Mode is required: -mode=create|seed|download")
	}

	// Validate arguments based on mode
	switch *mode {
	case "create":
		if *filePath == "" {
			log.Fatal("File path is required for create mode: -file=<path>")
		}
	case "seed":
		if *filePath == "" || *torrentPath == "" {
			log.Fatal("File path and torrent path are required for seed mode: -file=<path> -torrent=<path>")
		}
	case "download":
		if *torrentPath == "" || *userID == "" || *keyHex == "" {
			log.Fatal("Torrent path, user ID, and key are required for download mode: -torrent=<path> -user=<id> -key=<key>")
		}
	}

	return Config{
		Mode:        *mode,
		FilePath:    *filePath,
		TorrentPath: *torrentPath,
		OutputDir:   *outputDir,
		UserID:      *userID,
		KeyHex:      *keyHex,
		TrackerURL:  *trackerURL,
		ListenAddr:  *listenAddr,
	}
}

// downloadFile downloads a file using an encrypted torrent
func downloadFile(config Config) error {
	// 1. Decrypt the torrent file
	log.Printf("Decrypting torrent file: %s", config.TorrentPath)

	// Decode the key from hex
	key, err := hex.DecodeString(config.KeyHex)
	if err != nil {
		return fmt.Errorf("invalid key hex: %w", err)
	}

	// Decrypt the torrent file
	decryptedPath := config.TorrentPath + ".decrypted"
	if err := decryptFile(config.TorrentPath, decryptedPath, key); err != nil {
		return fmt.Errorf("failed to decrypt torrent: %w", err)
	}
	log.Printf("Decrypted torrent to: %s", decryptedPath)

	// 2. Create output directory if it doesn't exist
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 3. Create a torrent client
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = config.OutputDir
	cfg.SetListenAddr(config.ListenAddr)

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	// 4. Add the torrent
	t, err := client.AddTorrentFromFile(decryptedPath)
	if err != nil {
		return err
	}

	// 5. Wait for the torrent to be ready
	<-t.GotInfo()

	// Print torrent info
	log.Printf("Torrent info: %s", t.Name())
	log.Printf("Size: %d bytes", t.Length())
	log.Printf("Downloading to: %s", config.OutputDir)

	// 6. Start downloading
	t.DownloadAll()

	// 7. Wait for download to complete
	log.Println("Downloading... ")

	for {
		stats := t.Stats()
		progress := float64(stats.BytesWritten.Int64()) / float64(t.Length()) * 100
		log.Printf("Progress: %.2f%% (%d/%d bytes)",
			progress, stats.BytesWritten.Int64(), t.Length())

		if t.Complete().Bool() {
			break
		}

		time.Sleep(2 * time.Second)
	}

	log.Println("Download complete!")

	// 8. Keep client running for a moment to allow for proper completion
	time.Sleep(3 * time.Second)

	return nil
}

package main

import (
	"flag"
	"log"
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
	SeedAddr    string // Address of the seeder (for download mode)
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
	listenAddr := flag.String("listen", "0.0.0.0:50007", "Address to listen on")
	seedAddr := flag.String("seed", "127.0.0.1:50007", "Address of the seeder (for download mode)")

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
		SeedAddr:    *seedAddr,
	}
}

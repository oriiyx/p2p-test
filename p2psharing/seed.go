package main

import (
	"log"
	"path/filepath"

	"github.com/anacrolix/torrent"
)

// seedFile starts seeding the given file using the torrent file
func seedFile(config Config) error {
	log.Printf("Seeding file: %s using torrent: %s", config.FilePath, config.TorrentPath)

	// Get the absolute file path
	absFilePath, err := filepath.Abs(config.FilePath)
	if err != nil {
		return err
	}

	// Get the parent directory
	parentDir := filepath.Dir(absFilePath)

	// Create a torrent client config
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = parentDir
	cfg.SetListenAddr(config.ListenAddr)
	cfg.DisableIPv6 = true // Add this line if you're having IPv6-related issues
	cfg.Seed = true        // Add this for the seeder

	// Create a client
	client, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	// Add the torrent
	t, err := client.AddTorrentFromFile(config.TorrentPath)
	if err != nil {
		return err
	}

	// Wait for the torrent to be ready
	<-t.GotInfo()

	// Print torrent info
	log.Printf("Torrent info: %s", t.Name())
	log.Printf("Size: %d bytes", t.Length())
	log.Printf("Seeding on %s", config.ListenAddr)

	// Wait indefinitely while seeding
	t.DownloadAll()
	log.Println("Seeding... Press Ctrl+C to stop")
	select {}
}

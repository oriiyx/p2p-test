package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

type InternalAddress string

func (ia *InternalAddress) String() string {
	return string(*ia)
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
	// cfg.Debug = true
	cfg.SetListenAddr(config.ListenAddr)
	cfg.DisableIPv6 = true

	cfg.DisablePEX = false
	cfg.NoDHT = false
	// cfg.ListenPort = 0 // Let the system choose a free port

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

	addr := InternalAddress(config.SeedAddr)
	t.AddPeers([]torrent.PeerInfo{
		{
			Addr: &addr, // Need to add this to Config struct
		},
	})

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
		log.Printf("total peers: %v", stats.TotalPeers)
		log.Printf("active peers: %v", stats.ActivePeers)
		log.Printf("bytes read: %v", stats.BytesRead.Int64())
		log.Printf("Connected peers: %d", stats.ActivePeers)

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

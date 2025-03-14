# Private P2P File Sharing Implementation

This document explains the implementation of a private peer-to-peer file sharing system based on the BitTorrent protocol. This system allows for sharing large files or multiple files while restricting access to authorized users through encryption.

## Overview

The system implements a privacy layer on top of the BitTorrent protocol by encrypting the torrent file itself, rather than modifying the BitTorrent protocol. Only authorized users with the correct decryption keys can access the torrent file and therefore participate in the file sharing.

### Key Components

1. **Torrent Creation**: Generates standard .torrent files from source files or directories
2. **Encryption Layer**: AES-GCM encryption of the torrent file
3. **Key Management**: Simple user-to-key mapping system
4. **Seeding**: Standard BitTorrent seeding functionality
5. **Downloading**: Decryption of the torrent file followed by standard BitTorrent downloading

## Implementation Details

### 1. Data Structures

- **UserKey**: Represents an authorized user and their encryption key
- **Config**: Holds application configuration parameters

### 2. Core Functionality

#### Torrent Creation

The `createTorrentFile` function handles:
- Analyzing the source file or directory
- Setting piece length (256KB by default)
- Generating piece hashes
- Setting tracker information
- Writing the .torrent file

#### Encryption/Decryption

The system uses AES-GCM for strong encryption:
- `encryptFile`: Encrypts a file using a provided key
- `decryptFile`: Decrypts a file using a provided key

Key features:
- 256-bit AES keys for strong security
- Random initialization vectors (IV) for each encryption
- Authenticated encryption to prevent tampering

#### File Seeding

The `seedFile` function:
- Creates a BitTorrent client with specified configuration
- Adds the torrent to the client
- Runs indefinitely to seed the file(s)

#### File Downloading

The `downloadFile` function:
- Decrypts the encrypted torrent file using the provided key
- Creates a BitTorrent client with specified configuration
- Adds the decrypted torrent to the client
- Downloads the file(s) to the specified output directory
- Displays progress information

## Usage

The application supports three modes:

### 1. Create Mode

Creates a torrent file and encrypts it for authorized users.

```bash
./p2p-sharing -mode=create -file=/path/to/file -tracker=udp://tracker.example.com:1337/announce
```

Options:
- `-file`: Path to the file or directory to share (required)
- `-tracker`: URL of the tracker (optional, has default)
- `-key`: Hex-encoded encryption key (optional, will generate if not provided)
- `-user`: User ID for key association (optional, defaults to "user1")

Output:
- Creates a .torrent file
- Creates an encrypted .torrent.encrypted file
- Outputs the encryption key for sharing with authorized users

### 2. Seed Mode

Seeds a file using its torrent file.

```bash
./p2p-sharing -mode=seed -file=/path/to/file -torrent=/path/to/file.torrent -listen=:50007
```

Options:
- `-file`: Path to the file or directory to seed (required)
- `-torrent`: Path to the torrent file (required)
- `-listen`: Network address to listen on (optional, defaults to :50007)

### 3. Download Mode

Downloads a file using an encrypted torrent file.

```bash
./p2p-sharing -mode=download -torrent=/path/to/file.torrent.encrypted -user=user1 -key=<hex_key> -output=/path/to/output
```

Options:
- `-torrent`: Path to the encrypted torrent file (required)
- `-user`: User ID for decryption (required)
- `-key`: Hex-encoded decryption key (required)
- `-output`: Directory to download files to (optional, defaults to current directory)
- `-listen`: Network address to listen on (optional, defaults to :50007)

## Security Considerations

1. **Key Distribution**: This prototype does not implement a secure key distribution mechanism. In a production environment, you would need to securely distribute keys to authorized users.

2. **User Authentication**: The prototype uses a simple user ID to key mapping. A production system would require proper user authentication.

3. **Torrent Privacy**: While the torrent file is encrypted, once decrypted, it's a standard torrent file. The BitTorrent protocol itself doesn't provide privacy protections, so users can still see peer IPs.

4. **Tracker Privacy**: The system uses public trackers by default. For better privacy, you should use a private tracker that requires authentication.

## Limitations and Future Improvements

1. **Key Management**: Implement a more robust key management system with key rotation and revocation.

2. **Authentication**: Add authentication to the tracker to ensure only authorized peers can participate.

3. **DHT and PEX Disabling**: Disable DHT (Distributed Hash Table) and PEX (Peer Exchange) for better privacy.

4. **Traffic Encryption**: Add an option to encrypt all BitTorrent traffic, not just the torrent file.

5. **Web Interface**: Create a web interface for easier management and monitoring.

6. **User Management**: Add proper user management with registration, authentication, and authorization.

## Testing Locally

To test the system locally, follow these steps:

### 1. Create an Encrypted Torrent

```bash
./p2p-sharing -mode=create -file=./testfile.txt
```

This will output a user ID and key. Note these down.

### 2. Start Seeding

In one terminal:
```bash
./p2p-sharing -mode=seed -file=./testfile.txt -torrent=./testfile.txt.torrent -listen=:50007
```

### 3. Download the File

In another terminal:
```bash
./p2p-sharing -mode=download -torrent=./testfile.txt.torrent.encrypted -user=user1 -key=<key_from_step_1> -output=./downloads -listen=:50008
```

The file should be downloaded to the ./downloads directory.

## Dependencies

The implementation uses the following external packages:

- `github.com/anacrolix/torrent`: BitTorrent client library
- Standard Go libraries for file handling, encryption, etc.

## Conclusion

This private P2P file sharing implementation provides a basic yet effective way to share files securely with authorized users. It leverages the efficiency of the BitTorrent protocol while adding a layer of privacy through encryption of the torrent files.

By controlling who can access the encrypted torrent files, you control who can participate in the file sharing, making it suitable for private file distribution scenarios.
# Build and Test Instructions

This document provides step-by-step instructions to test the Private P2P File Sharing implementation.

## Testing the Implementation Locally

To test the application on your local machine, we'll:
1. Create a test file
2. Create and encrypt a torrent for that file
3. Seed the file from one instance
4. Download the file from another instance

### 1. Create a Test File

```bash
echo "This is a test file for our P2P sharing application." > testfile.txt
```

For testing with larger files:
```bash
dd if=/dev/urandom of=largefile.bin bs=1M count=100
```

### 2. Create an Encrypted Torrent

```bash
./p2psharing -mode=create -file=./testfile.txt
```

This will output:
- The path to the torrent file (`testfile.txt.torrent`)
- The path to the encrypted torrent file (`testfile.txt.torrent.encrypted`)
- A user ID (e.g., `user1`)
- A hexadecimal encryption key

Take note of the user ID and key, as you'll need them for downloading.

### 3. Start Seeding the File

In one terminal window:

```bash
./p2psharing -mode=seed -file=./testfile.txt -torrent=./testfile.txt.torrent -listen=:50007
```

This will start seeding the file on port 50007. The program will continue running until you stop it with Ctrl+C.

### 4. Download the File

In another terminal window:

```bash
./p2psharing -mode=download -torrent=./testfile.txt.torrent.encrypted -user=user1 -key=YOUR_KEY_HERE -output=./downloads -listen=:50008
```

Replace `YOUR_KEY_HERE` with the key generated in step 2.

This will:
- Decrypt the torrent file
- Connect to the seeder
- Download the file to the `./downloads` directory
- Show download progress
- Exit when the download is complete

### 5. Verify the Download

```bash
cat ./downloads/testfile.txt
```

This should display the contents of the original file.

## Testing with Multiple Files

To test with a directory of files:

### 1. Create a Test Directory with Files

```bash
mkdir testdir
echo "File 1" > testdir/file1.txt
echo "File 2" > testdir/file2.txt
mkdir -p testdir/subdir
echo "File 3" > testdir/subdir/file3.txt
```

### 2. Create an Encrypted Torrent for the Directory

```bash
./p2psharing -mode=create -file=./testdir
```

### 3. Seed the Directory

```bash
./p2psharing -mode=seed -file=./testdir -torrent=./testdir.torrent -listen=:50007
```

### 4. Download the Directory

```bash
./p2psharing -mode=download -torrent=./testdir.torrent.encrypted -user=user1 -key=YOUR_KEY_HERE -output=./downloads -listen=:50008
```

## Testing on Multiple Machines

To test between different machines on the same network:

1. On machine A (seeder):
   ```bash
   ./p2psharing -mode=seed -file=./testfile.txt -torrent=./testfile.txt.torrent -listen=:50007
   ```

2. On machine B (downloader):
    - Copy the encrypted torrent file (`testfile.txt.torrent.encrypted`) to machine B
    - Run:
      ```bash
      ./p2psharing -mode=download -torrent=./testfile.txt.torrent.encrypted -user=user1 -key=YOUR_KEY_HERE -output=./downloads -listen=:50008
      ```

Note: Ensure that port 50007 is accessible from machine B to machine A. You may need to configure firewall settings.

## Common Issues and Troubleshooting

1. **Connection Refused**:
    - Ensure the seeder is running
    - Check firewall settings
    - Verify the correct port is being used

2. **Decryption Failed**:
    - Double-check the key value
    - Ensure you're using the correct encrypted torrent file

3. **Download Stalls**:
    - The seeder might not be reachable
    - There might be network issues between peers
    - Try using a public tracker instead of the default

4. **Build Errors**:
    - Run `go mod tidy` to ensure all dependencies are correctly resolved
    - Check for Go version compatibility

## Performance Testing

For larger files or high-volume scenarios:

1. Create a large test file:
   ```bash
   dd if=/dev/urandom of=largefile.bin bs=1M count=1000
   ```

2. Create, seed, and download as before, but monitor:
    - CPU usage: `top` or `htop`
    - Memory usage: `free -m`
    - Network usage: `iftop` or `nethogs`

## Conclusion

These instructions should help you build and test the private P2P file sharing implementation. The application provides a basic yet effective way to share files securely with authorized users, leveraging the BitTorrent protocol while adding privacy through encryption.
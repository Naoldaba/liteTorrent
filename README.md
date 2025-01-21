# xorTorrent

## Overview

This project involves developing a command-line torrent client in Go, designed for efficient file downloading and uploading using the BitTorrent protocol. Leveraging Go’s powerful concurrency model and networking capabilities, the client ensures high-performance data transfers. It supports features such as peer discovery, magnet link handling, and multi-threaded downloads, making it a robust and scalable solution for decentralized file sharing.

# Torrent Client Project Features

## Key Components

### Efficient Download/Upload

- The client adeptly handles extensive file transfers, leveraging Go's concurrency features for simultaneous downloads and uploads.

### Parser Accuracy

- Torrent files are parsed with precision, extracting vital information such as tracker URL, file name, size, and pieces.

### Robust Tracker Communication

- Establishes reliable communication with the tracker, ensuring up-to-date information on available peers.

### Piece-Level Management

- Effective piece management guarantees a systematic download and upload process, ensuring file completeness.

### Concurrency for Scalability

- Utilizes Go's concurrency features to seamlessly handle multiple downloads and uploads concurrently.

### Error Resilience

- Robust error handling mechanisms ensure the client gracefully manages unexpected scenarios, maintaining stability.

### Fault Tolerance and Recovery

- The client exhibits resilience, recovering from errors and persisting bitfield information to ensure continuity in processes.


## Setup

1. **Initialize Go Modules:**

   - Run the following command to download and install dependencies:
     ```
     go mod tidy
     ```

2. **Add Torrent File in `main.go`:**

   - Open the `main.go` file and include a valid `.torrent` file. For example, you can use a Debian ISO file.

3. **Run the Main Program:**
   - Execute the following command to start the GoTorrent client:
     ```
     go run main.go
     ```
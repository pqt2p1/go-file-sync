# Go File Sync 🚀

A high-performance, concurrent file synchronization tool written in Go. Watch directories for changes and automatically sync files with checksum verification.

## Features

- **Real-time File Watching** - Monitors directories for changes using OS-level notifications
- **Concurrent Processing** - Worker pool pattern for parallel file operations
- **Data Integrity** - SHA256 checksum verification ensures files are copied correctly
- **Recursive Directory Support** - Automatically watches subdirectories
- **Progress Tracking** - Real-time progress updates with atomic counters
- **Smart Debouncing** - Prevents duplicate sync operations
- **Bidirectional Sync** - Handles create, modify, and delete operations

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-file-sync.git
cd go-file-sync

# Build the binary
go build -o filesync cmd/filesync/main.go
```

## Usage

```bash
# Basic usage
./filesync watch <source-directory> <destination-directory>

# Example
./filesync watch ~/Documents/project ~/Backup/project
```

## How It Works

1. **File Watching**: Uses `fsnotify` to monitor file system events
2. **Worker Pool**: 5 concurrent workers process sync jobs
3. **Checksum Verification**: Each file is verified after copying using SHA256
4. **Debouncing**: Duplicate events within 500ms are filtered out

## Architecture

```
go-file-sync/
├── cmd/
│   └── filesync/
│       └── main.go          # Application entry point
├── internal/
│   ├── checksum/            # SHA256 checksum verification
│   │   └── checksum.go
│   ├── sync/                # Core synchronization logic
│   │   ├── job.go           # Sync job definition
│   │   ├── progress.go      # Progress tracking
│   │   ├── syncer.go        # File sync implementation
│   │   ├── worker.go        # Worker function
│   │   └── worker_pool.go   # Worker pool management
│   └── watcher/             # File system watching
│       └── watcher.go       # FSNotify wrapper
└── README.md
```

## Key Concepts Demonstrated

### Concurrency Patterns
- **Worker Pool**: Efficient resource management with controlled goroutines
- **Channel Communication**: Type-safe job distribution
- **Atomic Operations**: Thread-safe progress tracking

### File System Operations
- **Recursive Directory Walking**: `filepath.Walk` for directory traversal
- **OS Notifications**: Real-time events instead of polling
- **Path Manipulation**: Cross-platform path handling

### Data Integrity
- **SHA256 Checksums**: Verify file contents after copying
- **Atomic File Operations**: Ensure data consistency

## Example Output

```
2025/06/22 23:09:46 Adding watch: /home/user/source
2025/06/22 23:09:46 Adding watch: /home/user/source/subfolder
2025/06/22 23:09:46 Watching... Press Ctrl+C to stop
Files: 10 completed, 0 failed | 15.75 MB
Worker 2 processing: /home/user/source/document.pdf
Worker 2 completed: /home/user/source/document.pdf
Files: 11 completed, 0 failed | 18.25 MB
```

## Performance

- Processes multiple files concurrently
- Minimal CPU usage with OS-level file watching
- Efficient memory usage with streaming file operations
- Debouncing reduces unnecessary sync operations

## Limitations

- New deeply nested directories created with `mkdir -p` require watcher restart
- Limited to local file systems (no network sync yet)
- No conflict resolution for simultaneous edits

## Future Enhancements

- [ ] Network synchronization support
- [ ] Configurable worker pool size
- [ ] Include/exclude patterns
- [ ] Compression support
- [ ] Bandwidth limiting
- [ ] Web UI for monitoring

## Contributing

Feel free to open issues or submit pull requests!

## License

MIT License

## Acknowledgments

Built as a learning project to explore Go's concurrency patterns and file system operations. Special thanks to the Go community for excellent libraries like `fsnotify`.
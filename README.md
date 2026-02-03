# Scheduler

A lightweight, cross-platform CLI tool for executing commands at fixed time intervals aligned to the clock. Built with Go for reliability and ease of deployment.

## Overview

Scheduler runs commands at precise intervals (e.g., every 5, 10, 15, 30, 60 minutes) aligned to the clock boundaries. For example, with a 15-minute interval, commands execute at :00, :15, :30, and :45 past each hour—not 15 minutes after the scheduler starts.

## Features

- ⏰ **Clock-Aligned Scheduling**: Commands execute at exact time boundaries, not relative to start time
- 🖥️ **Cross-Platform**: Runs on macOS, Linux, and Windows
- ⏱️ **Real-Time Countdown**: Visual countdown display showing time until next execution
- 📊 **Execution Tracking**: Displays start time, end time, and duration for each run
- 🛡️ **Graceful Shutdown**: Handles Ctrl+C (SIGINT) and SIGTERM signals cleanly
- 🔄 **Continuous Operation**: Continues executing even if a command fails
- 📝 **Output Control**: Optional stdout/stderr display from executed commands

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/rubendguez/scheduler.git
cd scheduler

# Build for your platform
go build -o scheduler main.go
```

### Pre-built Binaries

Use the build scripts to create binaries for all platforms:

**macOS/Linux:**
```bash
chmod +x build.sh
./build.sh
```

**Windows:**
```bat
build.bat
```

Binaries will be generated in the `bin/` directory:
- `scheduler-darwin-amd64` - macOS Intel
- `scheduler-darwin-arm64` - macOS Apple Silicon
- `scheduler-linux-amd64` - Linux
- `scheduler-windows.exe` - Windows

## Usage

### Basic Syntax

```bash
scheduler -cmd "your-command" -interval-min <minutes> [options]
```

### Required Flags

- `-cmd` - Command to execute (will be run through the OS shell)
- `-interval-min` - Interval in minutes for clock-aligned execution (must be > 0)

### Optional Flags

- `-stdout` - Print stdout from command execution (default: false)

### Examples

**Run a backup script every 30 minutes:**
```bash
./scheduler -cmd "./backup.sh" -interval-min 30
```

**Execute a Python script every 15 minutes with stdout output:**
```bash
./scheduler -cmd "python3 data_sync.py" -interval-min 15 -stdout
```

**Run a database cleanup every hour:**
```bash
./scheduler -cmd "psql -U user -d mydb -f cleanup.sql" -interval-min 60
```

**Execute multiple commands (shell syntax):**
```bash
./scheduler -cmd "cd /app && npm run sync" -interval-min 10
```

## How It Works

1. **Calculates Next Boundary**: On startup, the scheduler calculates the next clock-aligned time boundary based on your interval
2. **Countdown Display**: Shows a real-time countdown to the next execution
3. **Executes Command**: Runs your command through the system shell at the precise boundary
4. **Reports Results**: Displays execution time, duration, and any errors
5. **Repeats**: Immediately calculates the next boundary and continues

### Clock Alignment Example

If you start the scheduler at 10:07 with a 15-minute interval:
- Next execution: 10:15
- Following executions: 10:30, 10:45, 11:00, 11:15, etc.

## Output Format

During countdown:
```
Next run at 3:30PM (in 02:15)
```

After execution:
```
[started: 3:30PM] [ended: 3:30PM][duration: 245ms]
```

With `-stdout` flag:
```
STDOUT:
<command output>
```

Errors (if any):
```
STDERR:
<error output>

Command failed (will continue execution):
<error details>
```

## Signal Handling

The scheduler gracefully handles shutdown signals:
- `Ctrl+C` (SIGINT)
- SIGTERM

On receiving a shutdown signal, it displays:
```
Shutting down gracefully.
```

## Requirements

- Go 1.25.0 or higher (for building from source)
- No external dependencies

## Platform-Specific Notes

### macOS/Linux
Commands are executed using `sh -c "your-command"`

### Windows
Commands are executed using `cmd.exe /C "your-command"`

## Use Cases

- **Data Synchronization**: Sync data between systems at regular intervals
- **Backup Operations**: Schedule regular backups aligned to specific times
- **Health Checks**: Monitor services or endpoints periodically
- **Report Generation**: Generate reports at fixed intervals
- **Cache Warming**: Pre-populate caches at predictable times
- **Log Rotation**: Trigger log rotation scripts at specific intervals
- **API Polling**: Call external APIs at regular, clock-aligned intervals

## Development

### Project Structure

```
scheduler/
├── main.go           # Main application code
├── go.mod            # Go module definition
├── build.sh          # Unix build script
├── build.bat         # Windows build script
├── .gitignore        # Git ignore rules
└── bin/              # Compiled binaries (generated)
```

### Building

```bash
# Build for current platform
go build -o scheduler main.go

# Build for all platforms
./build.sh  # macOS/Linux
build.bat   # Windows
```

### Running Tests

```bash
go test ./...
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

Created by Argenis Dominguez ([@rubendguez](https://github.com/rubendguez))

## Version

Current version: v0.0.0

---

**Note**: This scheduler is designed for reliability and simplicity. For complex cron-like scheduling with multiple patterns, consider using system cron or a dedicated job scheduler.

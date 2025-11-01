# GoVC — MKV to MP4 Batch Converter

**Professional tool in Go with Hexagonal Architecture to batch convert MKV files to MP4, with parallelism and subtitle support.**

---

## 🚀 Quick Start

### Build

```bash
go build -o govc ./cmd/govc
```

### Run

```bash
# Show help
./govc --help
./govc convert --help

# Convert with default settings (uses system CPU count)
./govc convert /path/to/videos

# Convert with 4 parallel workers
./govc convert -p 4 /path/to/videos

# Without saving temporary logs
./govc convert -p 4 --logs=false /path/to/videos

# With go run
go run ./cmd/govc convert -p 4 /path/to/videos
```

---

## 📋 `convert` Command Behaviors

### Flags

**`-p, --workers N`** (Integer)

- Controls how many `ffmpeg` processes run simultaneously
- **Default**: number of CPUs on the machine
- Recommended: `-p 2` on machines with few cores, `-p 4` on modern machines

**`--logs BOOLEAN`** (Boolean)

- **Default**: `true` - keeps logs in `mp4/<name>.log` for each converted video
- **With `--logs=false`**: temporary logs are removed after successful conversion
- **Note**: Error logs are always kept for diagnostics (regardless of flag)
- Logs contain ffmpeg stderr output for troubleshooting

### Features

**Subtitle Support**

Supports **two types**:

1. **Embedded**: internal subtitles in MKV (mapped automatically)
2. **External**: `.srt` file with same name as video
   - Example: `video.mkv` → looks for `video.srt`
   - Both are converted to `mov_text` (MP4 compatible)

⚠️ **Limitation**: Image-based subtitles (PGS, VobSub) are not supported (would require re-encoding)

---

## 🏛️ Architecture

This project follows **Hexagonal Architecture (Ports & Adapters)** for maximum testability and extensibility:

```
┌─────────────────────────────────────────┐
│         CLI Input Adapter               │
└────────────────────┬────────────────────┘
                     │
                     ↓
       ┌─────────────────────────┐
       │  ConversionService      │ ← Use Case (core)
       │  (Orchestration)        │
       └────────┬──────┬─────────┘
                │      │
        ┌───────┘      └────────┐
        ↓                       ↓
   Filesystem           FFmpeg Adapter
   (Video Discovery)    (Converter)
```

**Benefits**:

- ✅ **Testable**: Mock adapters without running real ffmpeg
- ✅ **Decoupled**: Replace ffmpeg with another tool? New adapter, done
- ✅ **Scalable**: Add REST API? New input adapter
- ✅ **Readable**: Structure reflects business domain

---

## 📁 Code Structure

```
GoVC/
├── cmd/govc/main.go                    ← Entry point (Bootstrap)
├── internal/core/
│   ├── domain/                         ← Pure entities (Video, Conversion)
│   ├── ports/                          ← Interfaces (contracts)
│   └── services/conversion_service.go  ← Use case
└── internal/adapters/
    ├── cli/                            ← Input: CLI adapter
    ├── commands/                       ← urfave/cli command handlers
    ├── filesystem/                     ← Output: File system
    └── ffmpeg/                         ← Output: Converter
```

For detailed documentation → see `HEXAGONAL_ARCHITECTURE.md`

---

## 📦 Requirements

- **Go 1.20+**
- **ffmpeg** and **ffprobe** installed and in PATH

```bash
# macOS
brew install ffmpeg

# Linux (Ubuntu/Debian)
sudo apt-get install ffmpeg

# Windows (with Chocolatey)
choco install ffmpeg
```

---

## 💡 Usage Examples

### Example 1: Simple Conversion

```bash
./govc convert /videos
```

- Uses all CPUs
- Saves logs in `/videos/mp4/`

### Example 2: Control Workers and Logs

```bash
./govc convert -p 2 --logs=false /videos
```

- 2 parallel workers
- Deletes successful logs after conversion
- Error logs are still kept

### Example 3: Release Build

```bash
go build -o govc-v1.0 ./cmd/govc
./govc-v1.0 convert -p 4 /media/movies
```

---

## 🧪 Testing

To test the domain in isolation (without ffmpeg):

```go
package domain

import "testing"

func TestProgressTracker(t *testing.T) {
    tracker := domain.NewProgressTracker(3)
    tracker.Update("video1", 50)
    if tracker.GetSnapshot()["video1"] != 50 {
        t.Fatal("Progress not updated")
    }
}
```

All adapters can be mocked for pure unit tests.

---

## 🔧 Extending

Adding a new adapter is simple. Example: support for HTTP API to submit conversions:

1. Create `internal/adapters/http/adapter.go`
2. Implement `ConfigPort` to read config via HTTP
3. Inject into `ConversionService` (in bootstrap)

For complete guide → see `EXTENSION_GUIDE.md`

---

## 📝 Changelog

- **v1.0** (current): Refactoring with Hexagonal Architecture
  - ✅ Pure core domain
  - ✅ Well-defined ports
  - ✅ Decoupled adapters
  - ✅ Testable and extensible

---

## 📖 References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Developed with focus on Clean Code and Quality Architecture.**

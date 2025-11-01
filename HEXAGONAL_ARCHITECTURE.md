# GoVC — Hexagonal Architecture Refactoring

## 📐 Hexagonal Architecture (Ports & Adapters)

This is the reorganization of the GoVC project applying **Hexagonal Architecture**, making the code **testable, loosely coupled, and ready for evolution**.

### 🔷 Folder Structure

```
GoVC/
├── cmd/
│   └── govc/
│       └── main.go              ← Entry point (Bootstrap) - minimal and clean
├── internal/                    ← Private code (not exportable)
│   ├── core/                    ← Heart: pure business logic
│   │   ├── domain/              ← Entities (Video, Conversion, Progress)
│   │   │   ├── video.go
│   │   │   ├── conversion.go
│   │   │   └── progress.go
│   │   ├── ports/               ← Interfaces (abstract contracts) - ONE PER FILE
│   │   │   ├── config.go
│   │   │   ├── executor.go
│   │   │   ├── file_system.go
│   │   │   ├── progress_reporter.go
│   │   │   ├── service_command.go
│   │   │   ├── video_converter.go
│   │   │   ├── video_discovery.go
│   │   │   └── command_executor.go
│   │   └── services/            ← Use Cases (application)
│   │       └── conversion_service.go
│   └── adapters/                ← Concrete implementations
│       ├── cli/                 ← Input: CLI arguments
│       │   ├── config.go        ← Implements ConfigPort
│       │   ├── logger.go        ← Implements ProgressReporterPort
│       │   ├── command_executor.go       ← Implements CommandExecutorPort
│       │   ├── convert_command.go        ← Implements ServiceCommand
│       │   ├── *_mock.go        ← Mocks for testing (one per adapter)
│       ├── commands/            ← Command Handlers (urfave/cli integration)
│       │   ├── convert.go       ← ConvertCommandHandler (handles convert command)
│       │   └── factory.go       ← CommandFactory (builds all CLI commands)
│       ├── filesystem/          ← Output: File system operations
│       │   ├── adapter.go       ← Implements VideoDiscoveryPort, FileSystemPort
│       │   └── adapter_mock.go  ← Mock for testing
│       └── ffmpeg/              ← Output: Converter tool
│           ├── adapter.go       ← Implements VideoConverterPort
│           └── adapter_mock.go  ← Mock for testing
├── main.go                      ← Stub (points to cmd/govc)
├── go.mod
└── README.md
```

---

### 🎯 Key Concepts

#### 1️⃣ **Core (Heart)**

- **Domain**: Pure entities, no external dependencies
  - `Video`, `ConversionResult`, `ProgressTracker`
- **Ports**: Interfaces that define contracts (one interface per file)
  - `ConfigPort`, `Executor`, `FileSystemPort`, `ProgressReporterPort`, `ServiceCommand`, `VideoConverterPort`, `VideoDiscoveryPort`, `CommandExecutorPort`
- **Services**: Use Cases that orchestrate domain + ports
  - `ConversionService` (does orchestration, nothing more)

#### 2️⃣ **Adapters (Sides — Inputs)**

- **CLI** (`internal/adapters/cli/`): Reads user arguments and configurations

  - `CLIConfig` → implements `ConfigPort`
  - `LoggerReporter` → implements `ProgressReporterPort`
  - `CommandExecutor` → implements `CommandExecutorPort` (orchestrates commands)
  - `ConvertCommand` → implements `ServiceCommand` (wraps ConversionService)

- **Commands** (`internal/adapters/commands/`): urfave/cli command handlers
  - `ConvertCommandHandler` → handles the convert command execution
  - `CommandFactory` → builds all available CLI commands (easy to add new commands)

#### 3️⃣ **Adapters (Sides — Outputs)**

- **Filesystem** (`internal/adapters/filesystem/`): Interacts with FS
  - `FilesystemAdapter` → implements `VideoDiscoveryPort` + `FileSystemPort`
- **FFmpeg** (`internal/adapters/ffmpeg/`): Interacts with ffmpeg
  - `FFmpegAdapter` → implements `VideoConverterPort`

#### 4️⃣ **Bootstrap** (`cmd/govc/main.go`)

- Creates adapter instances
- Injects into `ConversionService`
- Executes

---

### ✅ Advantages of This Architecture

| Advantage                       | Explanation                                    |
| ------------------------------- | ---------------------------------------------- |
| **Testability**                 | Mock adapters easily; domain doesn't touch I/O |
| **Implementation Independence** | Replace ffmpeg? New adapter, done              |
| **Clarity**                     | Flow: CLI → Service → Adapters → Output        |
| **Maintainability**             | Change in one part doesn't break another       |
| **Scalability**                 | Add REST API? New input adapter                |

---

### 🚀 How to Run

#### Build

```bash
cd /Users/wallissonmarinho/www/GoVC
go build -o govc ./cmd/govc
```

#### Run

```bash
# Convert with default settings (uses system CPU count)
./govc convert /path/with/videos.mkv

# Convert with 2 parallel workers
./govc convert -p 2 /path/with/videos.mkv

# Convert and keep logs
./govc convert -p 4 --logs /path/with/videos.mkv

# Convert and delete successful logs
./govc convert -p 4 --logs=false /path/with/videos.mkv

# Show help
./govc --help
./govc convert --help
```

**Flags:**

- `-p, --workers N` : number of parallel workers (default: #CPUs)
- `--logs` : save per-file logs in mp4/ directory (default: true; use --logs=false to delete successful logs)

**Framework:** Uses [urfave/cli/v2](https://cli.urfave.org/) for robust CLI command handling.

---

### 📋 Execution Flow

```
┌─────────────────────────────────────────────────────────────┐
│                 urfave/cli App Bootstrap                    │
│                   cmd/govc/main.go                          │
│            (CommandFactory builds all commands)             │
└────────────────────────┬────────────────────────────────────┘
                         │ routes to command handler
                         ↓
┌─────────────────────────────────────────────────────────────┐
│         ConvertCommandHandler (adapter)                     │
│      (parses urfave/cli context, extracts args)            │
│     (creates adapters, injects into service)               │
└────────────────────────┬────────────────────────────────────┘
                         │ creates service + calls Execute
                         ↓
┌─────────────────────────────────────────────────────────────┐
│            ConversionService (Use Case)                     │
│        ┌─────────────────────────────────────┐              │
│        │ Execute() {                         │              │
│        │  - Discover videos                  │              │
│        │  - Setup progress tracking          │              │
│        │  - Orchestrate parallel conversion  │              │
│        │ }                                   │              │
│        └─────────────────────────────────────┘              │
└────────────────────────┬────────────────────────────────────┘
                         │ uses ports (interfaces)
                ┌────────┼────────┬─────────────┐
                ↓        ↓        ↓             ↓
            ┌──────┐ ┌──────┐ ┌────┐ ┌─────────────┐
            │ CLI  │ │  FS  │ │FFM │ │ ProgressRep │
            │Adapt │ │Adapt │ │Adap│ │   Adapter   │
            └──────┘ └──────┘ └────┘ └─────────────┘
               ↓        ↓        ↓             ↓
         Parse Flags Discover Convert      Report
              & Logs   Videos    & Embed   Progress
                        FS    Subtitles      Logs
```

---

### 🧪 Test an Adapter Isolated (Example)

```go
// Without depending on anything else!
adapter := ffmpeg.NewFFmpegAdapter(func(p float64) { fmt.Println(p) })
duration, _ := adapter.GetDuration("/video.mkv")
fmt.Printf("Duration: %.2f seconds\n", duration)
```

---

### 🔄 Possible Future Evolutions

1. **New Input**: Web API

   - New input adapter (`internal/adapters/http/`)
   - Register new handler in `CommandFactory`
   - Same service, new entry point

2. **New Command**: Video validation

   - New handler (`internal/adapters/commands/validate.go`)
   - Add to `CommandFactory.BuildCommands()`
   - Clean separation of concerns

3. **New Output**: AWS S3

   - New adapter (`internal/adapters/s3/`)
   - Same service, new output

4. **Unit Tests**:

   ```go
   // Mock adapter
   mockConverter := &MockConverterPort{}
   service := services.NewConversionService(..., mockConverter, ...)
   // Test without real ffmpeg!
   ```

---

### 📊 Comparison: Before vs. After

| Aspect      | Before                  | After                     |
| ----------- | ----------------------- | ------------------------- |
| Files       | 1 (`main.go` 418 lines) | 13+ files, well organized |
| Testability | Low (ffmpeg hardcoded)  | High (easy mocks)         |
| Coupling    | High (everything mixed) | Low (via interfaces)      |
| Scalability | Difficult               | Easy (new adapters)       |
| Clarity     | Confusing               | Clear (hexagonal flow)    |

---

### 📖 References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports & Adapters Pattern](https://wiki.c2.com/?HexagonalArchitecture)

---

**Ready for production with robust design!** 🎉

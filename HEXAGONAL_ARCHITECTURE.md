# GoVC — Hexagonal Architecture Refactoring

## 📐 Hexagonal Architecture (Ports & Adapters)

This is the reorganization of the GoVC project applying **Hexagonal Architecture**, making the code **testable, loosely coupled, and ready for evolution**.

### 🔷 Folder Structure

```
GoVC/
├── cmd/
│   └── govc/
│       ├── main.go              ← Entry point (Bootstrap)
│       └── ...
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
│       │   ├── config_mock.go   ← Mock for testing
│       │   ├── logger_mock.go   ← Mock for testing
│       │   ├── command_executor_mock.go  ← Mock for testing
│       │   └── convert_command_mock.go   ← Mock for testing
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

- **CLI** (`internal/adapters/cli/`): Reads user arguments and commands
  - `CLIConfig` → implements `ConfigPort`
  - `LoggerReporter` → implements `ProgressReporterPort`
  - `CommandExecutor` → implements `CommandExecutorPort` (orchestrates commands)
  - `ConvertCommand` → implements `ServiceCommand` (wraps ConversionService)

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
./govc -p 4 -logs=true /path/with/videos.mkv
```

Or directly:

```bash
go run ./cmd/govc -p 4 /path/with/videos.mkv
```

**Flags:**

- `-p N` : number of parallel workers (default: #CPUs)
- `-logs=true|false` : keep per-file logs in mp4/ directory (default: true; false deletes successful logs)

---

### 📋 Execution Flow

```
┌─────────────────────────────────────────────────────────────┐
│                      cmd/govc/main.go                       │ ← Bootstrap
│                    (Dependency Injection)                   │
└────────────────────────┬────────────────────────────────────┘
                         │ creates adapters + service
                         ↓
┌─────────────────────────────────────────────────────────────┐
│            ConversionService (Use Case)                     │
│        ┌─────────────────────────────────────┐              │
│        │ Execute() {                         │              │
│        │  - Discover videos                  │              │
│        │  - Setup progress                   │              │
│        │  - Orchestrate conversion           │              │
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
         Parse Flags Discover Videos Convert  Report
                         FS                   Logs
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
   - Same service, new entry point

2. **New Output**: AWS S3

   - New adapter (`internal/adapters/s3/`)
   - Same service, new output

3. **Unit Tests**:

   ```go
   // Mock adapter
   mockConverter := &MockConverterPort{}
   service := services.NewConversionService(..., mockConverter, ...)
   // Test without real ffmpeg!
   ```

4. **Enhanced CLI**: Use `cobra` or `urfave/cli`
   - More robust CLI adapter

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

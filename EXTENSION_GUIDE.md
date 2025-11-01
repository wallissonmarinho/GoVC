# 🔧 Extension Guide — How to Add New Adapters# 🔧 Extension Guide — How to Add New Adapters

## 📋 Common Scenarios## 📋 Common Scenarios

This guide shows how to extend the system without breaking anything!This guide shows how to extend the system without breaking anything!

---

## 🔌 Case 1: Adding a New Output Adapter (Ex: S3)## 🔌 Case 1: Adding a New Output Adapter (Ex: S3)## 🔌 Case 1: Adding a New Output Adapter (Ex: S3)

### Problem### Problem### Problem

You want to save converted videos directly to S3 instead of local filesystem.You want to save converted videos directly to S3 instead of local filesystem.You want to save converted videos directly to S3 instead of local filesystem.

### Solution### Solution### Solution

1. **Extend the interface** in `ports/ports.go` (if needed)1. **Extend the interface** in `ports/ports.go` (if needed)1. **Extend the interface** in `ports/ports.go` (if needed)

2. **Create new adapter** in `internal/adapters/s3/`

3. **Inject in bootstrap** in `cmd/govc/main.go`2. **Create new adapter** in `internal/adapters/s3/`2. **Create new adapter** in `internal/adapters/s3/`

### Step by Step3. **Inject in bootstrap** in `cmd/govc/main.go`3. **Inject in bootstrap** in `cmd/govc/main.go`

#### 1️⃣ Create the S3 Adapter (`internal/adapters/s3/adapter.go`)### Step by Step### Passo a Passo

```go#### 1️⃣ Create the S3 adapter (`internal/adapters/s3/adapter.go`)#### 1️⃣ Criar o adapter S3 (`internal/adapters/s3/adapter.go`)

package s3

`go`go

import (

    "github.com/wallissonmarinho/GoVC/internal/core/domain"package s3package s3

    "github.com/wallissonmarinho/GoVC/internal/core/ports"

    "github.com/aws/aws-sdk-go/aws"import (import (

    "github.com/aws/aws-sdk-go/aws/session"

    "github.com/aws/aws-sdk-go/service/s3/s3manager"    "github.com/wallissonmarinho/GoVC/internal/core/domain"    "github.com/wallissonmarinho/GoVC/internal/core/domain"

)

    "github.com/wallissonmarinho/GoVC/internal/core/ports"    "github.com/wallissonmarinho/GoVC/internal/core/ports"

// Ensure S3Adapter implements FileSystemPort

var \_ ports.FileSystemPort = (\*S3Adapter)(nil) "github.com/aws/aws-sdk-go/aws" "github.com/aws/aws-sdk-go/aws"

type S3Adapter struct { "github.com/aws/aws-sdk-go/aws/session" "github.com/aws/aws-sdk-go/aws/session"

    bucket   string

    uploader *s3manager.Uploader    "github.com/aws/aws-sdk-go/service/s3/s3manager"    "github.com/aws/aws-sdk-go/service/s3/s3manager"

}

))

func NewS3Adapter(bucket string) (\*S3Adapter, error) {

    sess, err := session.NewSession(&aws.Config{})// Ensure S3Adapter implements FileSystemPort// Ensure S3Adapter implements FileSystemPort

    if err != nil {

        return nil, errvar _ ports.FileSystemPort = (\*S3Adapter)(nil)var _ ports.FileSystemPort = (\*S3Adapter)(nil)

    }

type S3Adapter struct {type S3Adapter struct {

    return &S3Adapter{

        bucket:   bucket,    bucket   string    bucket string

        uploader: s3manager.NewUploader(sess),

    }, nil    uploader *s3manager.Uploader    uploader *s3manager.Uploader

}

}}

func (s \*S3Adapter) FileExists(path string) bool {

    // Check if file exists in S3func NewS3Adapter(bucket string) (*S3Adapter, error) {func NewS3Adapter(bucket string) (*S3Adapter, error) {

    return false

} sess, err := session.NewSession(&aws.Config{}) sess, err := session.NewSession(&aws.Config{})

func (s \*S3Adapter) IsValidOutput(path string) bool { if err != nil { if err != nil {

    // Validate path in S3

    return true        return nil, err        return nil, err

}

    }    }

func (s \*S3Adapter) RemoveFile(path string) error {

    // Remove file from S3

    return nil

} return &S3Adapter{ return &S3Adapter{

func (s \*S3Adapter) WriteLog(logPath string, lines []string) error { bucket: bucket, bucket: bucket,

    // Write log to S3

    return nil        uploader: s3manager.NewUploader(sess),        uploader: s3manager.NewUploader(sess),

}

`````}, nil    }, nil



#### 2️⃣ Update Bootstrap (`cmd/govc/main.go`)}}



```gofunc (s *S3Adapter) FileExists(path string) bool {func (s *S3Adapter) FileExists(path string) bool {

import (

    "log"    // Implement S3 check    // Implementar check em S3

    "github.com/wallissonmarinho/GoVC/internal/adapters/cli"

    "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"    return false    return false

    "github.com/wallissonmarinho/GoVC/internal/adapters/s3"      // ← NEW

    "github.com/wallissonmarinho/GoVC/internal/core/services"}}

)

func (s *S3Adapter) IsValidOutput(path string) bool {func (s *S3Adapter) IsValidOutput(path string) bool {

func main() {

    cliConfig, err := cli.NewCLIConfig()    // Validate in S3    // Validar em S3

    if err != nil {

        log.Fatalf("Config error: %v", err)    return true    return true

    }

}}

    converterAdapter := ffmpeg.NewFFmpegAdapter(func(p float64) {})

func (s *S3Adapter) RemoveFile(path string) error {func (s *S3Adapter) RemoveFile(path string) error {

    // Replace Filesystem with S3

    fileSystemAdapter, err := s3.NewS3Adapter("my-bucket")    // Remove from S3    // Remover em S3

    if err != nil {

        log.Fatalf("S3 error: %v", err)    return nil    return nil

    }

}}

    reporterAdapter := cli.NewLoggerReporter()

    discoveryAdapter := filesystem.NewFilesystemAdapter()func (s *S3Adapter) WriteLog(logPath string, lines []string) error {func (s *S3Adapter) WriteLog(logPath string, lines []string) error {



    service := services.NewConversionService(    // Write log to S3    // Escrever log em S3

        discoveryAdapter,

        converterAdapter,    return nil    return nil

        fileSystemAdapter,  // ← Now S3!

        reporterAdapter,}}

        cliConfig,

    )````



    if err := service.Execute(); err != nil {

        log.Fatalf("Error: %v", err)

    }#### 2️⃣ Update Bootstrap (`cmd/govc/main.go`)#### 2️⃣ Atualizar Bootstrap (`cmd/govc/main.go`)

}

`````

**Done!** Videos now save to S3, without touching anything else.`go`go

---import (import (

## 🌐 Case 2: Adding a New Input Adapter (Ex: HTTP API) "log" "log"

### Problem "github.com/wallissonmarinho/GoVC/internal/adapters/cli" "github.com/wallissonmarinho/GoVC/internal/adapters/cli"

You want the system to accept HTTP requests to convert videos. "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg" "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"

### Solution "github.com/wallissonmarinho/GoVC/internal/adapters/s3" // ← NEW "github.com/wallissonmarinho/GoVC/internal/adapters/s3" // ← NEW

1. **Create new adapter** in `internal/adapters/http/` "github.com/wallissonmarinho/GoVC/internal/core/services" "github.com/wallissonmarinho/GoVC/internal/core/services"

2. **Implement ConfigPort** and **VideoDiscoveryPort**

3. **New bootstrap in** `cmd/api/main.go`))

### Step by Step

#### 1️⃣ Create HTTP Adapter (`internal/adapters/http/config.go`)func main() {func main() {

`````go cliConfig, err := cli.NewCLIConfig()    cliConfig, err := cli.NewCLIConfig()

package http

    if err != nil {    if err != nil {

// HTTPConfig implements ConfigPort

type HTTPConfig struct {        log.Fatalf("Config error: %v", err)        log.Fatalf("Config error: %v", err)

    inputDir  string

    outputDir string    }    }

    workers   int

    saveLogs  bool

}

    converterAdapter := ffmpeg.NewFFmpegAdapter(func(p float64) {})    converterAdapter := ffmpeg.NewFFmpegAdapter(func(p float64) {})

func (h *HTTPConfig) GetInputDir() string {

    return h.inputDir

}

    // ✨ New: S3 instead of Filesystem    // ✨ Novo: S3 em vez de Filesystem

func (h *HTTPConfig) GetOutputDir() string {

    return h.outputDir    fileSystemAdapter, err := s3.NewS3Adapter("my-bucket")    fileSystemAdapter, err := s3.NewS3Adapter("my-bucket")

}

    if err != nil {    if err != nil {

func (h *HTTPConfig) GetWorkers() int {

    return h.workers        log.Fatalf("S3 error: %v", err)        log.Fatalf("S3 error: %v", err)

}

    }    }

func (h *HTTPConfig) IsSaveLogs() bool {

    return h.saveLogs

}

```    reporterAdapter := cli.NewLoggerReporter()    reporterAdapter := cli.NewLoggerReporter()



#### 2️⃣ Create HTTP Handler (`internal/adapters/http/handler.go`)    discoveryAdapter := filesystem.NewFilesystemAdapter()    discoveryAdapter := filesystem.NewFilesystemAdapter()



```go

package http

    service := services.NewConversionService(    service := services.NewConversionService(

import (

    "encoding/json"        discoveryAdapter,        discoveryAdapter,

    "net/http"

    "github.com/wallissonmarinho/GoVC/internal/core/services"        converterAdapter,        converterAdapter,

)

        fileSystemAdapter,  // ← Now S3!        fileSystemAdapter,  // ← Agora S3!

type ConversionRequest struct {

    VideoPath string `json:"video_path"`        reporterAdapter,        reporterAdapter,

    Workers   int    `json:"workers,omitempty"`

}        cliConfig,        cliConfig,



func ConvertHandler(service *services.ConversionService) http.HandlerFunc {    )    )

    return func(w http.ResponseWriter, r *http.Request) {

        var req ConversionRequest

        json.NewDecoder(r.Body).Decode(&req)

    if err := service.Execute(); err != nil {    if err := service.Execute(); err != nil {

        // Execute conversion

        err := service.Execute()        log.Fatalf("Error: %v", err)        log.Fatalf("Error: %v", err)



        w.Header().Set("Content-Type", "application/json")    }    }

        json.NewEncoder(w).Encode(map[string]interface{}{

            "success": err == nil,}}

            "error":   err,

        })````

    }

}**Done!** Videos now save to S3, without touching anything else.**Pronto!** Vídeos agora salvam em S3, sem tocar em nada else.

`````

---

#### 3️⃣ Create API Bootstrap (`cmd/api/main.go`)

## 🌐 Case 2: Adding a New Input Adapter (Ex: HTTP API)## 🌐 Caso 2: Adicionar um Novo Adapter de Input (Ex: HTTP API)

```go

package main### Problem### Problema



import (You want the system to accept HTTP requests to convert videos.Você quer que o sistema aceite requisições HTTP para converter vídeos.

    "log"

    "net/http"### Solution### Solução

    "github.com/wallissonmarinho/GoVC/internal/adapters/http"

    "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"1. **Create new adapter** in `internal/adapters/http/`1. **Criar novo adapter** em `internal/adapters/http/`

    "github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"

    "github.com/wallissonmarinho/GoVC/internal/core/services"2. **Implement ConfigPort** and **VideoDiscoveryPort**2. **Implementar ConfigPort** e **VideoDiscoveryPort**

)

3. **New bootstrap in** `cmd/api/main.go`3. **Bootstrap novo em** `cmd/api/main.go`

func main() {

    httpConfig := &http.HTTPConfig{### Step by Step### Passo a Passo

        inputDir:  "/videos",

        outputDir: "/videos/output",#### 1️⃣ Create HTTP adapter (`internal/adapters/http/config.go`)#### 1️⃣ Criar HTTP adapter (`internal/adapters/http/config.go`)

        workers:   4,

        saveLogs:  true,`go`go

    }

package httppackage http

    discoveryAdapter := filesystem.NewFilesystemAdapter()

    converterAdapter := ffmpeg.NewFFmpegAdapter(func(p float64) {})// HTTPConfig implements ConfigPort// HTTPConfig implementa ConfigPort

    fileSystemAdapter := filesystem.NewFilesystemAdapter()

    reporterAdapter := http.NewHTTPReporter()type HTTPConfig struct {type HTTPConfig struct {



    service := services.NewConversionService(    inputDir  string    inputDir  string

        discoveryAdapter,

        converterAdapter,    outputDir string    outputDir string

        fileSystemAdapter,

        reporterAdapter,    workers   int    workers   int

        httpConfig,

    )    saveLogs  bool    saveLogs  bool



    // HTTP Routes}}

    http.HandleFunc("/convert", http.ConvertHandler(service))

func (h *HTTPConfig) GetInputDir() string {func (h *HTTPConfig) GetInputDir() string {

    log.Fatal(http.ListenAndServe(":8080", nil))

}    return h.inputDir    return h.inputDir

```

}}

**Now you have:**

- CLI: `go run ./cmd/govc -p 4 /path`// ... etc// ... etc

- HTTP API: `curl -X POST http://localhost:8080/convert`

`````

---



## 🔄 Case 3: Replace FFmpeg with MediaInfo

#### 2️⃣ Create HTTP handler (`internal/adapters/http/handler.go`)#### 2️⃣ Criar HTTP handler (`internal/adapters/http/handler.go`)

### Problem



You want to use `mediainfo` instead of `ffmpeg` to extract metadata.

```go```go

### Solution

package httppackage http

**Create new adapter** `internal/adapters/mediainfo/adapter.go` implementing `VideoConverterPort`.



```go

package mediainfoimport (import (



import (    "encoding/json"    "encoding/json"

    "os/exec"

    "github.com/wallissonmarinho/GoVC/internal/core/domain"    "net/http"    "net/http"

)

    "github.com/wallissonmarinho/GoVC/internal/core/services"    "github.com/wallissonmarinho/GoVC/internal/core/services"

type MediaInfoAdapter struct{}

))

func (m *MediaInfoAdapter) Convert(video *domain.Video, inputDir string) error {

    // Use mediainfo instead of ffmpeg

    cmd := exec.Command("mediainfo", video.Path)

    return cmd.Run()type ConversionRequest struct {type ConversionRequest struct {

}

    VideoPath string `json:"video_path"`    VideoPath string `json:"video_path"`

func (m *MediaInfoAdapter) GetDuration(videoPath string) (float64, error) {

    // Parse mediainfo output    Workers   int    `json:"workers,omitempty"`    Workers   int    `json:"workers,omitempty"`

    return 0, nil

}}}

```



**Then, just replace in bootstrap:**

func ConvertHandler(service *services.ConversionService) http.HandlerFunc {func ConvertHandler(service *services.ConversionService) http.HandlerFunc {

```go

converterAdapter := mediainfo.NewMediaInfoAdapter() // ← New!    return func(w http.ResponseWriter, r *http.Request) {    return func(w http.ResponseWriter, r *http.Request) {

```

        var req ConversionRequest        var req ConversionRequest

---

        json.NewDecoder(r.Body).Decode(&req)        json.NewDecoder(r.Body).Decode(&req)

## 🗄️ Case 4: Add Database (PostgreSQL)



### New Adapter

        // Execute conversion        // Execute conversion

Create `internal/adapters/postgres/adapter.go` to track conversions in DB.

        err := service.Execute()        err := service.Execute()

```go

package postgres



import (        w.Header().Set("Content-Type", "application/json")        w.Header().Set("Content-Type", "application/json")

    "database/sql"

    _ "github.com/lib/pq"        json.NewEncoder(w).Encode(map[string]interface{}{        json.NewEncoder(w).Encode(map[string]interface{}{

)

            "success": err == nil,            "success": err == nil,

type PostgresAdapter struct {

    db *sql.DB            "error":   err,            "error":   err,

}

        })        })

func NewPostgresAdapter(connStr string) (*PostgresAdapter, error) {

    db, err := sql.Open("postgres", connStr)    }    }

    if err != nil {

        return nil, err}}

    }

    return &PostgresAdapter{db: db}, nil````

}

#### 3️⃣ Create API bootstrap (`cmd/api/main.go`)#### 3️⃣ Criar bootstrap API (`cmd/api/main.go`)

func (p *PostgresAdapter) LogConversion(filename string, success bool) error {

    _, err := p.db.Exec(`go`go

        "INSERT INTO conversions (filename, success, timestamp) VALUES ($1, $2, NOW())",

        filename, success,package mainpackage main

    )

    return errimport (import (

}

```    "log"    "log"



**Use in service** (can extend `ProgressReporterPort` or create new port).    "net/http"    "net/http"



---    "github.com/wallissonmarinho/GoVC/internal/adapters/http"    "github.com/wallissonmarinho/GoVC/internal/adapters/http"



## 📊 Quick Reference    "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"    "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"



### Add New Adapter (Checklist)    "github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"    "github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"



- [ ] Create folder: `internal/adapters/new_adapter/`    "github.com/wallissonmarinho/GoVC/internal/core/services"    "github.com/wallissonmarinho/GoVC/internal/core/services"

- [ ] Create file: `adapter.go`

- [ ] Implement relevant interface(s) from `ports/ports.go`))

- [ ] Add assertion: `var _ ports.YourPort = (*YourAdapter)(nil)`

- [ ] Update bootstrap in `cmd/govc/main.go`func main() {func main() {

- [ ] Build & test: `go build ./cmd/govc`

- [ ] Document usage in your README    httpConfig := &http.HTTPConfig{    httpConfig := &http.HTTPConfig{



---        inputDir:  "/videos",        inputDir:  "/videos",



## 🧪 Test New Adapter (Mock)        outputDir: "/videos/output",        outputDir: "/videos/output",



```go        workers:   4,        workers:   4,

package main

        saveLogs:  true,        saveLogs:  true,

import (

    "testing"    }    }

    "github.com/wallissonmarinho/GoVC/internal/adapters/your_adapter"

)



func TestYourAdapter(t *testing.T) {    discoveryAdapter := filesystem.NewFilesystemAdapter()    discoveryAdapter := filesystem.NewFilesystemAdapter()

    adapter := your_adapter.New()

    converterAdapter := ffmpeg.NewFFmpegAdapter(func(p float64) {})    converterAdapter := ffmpeg.NewFFmpegAdapter(func(p float64) {})

    // Test method

    result, err := adapter.YourMethod()    fileSystemAdapter := filesystem.NewFilesystemAdapter()    fileSystemAdapter := filesystem.NewFilesystemAdapter()



    if err != nil {    reporterAdapter := http.NewHTTPReporter()    reporterAdapter := http.NewHTTPReporter()

        t.Fatalf("Error: %v", err)

    }



    if result != expected {    service := services.NewConversionService(    service := services.NewConversionService(

        t.Errorf("got %v, want %v", result, expected)

    }        discoveryAdapter,        discoveryAdapter,

}

```        converterAdapter,        converterAdapter,



---        fileSystemAdapter,        fileSystemAdapter,



## 📈 Escalation Path        reporterAdapter,        reporterAdapter,



```        httpConfig,        httpConfig,

Phase 1: CLI (current) ✅

└─ cmd/govc/main.go    )    )



Phase 2: Add HTTP API

└─ cmd/api/main.go

└─ internal/adapters/http/    // HTTP Routes    // HTTP Routes



Phase 3: Add Web UI    http.HandleFunc("/convert", http.ConvertHandler(service))    http.HandleFunc("/convert", http.ConvertHandler(service))

└─ cmd/web/main.go

└─ internal/adapters/web/



Phase 4: Add Database    log.Fatal(http.ListenAndServe(":8080", nil))    log.Fatal(http.ListenAndServe(":8080", nil))

└─ internal/adapters/postgres/

}}

Phase 5: Add Cloud Storage

└─ internal/adapters/s3/````

└─ internal/adapters/gcs/



→ ALL with SAME core logic!

```**Now you have**:**Agora você tem**:



---



## 🚀 Example: Your Own Extension- CLI: `go run ./cmd/govc -p 4 /path`- CLI: `go run ./cmd/govc -p 4 /path`



**Want to add feature X? Follow this template:**- HTTP API: `curl -X POST http://localhost:8080/convert`- HTTP API: `curl -X POST http://localhost:8080/convert`



```go

package your_adapter

------

import "github.com/wallissonmarinho/GoVC/internal/core/ports"



// Declare which port you implement

var _ ports.YourPort = (*YourAdapter)(nil)## 🔄 Case 3: Replace FFmpeg with MediaInfo## 🔄 Caso 3: Trocar FFmpeg por MediaInfo



type YourAdapter struct {

    // Your configuration

}### Problem### Problema



func New() *YourAdapter {

    return &YourAdapter{}

}You want to use `mediainfo` instead of `ffmpeg` to extract metadata.Você quer usar `mediainfo` em vez de `ffmpeg` para extrair metadados.



// Implement all methods from ports.YourPort

func (s *YourAdapter) Method1() error {

    // TODO: Your logic### Solution### Solução

    return nil

}



func (s *YourAdapter) Method2() string {**Create new adapter** `internal/adapters/mediainfo/adapter.go` implementing `VideoConverterPort`.**Criar novo adapter** `internal/adapters/mediainfo/adapter.go` implementando `VideoConverterPort`.

    // TODO: Your logic

    return ""

}

```go```go

// Add in bootstrap (cmd/govc/main.go)

// adapter := your_adapter.New()package mediainfopackage mediainfo

```



---

import (import (

**Easy extension, without breaking anything!** 🎉

    "os/exec"    "os/exec"

    "github.com/wallissonmarinho/GoVC/internal/core/domain"    "github.com/wallissonmarinho/GoVC/internal/core/domain"

))



type MediaInfoAdapter struct {}type MediaInfoAdapter struct {}



func (m *MediaInfoAdapter) Convert(video *domain.Video, inputDir string) error {func (m *MediaInfoAdapter) Convert(video *domain.Video, inputDir string) error {

    // Use mediainfo instead of ffmpeg    // Usar mediainfo em vez de ffmpeg

    cmd := exec.Command("mediainfo", video.Path)    cmd := exec.Command("mediainfo", video.Path)

    return cmd.Run()    return cmd.Run()

}}



func (m *MediaInfoAdapter) GetDuration(videoPath string) (float64, error) {func (m *MediaInfoAdapter) GetDuration(videoPath string) (float64, error) {

    // Parse mediainfo output    // Parse mediainfo output

    return 0, nil    return 0, nil

}}



// ... etc// ... etc

`````

**Then, just replace in bootstrap:\*\***Depois, basta trocar no bootstrap:\*\*

`go`go

converterAdapter := mediainfo.NewMediaInfoAdapter() // ← New!converterAdapter := mediainfo.NewMediaInfoAdapter() // ← Novo!

````



------



## 🗄️ Case 4: Add Database (PostgreSQL)## 🗄️ Caso 4: Adicionar Banco de Dados (PostgreSQL)



### New Adapter### Novo Adapter



Create `internal/adapters/postgres/adapter.go` to track conversions in DB.Criar `internal/adapters/postgres/adapter.go` para rastrear conversões no BD.



```go```go

package postgrespackage postgres



import (import (

    "database/sql"    "database/sql"

    _ "github.com/lib/pq"    _ "github.com/lib/pq"

))



type PostgresAdapter struct {type PostgresAdapter struct {

    db *sql.DB    db *sql.DB

}}



func NewPostgresAdapter(connStr string) (*PostgresAdapter, error) {func NewPostgresAdapter(connStr string) (*PostgresAdapter, error) {

    db, err := sql.Open("postgres", connStr)    db, err := sql.Open("postgres", connStr)

    if err != nil {    if err != nil {

        return nil, err        return nil, err

    }    }

    return &PostgresAdapter{db: db}, nil    return &PostgresAdapter{db: db}, nil

}}



func (p *PostgresAdapter) LogConversion(filename string, success bool) error {func (p *PostgresAdapter) LogConversion(filename string, success bool) error {

    _, err := p.db.Exec(    _, err := p.db.Exec(

        "INSERT INTO conversions (filename, success, timestamp) VALUES ($1, $2, NOW())",        "INSERT INTO conversions (filename, success, timestamp) VALUES ($1, $2, NOW())",

        filename, success,        filename, success,

    )    )

    return err    return err

}}

````

**Use in service** (can extend `ProgressReporterPort` or create new port).**Usar no service** (pode estender `ProgressReporterPort` ou criar new port).

---

## 📊 Quick Reference## 📊 Quick Reference

### Add New Adapter (Checklist)### Adicionar Novo Adapter (Checklist)

- [ ] Create folder: `internal/adapters/new_adapter/`- [ ] Criar pasta: `internal/adapters/novo_adapter/`

- [ ] Create file: `adapter.go`- [ ] Criar arquivo: `adapter.go`

- [ ] Implement relevant interface(s) from `ports/ports.go`- [ ] Implementar interface(s) relevante(s) de `ports/ports.go`

- [ ] Add assertion: `var _ ports.YourPort = (*YourAdapter)(nil)`- [ ] Adicionar verificação: `var _ ports.SeuPort = (*SeuAdapter)(nil)`

- [ ] Update bootstrap in `cmd/govc/main.go`- [ ] Atualizar bootstrap em `cmd/govc/main.go`

- [ ] Build & test: `go build ./cmd/govc`- [ ] Build & test: `go build ./cmd/govc`

- [ ] Document usage in your README- [ ] Documentar uso em seu README

---

## 🧪 Test New Adapter (Mock)## 🧪 Testar Novo Adapter (Mock)

`go`go

package mainpackage main

import (import (

    "testing"    "testing"

    "github.com/wallissonmarinho/GoVC/internal/adapters/your_adapter"    "github.com/wallissonmarinho/GoVC/internal/adapters/seu_adapter"

))

func TestYourAdapter(t *testing.T) {func TestSeuAdapter(t *testing.T) {

    adapter := your_adapter.New()    adapter := seu_adapter.New()



    // Test method    // Test método

    result, err := adapter.YourMethod()    result, err := adapter.SeuMetodo()



    if err != nil {    if err != nil {

        t.Fatalf("Error: %v", err)        t.Fatalf("Error: %v", err)

    }    }



    if result != expected {    if result != expected {

        t.Errorf("got %v, want %v", result, expected)        t.Errorf("got %v, want %v", result, expected)

    }    }

}}

```



------



## 📈 Escalation Path## 📈 Escalation Path



```

Phase 1: CLI (current) ✅Phase 1: CLI (current) ✅

└─ cmd/govc/main.go └─ cmd/govc/main.go

Phase 2: Add HTTP APIPhase 2: Add HTTP API

└─ cmd/api/main.go └─ cmd/api/main.go

└─ internal/adapters/http/ └─ internal/adapters/http/

Phase 3: Add Web UIPhase 3: Add Web UI

└─ cmd/web/main.go └─ cmd/web/main.go

└─ internal/adapters/web/ └─ internal/adapters/web/

Phase 4: Add DatabasePhase 4: Add Database

└─ internal/adapters/postgres/ └─ internal/adapters/postgres/

Phase 5: Add Cloud StoragePhase 5: Add Cloud Storage

└─ internal/adapters/s3/ └─ internal/adapters/s3/

└─ internal/adapters/gcs/ └─ internal/adapters/gcs/

→ ALL with SAME core logic!→ ALL com MESMA core logic!

````



------



## 🚀 Example: Your Own Extension## 🚀 Exemplo: Sua Própria Extensão



**Want to do X? Follow this template:****Quer fazer X? Siga este template:**



```go```go

package your_adapterpackage seu_adapter



import "github.com/wallissonmarinho/GoVC/internal/core/ports"import "github.com/wallissonmarinho/GoVC/internal/core/ports"



// Declare which port you implement// Declare which port you implement

var _ ports.YourPort = (*YourAdapter)(nil)var _ ports.SeuPort = (*SeuAdapter)(nil)



type YourAdapter struct {type SeuAdapter struct {

    // Your configuration    // Sua configuração

}}



func New() *YourAdapter {func New() *SeuAdapter {

    return &YourAdapter{}    return &SeuAdapter{}

}}



// Implement all methods from ports.YourPort// Implement all methods from ports.SeuPort

func (s *YourAdapter) Method1() error {func (s *SeuAdapter) Metodo1() error {

    // TODO: Your logic    // TODO: Sua lógica

    return nil    return nil

}}



func (s *YourAdapter) Method2() string {func (s *SeuAdapter) Metodo2() string {

    // TODO: Your logic    // TODO: Sua lógica

    return ""    return ""

}}



// Add in bootstrap (cmd/govc/main.go)// Adicione no bootstrap (cmd/govc/main.go)

// adapter := your_adapter.New()// adapter := seu_adapter.New()

````

---

**Easy extension, without breaking anything!** 🎉**Extensão fácil, sem quebrar nada!** 🎉

# Hexagonal Architecture Best Practices — GoVC

## 📐 Princípios Fundamentais

### 1. Separação por Camadas

```
┌─────────────────────────────────────────────────────────┐
│                     External World                      │
├─────────────────────────────────────────────────────────┤
│                 Adapters (Input/Output)                 │
│  CLI  │  HTTP  │  Files  │  Database  │  Message Queue  │
├─────────────────────────────────────────────────────────┤
│                    Ports (Interfaces)                   │
│         Abstract contracts between core and adapters    │
├─────────────────────────────────────────────────────────┤
│                  Core (Business Logic)                  │
│      Domain Entities  │  Use Cases  │  Services         │
├─────────────────────────────────────────────────────────┤
│                    Ports (Interfaces)                   │
│         Abstract contracts between core and adapters    │
├─────────────────────────────────────────────────────────┤
│                 Adapters (Input/Output)                 │
│  CLI  │  HTTP  │  Files  │  Database  │  Message Queue  │
└─────────────────────────────────────────────────────────┘
```

---

## 📁 Estrutura de Pastas Correta

```
GoVC/
│
├── cmd/                                 ← Entry points (binários)
│   └── govc/
│       └── main.go                      ← Bootstrap (DI)
│
├── internal/                            ← Código privado (não exportável)
│   ├── core/                            ← ❤️  CORAÇÃO (lógica pura)
│   │   ├── domain/                      ← Entities (modelos puros)
│   │   │   ├── video.go
│   │   │   ├── conversion.go
│   │   │   └── progress.go
│   │   ├── ports/                       ← Interfaces (contratos)
│   │   │   ├── config.go
│   │   │   ├── video_discovery.go
│   │   │   ├── video_converter.go
│   │   │   ├── file_system.go
│   │   │   ├── progress_reporter.go
│   │   │   ├── service_command.go
│   │   │   ├── executor.go
│   │   │   └── command_executor.go
│   │   └── services/                    ← Use Cases (orquestração)
│   │       └── conversion_service.go
│   │
│   └── adapters/                        ← Implementações concretas
│       ├── cli/                         ← Input Adapter: CLI
│       │   ├── config.go                ✅ Implementa ConfigPort
│       │   ├── logger.go                ✅ Implementa ProgressReporterPort
│       │   ├── command_executor.go      ✅ Implementa CommandExecutorPort
│       │   ├── convert_command.go       ✅ Implementa ServiceCommand
│       │   ├── config_mock.go           ✅ Mock (para testes)
│       │   ├── logger_mock.go           ✅ Mock (para testes)
│       │   ├── command_executor_mock.go ✅ Mock (para testes)
│       │   └── convert_command_mock.go  ✅ Mock (para testes)
│       │
│       ├── commands/                    ← urfave/cli Handler Adapter
│       │   ├── convert.go               ✅ ConvertCommandHandler
│       │   ├── factory.go               ✅ CommandFactory
│       │   ├── convert_mock.go          ✅ Mock (para testes)
│       │   └── factory_mock.go          ✅ Mock (para testes)
│       │
│       ├── filesystem/                  ← Output Adapter: FileSystem
│       │   ├── adapter.go               ✅ Implementa VideoDiscoveryPort + FileSystemPort
│       │   └── adapter_mock.go          ✅ Mock (para testes)
│       │
│       └── ffmpeg/                      ← Output Adapter: FFmpeg
│           ├── adapter.go               ✅ Implementa VideoConverterPort
│           └── adapter_mock.go          ✅ Mock (para testes)
│
├── tests/                               ← ⭐ TESTES (separados!)
│   ├── unit/                            ← Testes unitários
│   │   ├── core/
│   │   │   ├── domain/
│   │   │   │   ├── video_test.go
│   │   │   │   ├── conversion_test.go
│   │   │   │   └── progress_test.go
│   │   │   └── services/
│   │   │       └── conversion_service_test.go
│   │   └── adapters/
│   │       ├── commands/
│   │       │   ├── convert_test.go
│   │       │   └── factory_test.go
│   │       ├── cli/
│   │       │   ├── config_test.go
│   │       │   └── logger_test.go
│   │       ├── filesystem/
│   │       │   └── adapter_test.go
│   │       └── ffmpeg/
│   │           └── adapter_test.go
│   │
│   ├── integration/                     ← Testes de integração
│   │   ├── conversion_flow_test.go
│   │   └── cli_integration_test.go
│   │
│   └── testdata/                        ← Fixtures para testes
│       ├── videos/
│       ├── configs/
│       └── mocks/
│
├── main.go                              ← Stub (aponta para cmd/govc)
├── go.mod
├── go.sum
└── README.md
```

---

## 🎯 Regras da Arquitetura Hexagonal

### ✅ PERMITIDO (segue padrão)

```go
// 1. Core pode usar APENAS:
//    - domain/: entities puras
//    - ports/: interfaces
//    - services/: use cases

// 2. Adapters podem usar:
//    - core (via ports): depende de interfaces
//    - external tools (ffmpeg, filesystem, etc)

// 3. Testes podem usar:
//    - Mocks (ficam nos adapters)
//    - Qualquer coisa para testar

// 4. Main/Bootstrap pode:
//    - Criar todas as instâncias
//    - Injetar dependências
//    - Iniciar a aplicação
```

### ❌ PROIBIDO (viola padrão)

```go
// ❌ Core usando adapters diretamente
import "github.com/wallissonmarinho/GoVC/internal/adapters/cli"

// ❌ Adapters sem implementar interface
type MyAdapter struct {}
func (a *MyAdapter) DoSomething() {}  // Não implementa nenhuma port!

// ❌ Testes dentro de adapters/core (misturar código)
// ✅ Melhor: tests/unit/adapters/commands/convert_test.go

// ❌ Bootstrap fora do main (misturar lógica)
// ✅ Melhor: cmd/govc/main.go apenas cria e injeta
```

---

## 📦 Onde Colocar Cada Coisa

### Arquivos de Implementação

| Tipo                       | Lugar              | Exemplo                 |
| -------------------------- | ------------------ | ----------------------- |
| Domain Entity              | `core/domain/`     | `video.go`              |
| Port Interface             | `core/ports/`      | `video_converter.go`    |
| Use Case Service           | `core/services/`   | `conversion_service.go` |
| **Adapter Implementation** | `adapters/{type}/` | `cli/config.go`         |
| **Adapter Mock**           | `adapters/{type}/` | `cli/config_mock.go`    |
| Bootstrap/DI               | `cmd/{app}/`       | `main.go`               |

### Arquivos de Teste

| Tipo              | Lugar                       | Exemplo                      |
| ----------------- | --------------------------- | ---------------------------- |
| Domain Tests      | `tests/unit/core/domain/`   | `video_test.go`              |
| Service Tests     | `tests/unit/core/services/` | `conversion_service_test.go` |
| Adapter Tests     | `tests/unit/adapters/`      | `cli/config_test.go`         |
| Integration Tests | `tests/integration/`        | `conversion_flow_test.go`    |

---

## 🔄 Fluxo de Dependência (sempre de fora para dentro)

```
main.go (bootstrap)
   ↓
Adapters (criam instâncias)
   ↓
Ports (implementadas)
   ↓
Core (recebe via injeção)
```

**Nunca ao contrário!**

```go
// ✅ CORRETO
func main() {
    adapter := cli.NewCLIConfig()           // Adapter criado aqui
    service := services.NewConversionService(adapter, ...) // Injetado
    service.Execute()
}

// ❌ ERRADO
// ConversionService criando CLIConfig diretamente
// (tight coupling, difícil de testar)
```

---

## 🧪 Padrão de Testes

### Teste Unitário (Domain/Service)

```go
// tests/unit/core/domain/video_test.go
package domain

import (
    "testing"
    "github.com/wallissonmarinho/GoVC/internal/core/domain"
)

func TestVideo_NewVideo(t *testing.T) {
    video := domain.NewVideo("/path/video.mkv")
    if video.Filename() != "video.mkv" {
        t.Fatal("Expected filename")
    }
}
```

### Teste Unitário (Adapter)

```go
// tests/unit/adapters/cli/config_test.go
package cli

import (
    "testing"
    "github.com/wallissonmarinho/GoVC/internal/adapters/cli"
)

func TestCLIConfig_NewCLIConfig(t *testing.T) {
    config, err := cli.NewCLIConfigFromContext(4, true, "/path")
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if config.GetWorkers() != 4 {
        t.Error("Expected 4 workers")
    }
}
```

### Teste Unitário (com Mock)

```go
// tests/unit/core/services/conversion_service_test.go
package services

import (
    "testing"
    "github.com/wallissonmarinho/GoVC/internal/adapters/cli"
    "github.com/wallissonmarinho/GoVC/internal/core/services"
)

func TestConversionService_Execute(t *testing.T) {
    // Usar mock
    mockConverter := &cli.MockFFmpegAdapter{}
    service := services.NewConversionService(
        mockDiscovery,
        mockConverter,  // Mock usado aqui!
        mockFS,
        mockReporter,
        mockConfig,
    )

    err := service.Execute()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
}
```

### Teste de Integração

```go
// tests/integration/conversion_flow_test.go
package integration

import (
    "testing"
    "github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
    "github.com/wallissonmarinho/GoVC/internal/core/services"
)

func TestConversionFlow_End2End(t *testing.T) {
    // Usa adaptadores REAIS
    discoveryAdapter := filesystem.NewFilesystemAdapter()
    converterAdapter := ffmpeg.NewFFmpegAdapter()

    service := services.NewConversionService(...)

    // Testa o fluxo todo
    err := service.Execute()
    if err != nil {
        t.Fatalf("End-to-end failed: %v", err)
    }
}
```

---

## 🛠️ Checklist para Novos Adapters

Quando adicionar um novo adapter, seguir:

- [ ] Criar pasta em `internal/adapters/{tipo}/`
- [ ] Implementar arquivo `adapter.go`
- [ ] Implementar `*_mock.go` para testes
- [ ] Criar testes em `tests/unit/adapters/{tipo}/`
- [ ] Injetar no bootstrap (`cmd/govc/main.go`)
- [ ] Documentar em `EXTENSION_GUIDE.md`

---

## 📋 Exemplo: Adicionar HTTP API Input Adapter

### 1. Criar estrutura

```
internal/adapters/http/
├── adapter.go           ← Implementação
├── config_adapter.go    ← Implementa ConfigPort
└── config_mock.go       ← Mock
```

### 2. Implementar

```go
// internal/adapters/http/config_adapter.go
package http

import "github.com/wallissonmarinho/GoVC/internal/core/ports"

type HTTPConfigAdapter struct {
    // campos
}

func (a *HTTPConfigAdapter) GetInputDir() string {
    // implementação
}

// Garante que implementa ConfigPort
var _ ports.ConfigPort = (*HTTPConfigAdapter)(nil)
```

### 3. Testar

```go
// tests/unit/adapters/http/config_adapter_test.go
package http

import (
    "testing"
    "github.com/wallissonmarinho/GoVC/internal/adapters/http"
)

func TestHTTPConfigAdapter_GetInputDir(t *testing.T) {
    // teste
}
```

### 4. Registrar no bootstrap

```go
// cmd/govc/main.go
httpConfig := http.NewHTTPConfigAdapter(":8080")
service := services.NewConversionService(
    discoveryAdapter,
    converterAdapter,
    fileSystemAdapter,
    reporterAdapter,
    httpConfig,  // ← novo adapter
)
```

---

## ✅ Verificação

```bash
# Compilar (sem warnings)
go build ./...

# Testes (todos passando)
go test ./...

# Go vet (sem problemas)
go vet ./...

# Imports (sem imports cruzados indevidos)
go build -v ./cmd/govc
```

---

## 📚 Referências

- [Alistair Cockburn - Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture - Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design)

---

**🎯 GoVC segue rigorosamente as práticas de Hexagonal Architecture!**

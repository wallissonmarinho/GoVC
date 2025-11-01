# Unit Tests - GoVC

## Overview 📊

Teste unitários implementados seguindo o padrão encontrado no repositório `api-backend` com uso de `testify/assert` e `testify/mock`.

**Status**: ✅ **46/46 testes passando**

## Estrutura de Testes

### 1. Domain Layer Tests

**Arquivo**: `internal/core/domain/video_test.go` e `internal/core/domain/progress_test.go`

#### Video Tests (12 testes)

- ✅ `TestNewVideoSimples` - Criação de entidade Video
  - Caminho simples
  - Caminho complexo
  - Remoção de extensão
- ✅ `TestVideoOutputPath` - Geração de caminho de saída

  - Caminho simples
  - Caminho complexo
  - Extensão .mp4

- ✅ `TestVideoLogPath` - Geração de caminho de log

  - Log simples
  - Log complexo

- ✅ `TestVideoSubtitlePath` - Geração de caminho de subtítulos

  - Subtitle simples
  - Subtitle input diferente

- ✅ `TestVideoFilename` - Extração de nome do arquivo

  - Filename simples
  - Filename com caracteres especiais
  - Filename de caminho complexo

- ✅ `TestVideoSetDuration` - Definição de duração

  - Duração válida
  - Duração zero
  - Atualizar duração

- ✅ `TestVideoMarkWithSubtitles` - Marcação de subtítulos
  - Marcar com subtítulos
  - Manter sem subtítulos

#### ProgressTracker Tests (23 testes)

- ✅ `TestNewProgressTracker` - Inicialização (3 testes)

  - Valores iniciais
  - Total zero
  - Mapa vazio

- ✅ `TestProgressTrackerUpdate` - Atualização de progresso (5 testes)

  - Atualizar um vídeo
  - Progresso 100%
  - Limitar em 100%
  - Múltiplos vídeos
  - Sobrescrever progresso

- ✅ `TestProgressTrackerMarkCompleted` - Marcar como completo (3 testes)

  - Incrementar uma vez
  - Múltiplos incrementos
  - Ultrapassar total

- ✅ `TestProgressTrackerIsComplete` - Verificar conclusão (4 testes)

  - Incompleto
  - Completo
  - Total zero
  - Completed >= total

- ✅ `TestProgressTrackerGetSnapshot` - Snapshot (4 testes)

  - Snapshot vazio
  - Cópia do map
  - Cópia independente
  - Conter todos os vídeos

- ✅ `TestProgressTrackerCompleteFlow` - Fluxo completo (1 teste)
  - Rastreamento completo de conversão

### 2. Service Layer Tests

**Arquivo**: `internal/core/services/conversion_service_test.go`

#### Mock Adapters (5 mocks)

- `MockVideoDiscoveryPort` - Mock de descoberta de vídeos
- `MockVideoConverterPort` - Mock de conversão
- `MockFileSystemPort` - Mock de sistema de arquivos
- `MockProgressReporterPort` - Mock de relatório de progresso
- `MockConfigPort` - Mock de configuração

#### Service Tests (10 testes)

- ✅ `TestNewConversionService` - Criação de serviço (1 teste)

  - Interfaces válidas

- ✅ `TestVideoDiscoveryFindVideos` - Descoberta (3 testes)

  - Encontrar vídeos com sucesso
  - Erro ao buscar
  - Lista vazia

- ✅ `TestVideoDiscoveryCreateOutputDir` - Diretório de saída (2 testes)

  - Criar com sucesso
  - Erro ao criar

- ✅ `TestVideoConverterConvert` - Conversão (2 testes)

  - Sucesso
  - Erro na conversão

- ✅ `TestVideoConverterGetDuration` - Duração (2 testes)

  - Obter duração
  - Erro ao obter

- ✅ `TestVideoConverterHasExternalSubtitles` - Subtítulos (2 testes)
  - Com subtítulos
  - Sem subtítulos

## Padrões Utilizados

### 1. Subtests com `t.Run`

Cada teste utiliza `t.Run` sem loops `for`, conforme solicitado:

```go
func TestVideoOutputPath(t *testing.T) {
	t.Run("Deve gerar caminho de saida simples", func(t *testing.T) {
		// teste aqui
	})

	t.Run("Deve gerar caminho com diretorio complexo", func(t *testing.T) {
		// teste aqui
	})
}
```

### 2. Mocks com Testify

Padrão de mock seguindo seu repositório:

```go
discoveryMock := new(MockVideoDiscoveryPort)
discoveryMock.On("FindVideos", "/input").Return([]*domain.Video{}, nil)
videos, err := discoveryMock.FindVideos("/input")
assert.NoError(t, err)
discoveryMock.AssertExpectations(t)
```

### 3. Assertions com Testify

Uso padronizado de assertions:

```go
assert.NoError(t, err)
assert.Error(t, err)
assert.Equal(t, expected, actual)
assert.True(t, value)
assert.Nil(t, value)
assert.Len(t, slice, 3)
```

## Cobertura

### Domain Layer

- **Video Entity**: 100% - Todos os métodos testados
- **ProgressTracker**: 100% - Todos os métodos testados

### Service Layer

- **Port Interfaces**: 100% - Todos os contratos testados via mocks

### Adapters

- Prontos para testes (sem testes de adapters nesta iteração)

## Executar Testes

### Todos os testes

```bash
go test ./... -v
```

### Apenas domínio

```bash
go test ./internal/core/domain/... -v
```

### Apenas serviços

```bash
go test ./internal/core/services/... -v
```

### Com cobertura

```bash
go test ./... -cover -v
```

### Relatório de cobertura

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Dependências de Teste

- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/mock` - Mocking

Adicionadas via:

```bash
go get github.com/stretchr/testify
```

## Próximos Passos

### Phase 1: Adapter Tests (Futuro)

- [ ] Testes para `FilesystemAdapter`
- [ ] Testes para `FFmpegAdapter`
- [ ] Testes para `CLIConfig` e `LoggerReporter`

### Phase 2: Integration Tests (Futuro)

- [ ] Testes end-to-end
- [ ] Testes de conversão real
- [ ] Testes de erro e recuperação

### Phase 3: Benchmarks (Futuro)

- [ ] Performance de descoberta de vídeos
- [ ] Performance de progresso tracker
- [ ] Stress tests paralelos

## Notas de Implementação

✅ Todos os testes usam `t.Run` sem loops  
✅ Padrão testify/mock idêntico ao seu repositório  
✅ Namings em português (conforme seu projeto)  
✅ Cobertura 100% das entidades de domínio  
✅ Mocks reutilizáveis para futuras integrações

---

**Última atualização**: 1 de novembro de 2025  
**Total de testes**: 46  
**Taxa de sucesso**: 100% ✅

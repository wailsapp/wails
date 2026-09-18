---
title: "Depurar falhas"
description: "Coletar relatórios de diagnóstico e minidumps de um aplicativo Wails"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

Quando algo dá errado em produção, os logs raramente são suficientes. Este guia apresenta o pacote `github.com/wailsapp/wails/v3/pkg/debug`, que produz relatórios de diagnóstico estruturados e, no Windows, minidumps completos que podem ser abertos no WinDbg ou no Visual Studio.

O pacote é intencionalmente pequeno e tem dois pontos de entrada:

- `debug.Dump(...)` — gravar um minidump do Windows do processo atual.
- `debug.Report(...)` — coletar um retrato detalhado de diagnóstico, com minidump opcional.

Ambos são executados com o token do usuário que faz a chamada. Sem elevação de privilégios, sem `SeDebugPrivilege`, sem `OpenProcess`: o dump do processo atual usa o pseudo-handle `GetCurrentProcess`, que dispensa verificações de ACL para acesso ao próprio processo.

## Referência rápida

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report` sempre retorna um `*CrashReport`, mesmo em uma falha parcial. Quem chama recebe tudo o que foi coletado com sucesso antes do erro, portanto o relatório nunca é `nil` quando `err != nil`.

## A estrutura `CrashReport`

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

Serialize com `encoding/json` para serviços de relatório de falhas ou percorra os campos diretamente se precisar apenas de parte deles.

## Integração com `PanicHandler`

O padrão mais útil é coletar um relatório e, opcionalmente, um minidump no momento em que o Wails captura um panic e enviar os dados combinados ao seu próprio fluxo de relatórios:

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

O Wails **não** chama `debug.Report` automaticamente. Relatórios de falhas frequentemente envolvem consentimento do usuário ou remoção de dados pessoais, portanto essa decisão pertence ao seu manipulador. Consulte o [guia de tratamento de panics](/guides/panic-handling/) para saber quando o Wails chama `PanicHandler` e como cobrir também goroutines criadas pelo seu código.

## Depuração manual

`debug.Dump` e `debug.Report` também são úteis fora de situações de panic:

- Detecção de travamentos: quando um watchdog disparar, gere um dump do processo para abri-lo no WinDbg e ver exatamente o que bloqueou cada thread.
- Relatórios de estados anormais: por exemplo, uma contagem de goroutines que cresce sem limite; capture um relatório, deixe o aplicativo continuar e investigue depois.
- Pacotes de suporte sob demanda: conecte uma opção de menu “Relatar um bug” a `debug.Report(debug.WithDump())` e anexe o resultado ao chamado de suporte do usuário.

## Minidumps do Windows na prática

Por padrão, os minidumps usam um conjunto de flags detalhado, mas compacto:

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

Isso produz dumps de aproximadamente 5–50 MB que preservam todas as pilhas de threads, tabelas de handles e listas de módulos, suficientes para reconstruir onde cada goroutine estava no momento do dump.

`debug.WithFullMemory()` (ou `debug.WithDumpFullMemory()` em `Report`) adiciona `MiniDumpWithFullMemory`, que captura todo o espaço de endereçamento do processo. É útil para investigar por que um ponteiro específico foi corrompido, mas os arquivos chegam a centenas de MB ou GB. Reserve essa opção para análises aprofundadas.

Abra o `.dmp` resultante no WinDbg com `windbg -z C:\path\to\your.dmp` ou em Arquivo → Abrir dump de falha no Visual Studio. Os símbolos normalmente são removidos dos builds de produção; combine o dump com o `.pdb` do mesmo build, ou com o próprio binário Go se ele contiver informações de depuração, para obter rastreamentos de pilha úteis.

## Plataformas diferentes do Windows

`debug.Dump` retorna um erro de recurso não implementado fora do Windows. Linux e macOS possuem outras ferramentas de análise pós-falha, como core dumps via `GOTRACEBACK=crash`, `rr` ou relatores específicos da plataforma, que não correspondem diretamente à API `MiniDumpWriteDump`. `debug.Report` continua funcionando: você recebe tudo, exceto um arquivo de dump.

O erro de `Dump` foi projetado para permitir verificação e continuidade controlada com funcionalidade reduzida:

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## Observações de segurança

- **Não é necessária elevação de privilégios.** O pacote evita deliberadamente `SeDebugPrivilege` e a abordagem `OpenProcess` / PID remoto usada por ferramentas de coleta de credenciais. Ele gera dumps apenas do processo que faz a chamada.
- **Os arquivos de dump contêm tudo o que está na memória.** Isso inclui credenciais em uso, tokens de autenticação, dados do usuário e chaves de criptografia. Trate um dump como uma imagem completa da memória: proteja-o em trânsito, remova dados sensíveis armazenados e considere expurgá-los antes de enviar a terceiros.
- **O caminho padrão é `os.TempDir()`.** No Windows, ele corresponde a `%LOCALAPPDATA%\Temp`: específico do usuário e não legível por todos, mas persistente. Mova os dumps para um local seguro se for mantê-los.

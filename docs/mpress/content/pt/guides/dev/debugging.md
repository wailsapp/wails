---
title: "Depuração"
description: "Investigação de problemas e análise do desempenho do aplicativo"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

Este guia demonstra várias ferramentas para inspecionar e investigar possíveis problemas de desempenho no seu aplicativo Wails usando

- [`runtime/trace`](https://pkg.go.dev/runtime/trace) para criar gráficos de desempenho e inspecioná-los em um navegador

## Criação de rastreamentos de desempenho

@steps
### Prepare seu aplicativo para o rastreamento
Próximo ao ponto de entrada do programa, verifique se existe um código semelhante ao seguinte

```go
// Create the file to store our trace data within
traceFile, err := os.Create("trace.out")
if err != nil {
  log.Fatalf("trace.out could not be created: %v", err)
}

// Start the trace
if err := trace.Start(traceFile); err != nil {
  _ = traceFile.Close()
  log.Fatalf("trace.start could not start: %v", err)
}

// Trace cleanup on exit. Alternatively,
defer func() {
  trace.Stop()
  _ = traceFile.Close()
}()

...Start your wails app here...
```

A saída será salva em `trace.out` no seu diretório de trabalho, com métricas que abrangem todo o tempo de execução do aplicativo. O rastreamento padrão contém dados bastante brutos; é altamente recomendável ler mais sobre rastreamento para adicionar informações contextuais e limitar o que realmente será registrado. Por exemplo: `WithRegion, NewTask, Log`

### Execute seu aplicativo para criar o rastreamento
Enquanto o aplicativo estiver em execução, um rastreamento contínuo será gerado até o momento do encerramento. Portanto, realize algumas ações no aplicativo e encerre-o quando terminar

### Instale os auxiliares de visualização
Algumas visualizações precisam do [Graphviz](https://graphviz.org/). Você pode verificar se ele está instalado executando `dot -V`

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### Visualize seu rastreamento
Com os dados de rastreamento disponíveis, podemos iniciar a interface web

```bash
go tool trace trace.out
```

Isso deve abrir o navegador padrão — recomenda-se usar um navegador baseado no Chrome — na página inicial do visualizador de eventos de rastreamento.

Para quem está começando, a tela `Syscall profile` provavelmente é a mais útil. Ela mostra em detalhes exatamente o que o programa está fazendo, juntamente com os tempos de execução

@end

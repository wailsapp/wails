---
title: "Wails v3.1.0: mais Go, menos obstáculos na compilação"
description: "Compilações experimentais de aplicações desktop sem CGO e Wake: uma forma mais simples de compilar, empacotar e assinar aplicações Wails."
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

Foi uma semana movimentada no Wails: compilações sem CGO, configuração de projetos mais fácil e menos dificuldades no empacotamento e na assinatura.

O Wails v3.1.0 dá início aos nossos lançamentos semanais de funcionalidades, disponibilizando novos recursos regularmente. As correções de bugs continuarão sendo lançadas em atualizações de patch noturnas assim que estiverem prontas.

## Correções e otimizações desta semana

Esta semana trouxe **100 correções, 21 melhorias de desempenho e quatro alterações de reforço de segurança** no framework, na CLI, nas ferramentas de compilação e na CI. Entre os destaques estão um posicionamento de janelas mais confiável em configurações com diferentes DPIs, correções na bandeja e no modo escuro do Linux, fornecimento de ativos mais eficiente e menor uso de memória durante o empacotamento MSIX. [Veja abaixo algumas das alterações](#alguns-destaques-desta-semana).

## Compilações experimentais sem CGO

O purego permite que o Go chame bibliotecas nativas do sistema sem CGo, a ponte habitual do Go para código C. Começamos a explorá-lo quando surgiu, mas manter dois backends era um desafio. Com o amadurecimento do purego e a IA generativa tornando essa manutenção mais administrável, estamos prontos para disponibilizar backends experimentais para macOS e Linux.

Agora você pode compilar e fazer compilação cruzada de aplicações desktop sem um compilador C:

```sh
WAILS_PUREGO=1 wails3 build
```

Depois de compilar os assets do frontend, você também pode compilar diretamente com Go:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

O backend do macOS usa AppKit e WebKit; o Linux é compatível com GTK4/WebKitGTK 6.0 e GTK3. Os usuários ainda precisam das bibliotecas nativas de runtime, e o empacotamento e a assinatura têm requisitos próprios. O CGo continua sendo o padrão enquanto coletamos feedback.

[Detalhes de implementação e das plataformas](https://github.com/wailsapp/wails/pull/6521)

## Por que Wake?

O Wails v2 ocultava detalhes demais do processo de compilação. No v3, os Taskfiles deixaram tudo acessível, mas a maioria das pessoas quer algumas personalizações sem precisar manter um sistema de compilação.

O Wake reúne em `wails.yaml` as opções mais comuns: aplicações, ferramentas CLI complementares, servidores, destinos e empacotamento. O Task continua sendo nosso sistema de compilação principal e está disponível quando você precisar da flexibilidade dele.

Agora o Wake oferece recarregamento ao vivo, compilações incrementais, cache de saídas e empacotamento paralelo. Já tem os binários? O `--skip-build` os empacota sem exigir a instalação do Go nem de ferramentas de frontend.

Ative-o com `WAILS_EXP_USE_WAKE=1`.

### Um início guiado

Com o Wake ativado, o `wails3 init -ui` orienta você nas escolhas de backend, destino e pacote, verifica suas ferramentas e salva a configuração do projeto.

@image{src="/assets/screenshots/wake/application.png" alt="O assistente de projetos do Wake exibindo opções de aplicação desktop e backend, incluindo a opção experimental Pure Go." expand}

Escolha os formatos de pacote e as opções de assinatura para cada destino:

@image{src="/assets/screenshots/wake/packages.png" alt="Seleção de pacotes do Wake para Linux, mostrando os formatos AppImage, DEB, RPM, Void e Arch Linux e a opção de assinar as saídas de lançamento." expand}

A verificação de ferramentas abaixo mostra uma compilação padrão para Linux:

@image{src="/assets/screenshots/wake/tools.png" alt="O Wake verificando o compilador C, Go, npm e pkg-config para uma compilação padrão para Linux." expand}

### O que posso compilar na minha máquina?

Compilação, empacotamento e assinatura têm requisitos diferentes. Veja o que você pode fazer em cada host:

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| Operação | Host Windows | Host macOS | Host Linux | Requisitos |
| Compilar aplicações desktop com purego | ✓ | ✓ | ✓ | Go e o backend experimental purego para destinos macOS/Linux |
| Criar pacotes MSIX | ✓ | ✓ | ✓ | Aplicação compilada; o Windows SDK não é necessário |
| Criar instaladores NSIS | ✓ | ✓ | ✓ | Aplicação compilada e `makensis` |
| Criar aplicações macOS e DMGs | ✓ | ✓ | ✓ | Aplicação compilada; o Go cuida da criação do pacote |
| Criar pacotes DEB, RPM, Arch e Void | ✓ | ✓ | ✓ | Aplicação compilada |
| Criar AppImages | ✓ | ✓ | ✓ | Runtime AppImage do destino; árvore do runtime de GUI do Linux para empacotamento cruzado |
| Assinar pacotes Windows com uma identidade exportada | ✓ | ✓ | ✓ | Certificado RSA PFX exportado e chave privada |
| Assinar aplicações Apple e DMGs com uma identidade exportada | ✓ | ✓ | ✓ | Identidade RSA PKCS#12 exportada; consulte a observação sobre validação abaixo |
| Assinar pacotes DEB, RPM e Arch | ✓ | ✓ | ✓ | Chave OpenPGP exportada |
| Enviar para notarização da Apple e anexar tíquetes | ✓ | ✓ | ✓ | Chave de API da equipe Apple e uma aplicação devidamente assinada |
| Usar o `codesign` da Apple ou perfis do Keychain | – | ✓ | – | macOS e a identidade de assinatura apropriada |
@end

As compilações com CGo ainda precisam de uma cadeia de ferramentas C. Hooks personalizados que modificam DMGs exigem macOS e `hdiutil`.

### Assinatura e notarização

O Wake permite assinar pacotes para Windows, Apple e Linux, incluindo assinaturas ad hoc da Apple. A assinatura multiplataforma usa identidades exportadas; identidades do repositório de certificados e do Keychain ainda exigem ferramentas nativas. Credenciais no nível do usuário mantêm as chaves privadas fora da configuração do projeto. A assinatura de AppImage e Void ainda não está incluída.

A notarização da Apple funciona por meio do Go com uma chave de API da equipe. Testamos o envio e a anexação de tíquetes no Linux, e o Gatekeeper aceitou o resultado. Essa aplicação foi assinada usando o [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html) da Apple; ainda não verificamos se a Apple aceita as assinaturas de Developer ID geradas pelo Go do Wake.

### Migração para o Wake

A migração é compatível com todos os nove modelos integrados. Ela é unidirecional: ativar `wails.yaml` substitui o roteamento do Taskfile. Comandos personalizados e substituições sem suporte são rejeitados para revisão.

[Funcionalidades, requisitos e detalhes de migração do Wake](https://github.com/wailsapp/wails/pull/6531)

## Alguns destaques desta semana

- **As janelas são posicionadas corretamente em telas com diferentes DPIs.** Os limites e as posições das telas agora são processados corretamente em configurações de macOS e Linux com diferentes escalas de exibição. [#6346](https://github.com/wailsapp/wails/pull/6346)
- **A bandeja e a detecção de tema no Linux funcionam de forma mais confiável.** Ícones da bandeja criados após a inicialização da aplicação não são mais descartados, as dicas de ferramenta são preenchidas e a detecção do modo escuro lê a configuração de aparência do ambiente de desktop desde a inicialização.
- **Os callbacks de eventos do Go mantêm a ordem para cada listener.** Cada listener é executado de forma serial na ordem de emissão; listeners diferentes podem ser executados simultaneamente, e emissões simultâneas não têm uma ordem definida. Os hooks de eventos são copiados sob um bloqueio para evitar uma condição de corrida de dados. [#6275](https://github.com/wailsapp/wails/pull/6275), [#6276](https://github.com/wailsapp/wails/pull/6276)
- **O cache de compilação funciona com mais projetos reais.** Espaços entre aspas ou escapados não causam mais falhas de cache nas saídas do Go, e os salvamentos simultâneos no cache são serializados para evitar arquivos corrompidos. [#6495](https://github.com/wailsapp/wails/pull/6495), [#6488](https://github.com/wailsapp/wails/pull/6488)
- **O fornecimento de ativos realiza menos trabalho desnecessário.** Extensões de arquivo conhecidas dispensam a detecção de MIME, e o registro de depuração desativado não prepara mais mensagens que seriam descartadas. Ambas as alterações reduzem o trabalho no caminho da requisição. [#6414](https://github.com/wailsapp/wails/pull/6414), [#6420](https://github.com/wailsapp/wails/pull/6420)
- **O empacotamento MSIX evita buffers do tamanho dos binários.** Os executáveis são transmitidos para o pacote em vez de serem lidos integralmente na memória. Anteriormente, uma cópia representativa de 16 MiB alocava cerca de 16.8 MB por operação. [#6472](https://github.com/wailsapp/wails/pull/6472)

## Experimente

Ambas as funcionalidades são experimentais. Experimente uma compilação multiplataforma, empacote sua aplicação e conte-nos onde encontrou dificuldades. Agradecemos a todos que estão testando, relatando problemas e contribuindo.

Um agradecimento especial a beesophie, da nossa comunidade no Discord, que criou o [Cook](https://github.com/SimTXNG/cook), um sistema de compilação incrível que nos mostrou algumas lacunas do nosso próprio sistema. Ver o que o Cook tornou possível nos ajudou a reconhecer onde poderíamos facilitar a compilação com Wails e serviu como uma verdadeira inspiração para este trabalho.

Obrigado, Beesophie!

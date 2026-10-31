---
title: "Builds YAML"
description: "Configure e execute o sistema experimental de build nativo do Wails v3 com um único arquivo wails.yaml."
sidebar:
  order: 2
  badge:
    text: "Experimental"
    variant: caution
---

@note{type="caution" title="Prévia experimental"}
A primeira versão aceita apenas configuração literal. Este experimento também cria protótipos de hooks de ciclo de vida com duração limitada, para que seu formato possa ser avaliado antes de serem propostos para lançamento. Estágios personalizados, expressões YAML e chamadas tipadas às ferramentas do Wails continuam adiados. Os projetos existentes continuam usando Taskfiles até que você ative explicitamente uma migração revisada.
@end

## Habilitar builds YAML

Defina `WAILS_EXP_USE_WAKE` em todos os terminais que executarem os comandos experimentais. Para Bash ou Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Para PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

A adesão é determinada pela presença da variável: até mesmo um valor vazio ou `false` habilita o experimento. Remova a definição dela (`unset WAILS_EXP_USE_WAKE` no Bash/Zsh ou `Remove-Item Env:WAILS_EXP_USE_WAKE` no PowerShell) para desabilitá-lo. Depois que `wails.yaml` for ativado, os comandos exigirão que o Wake seja reabilitado e nunca recorrerão a Taskfiles.
Sem a variável, o roteamento de builds e sessões de desenvolvimento por YAML e a geração de projetos YAML ficam desabilitados, e os comandos `migrate`, `config` e `eject` ficam indisponíveis.
Projetos que ainda não foram migrados continuam usando o runtime Task integrado separado.

## Configurar um projeto

Com a adesão habilitada, adicione `wails.yaml` à raiz do projeto. A presença desse arquivo seleciona o sistema de build nativo e faz o Wails ignorar todos os Taskfiles desse projeto. Mantenha a adesão habilitada para todos os comandos abaixo. A modeline do esquema habilita o preenchimento automático e a validação em linha em editores compatíveis com servidores de linguagem YAML.

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/wails.v3.json
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
frontend:
  directory: frontend
  install: [pnpm, install, --frozen-lockfile]
  build: [pnpm, run, build]
  dev: [pnpm, run, dev]
  output: frontend/dist
build:
  output: bin
  trim_path: true
  strip: true
profiles:
  release:
    targets:
      linux/amd64:
        formats: [deb, rpm, appimage]
```

Execute o build padrão para o host ou a solicitação de lançamento nomeada completa:

```bash
wails3 build
wails3 build release
```

Para um build avulso, selecione alvos e formatos separados por vírgulas sem adicionar um perfil:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm
```

## Fazer compilação cruzada de aplicativos desktop

Para o backend desktop experimental implementado exclusivamente em Go, defina `CGO_ENABLED: "0"` no ambiente de build. O Wake seleciona a tag `purego` para programas com GUI no macOS e no Linux e usa o Go diretamente em hosts Windows, macOS e Linux. Nesse modo, a criação de executáveis desktop não exige Docker, emulação, compilador C nem um SDK da plataforma de destino. As dependências do aplicativo também devem ser compatíveis com `CGO_ENABLED=0`.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  output: bin
  environment:
    CGO_ENABLED: "0"
```

```bash
wails3 build --targets windows/amd64,windows/arm64,linux/amd64,linux/arm64,darwin/amd64,darwin/arm64
wails3 build --targets darwin/universal --formats dmg
```

O segundo comando combina as duas arquiteturas do macOS e cria um DMG não assinado usando Go. Ele pode ser executado em qualquer host desktop compatível. A assinatura de código e a notarização da Apple continuam sendo operações separadas no macOS.

O modo resultante é selecionado de forma independente para cada executável. Um `build.environment.CGO_ENABLED` declarado, o ambiente do alvo ou o ambiente do binário substitui o valor do terminal que invocou o comando. Caso contrário, uma tag `purego` selecionada explicitamente, um programa CLI/servidor ou um alvo Windows usa `0` por padrão; programas com GUI no macOS e no Linux usam `1` por padrão. `--plan --json` informa o compilador, o modo CGo e as tags resultantes de cada executável. A compilação cruzada cria arquivos para o alvo; a execução de um alvo desktop incompatível ainda exige o host apropriado.

Com `CGO_ENABLED=1`, o Wake usa primeiro uma cadeia de ferramentas nativa compatível ou o compilador cruzado `CC` configurado; depois, usa o Zig para um build cruzado destinado ao Windows, quando disponível; por fim, recorre ao fallback existente com Docker/Podman, onde houver suporte. Um `CC` explícito deve fornecer o compilador, as bibliotecas e o sysroot do alvo; o Zig sozinho não fornece um sysroot do GTK/WebKit para Linux. O empacotamento e a assinatura podem exigir ferramentas adicionais, conforme descrito abaixo. Alvos móveis exigem CGo e os SDKs de suas plataformas.

A compilação somente com Go não inclui as bibliotecas de sistema do destino. Aplicativos Linux com GUI, inclusive executáveis complementares com GUI, ainda exigem o runtime GTK/WebKitGTK selecionado por suas tags de build. As dependências do pacote ou um ambiente de implantação apropriado devem fornecê-lo. Por padrão, pacotes que contêm somente programas CLI/servidor não adicionam essas dependências de GUI; substituições explícitas de dependências continuam disponíveis.

## Criar vários executáveis

O projeto define o programa principal. Defina `project.package` como o pacote Go dele (o padrão é `.`) e `project.kind` como `gui`, `cli` ou `server` (o padrão é `gui`). Programas adicionais usam entradas nomeadas em `binaries`:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
binaries:
  notes-tool:
    package: cmd/notes-tool
    kind: cli
    tags: [tools]
    ldflags: [-X, main.Mode=release]
    environment:
      CGO_ENABLED: "0"
profiles:
  release:
    targets:
      linux/amd64: {}
```

Cada programa tem seus próprios dados de entrada para compilação e cache. Os pacotes importados do projeto e os arquivos incorporados são rastreados; a edição de um programa independente não recompila os demais. Há suporte a executáveis adicionais no desktop, e eles são incluídos nos pacotes de aplicativo e instaladores padrão. Arquivos de pacote substitutos pertencentes ao usuário controlam o próprio conteúdo.

```bash
wails3 build --binaries notes,notes-tool --jobs 4
wails3 build --binaries notes-tool --skip-frontend --skip-bindings
```

Um perfil nomeado é uma solicitação completa; use estas substituições de seleção com um build anônimo. `--jobs` limita os workers de execução. Selecionar apenas programas CLI/servidor omite o trabalho de frontend.

Para ler uma versão sem executar o código do aplicativo, configure tanto `project.version_file` quanto `project.version_symbol`. O arquivo deve conter uma declaração estática de string em Go. O literal `project.version` poderá então ser omitido. Defina `commit_symbol: main.Commit` no projeto ou em um binário nomeado para injetar a revisão do repositório na variável de string desse executável. A ausência de metadados do repositório é um erro quando a injeção do commit é solicitada.

De `upx_level: 1` a `9`, o valor opcionalmente compacta executáveis de lançamento para Linux e Windows; zero desabilita a compactação. O UPX deve estar instalado. Builds de desenvolvimento e para macOS omitem a compactação. Uma falha de compilação ou compactação preserva o programa publicado anteriormente.

## Empacotar e assinar para distribuição

### Usar executáveis desktop pré-compilados

Use `--skip-build` quando outro sistema de build ou job de CI já tiver produzido os executáveis:

```bash
wails3 build --targets linux/amd64 --formats deb,rpm --skip-build
wails3 run --skip-build -- --example-mode
```

Coloque cada executável em `build.output` (o padrão é `bin`), usando `project.binary_name` para o executável principal e o nome de cada entrada `binaries` para os executáveis auxiliares. No Windows, os nomes terminam em `.exe`. Uma invocação para vários destinos lê os arquivos de cada destino em um subdiretório `<os>-<arch>`. Para uma GUI do macOS, o Wake lê o executável bruto em `bin/<name>`, se estiver presente; caso contrário, usa `bin/<name>.app/Contents/MacOS/<name>`. Uma entrada `darwin/universal` deve conter slices de executável para amd64 e arm64.

O Wake valida o tipo de arquivo e a arquitetura de destino, copia os executáveis para o workspace gerado e rastreia o conteúdo deles para o cache de pacotes. Entradas de executáveis ausentes, incompatíveis ou que sejam links simbólicos causam uma falha, em vez de serem recompiladas ou restauradas a partir de uma entrada antiga. Os executáveis fornecidos não são sobrescritos nem têm suas permissões alteradas com chmod. Alterações em qualquer executável invalidam os pacotes que o utilizam.

Esse modo ignora a instalação e os builds do frontend, os bindings, a compilação Go, o UPX e os hooks `before_build`/`after_build`. Ele ainda gera os recursos do pacote, monta pacotes de aplicativo do macOS e executa os hooks configurados de empacotamento e assinatura. As ferramentas de empacotamento e assinatura continuam sujeitas à tabela abaixo; as ferramentas de build de Go e do frontend não são necessárias. Os metadados do pacote vêm do manifesto atual, incluindo `version_file` quando configurado; nenhuma versão é inferida do binário. Use configurações correspondentes de build e empacotamento, inclusive quanto ao tipo GUI ou CLI.

`wails3 package --skip-build` também é compatível. Com `wails3 run`, os argumentos de inicialização e o ambiente ainda são aplicados, e o executável deve poder ser executado no host. Esse modo exige um projeto YAML ativo e só é compatível com alvos desktop de produção; ele não se aplica a sessões de desenvolvimento nem ao empacotamento para dispositivos móveis.

A criação de pacotes não assinados tem requisitos de host e ferramentas diferentes dos requisitos de assinatura:

| Saída | Hosts e ferramentas de criação | Assinatura |
| --- | --- | --- |
| MSIX do Windows | Windows, macOS ou Linux; gerador em Go, sem necessidade do SDK do Windows. | Host Windows e ferramentas de assinatura. |
| Instalador NSIS do Windows | Windows, macOS ou Linux; exige `makensis` nativo do host. | Host Windows e ferramentas de assinatura. |
| Aplicativo e DMG do macOS | Windows, macOS ou Linux; montagem em Go, mesclagem de binários universais e gerador portátil de DMG. | macOS com ferramentas e credenciais da Apple. |
| DEB, RPM, Arch e Void para Linux | Windows, macOS ou Linux; geradores nativos de pacotes em Go. | DEB/RPM usam suas ferramentas externas de assinatura; o Arch usa o assinador em Go abaixo. A saída do Void não é assinada. |
| AppImage do Linux | Host Linux e arquitetura de destino correspondentes; linuxdeploy externo e plugin de implantação do GTK. | Não há fluxo integrado de assinatura de AppImage. |

Os arquivos MSIX do Windows são criados diretamente em Go; a criação de um pacote não assinado não exige o SDK do Windows nem o Docker. A assinatura no Windows ainda exige um host Windows e suas ferramentas de assinatura. O NSIS usa o compilador NSIS externo e aceita definições literais por meio de `packages.nsis.defines`. No macOS, o Wails fornece uma localidade UTF-8 somente ao processo filho do compilador NSIS quando a localidade herdada não é UTF-8; os ambientes do terminal e do aplicativo permanecem inalterados. O instalador NSIS padrão omite o bootstrapper do WebView2 quando todos os executáveis empacotados são ferramentas de linha de comando ou servidores. Um executável principal ou complementar com interface gráfica mantém o WebView2. Scripts NSIS personalizados mantêm o bootstrapper para preservar a compatibilidade com sua própria lógica de instalação do runtime.

O Linux oferece suporte a pacotes AppImage, DEB, RPM, Arch Linux e Void. Os AppImages mantêm o limite externo do linuxdeploy e do plugin GTK. Defina `packages.appimage.update_info` para solicitar um arquivo zsync e `packages.appimage.checksum: true` para solicitar um arquivo SHA256. Os arquivos auxiliares são restaurados do cache de pacotes e relacionados no recibo de artefatos junto com o AppImage. Atualmente, o empacotamento de AppImage exige Linux e correspondência entre as arquiteturas do host e do destino. A detecção do GTK não encontra bibliotecas carregadas dinamicamente por uma interface gráfica implementada inteiramente em Go; portanto, este adaptador ainda não oferece suporte à implantação desse tipo de interface. Use um dos formatos de gerenciador de pacotes do Linux para esse modo de compilação.

Os pacotes Arch preservam versões de pré-lançamento nos metadados nativos do pacote; por exemplo, `1.2.3-beta.2` se torna `1.2.3beta.2-1`, que é ordenado antes da versão final. As dependências do Void aceitam nomes de pacotes simples, como `gtk4`, normalizados para `gtk4>=0`; comparações explícitas de versão e valores `name-version_revision` completos são preservados. Use nomes de pacotes do Void, não nomes de arquivos de biblioteca.

Para criar um pacote Arch e uma assinatura separada, use um arquivo exportado de chave privada PGP em formato ASCII Armor e, opcionalmente, indique o nome de uma variável de ambiente que contenha a frase secreta:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
linux:
  signing:
    certificate: signing/private-key.asc
    credential: ARCH_SIGNING_PASSPHRASE
profiles:
  release:
    targets:
      linux/amd64:
        formats: [archlinux]
        sign: true
```

Defina a variável de ambiente indicada no seu ambiente local ou no repositório de segredos da CI e execute `wails3 build release`. Omita `credential` se a chave não estiver criptografada. O assinador em Go é executado em todos os hosts de desktop compatíveis e não exige GPG. Ele publica o `.pkg.tar.zst` original e um arquivo `.pkg.tar.zst.sig` no mesmo diretório; ambos são registrados separadamente no recibo de artefatos. O conteúdo da chave não é preparado nem armazenado no cache de compilação. As assinaturas e sua publicação nunca são reutilizadas do cache. Se a assinatura falhar, o par publicado anteriormente será preservado; se a instalação de qualquer um dos arquivos falhar, a publicação reverterá os dois. Isso não torna atômicas para leitores simultâneos duas operações de renomeação no sistema de arquivos.

Os pacotes de aplicativos do macOS incluem programas complementares. Defina `darwin.microphone_usage_description` para adicionar o texto de uso e, quando nenhum arquivo de direitos personalizado for fornecido, o direito de entrada de áudio. As fontes personalizadas de plist e de direitos permanecem inalteradas. A assinatura e a notarização usam as ferramentas e credenciais existentes da Apple no macOS; a notarização exige um perfil de notário real. As saídas assinadas mantêm a extensão da plataforma, por exemplo, `notes.signed.app` e `notes.signed.dmg`.

## Controles de desenvolvimento

Durante `wails3 dev`, digite um comando seguido de Enter:

| Comando | Ação |
| --- | --- |
| `r` | Forçar uma nova compilação e reinicialização, inclusive quando o backend não tiver sido alterado. |
| `b` | Gerar novamente os bindings e recompilar. |
| `f` | Reiniciar o frontend. |
| `c` | Limpar o terminal. |
| `q` | Encerrar a sessão e seus processos filhos. |
| `h` | Exibir os controles. |

Quando uma compilação falha, o aplicativo atual continua em execução. Os controles exigem um terminal interativo; a entrada redirecionada por pipe não altera o comportamento da sessão.

## Inspecionar o trabalho antes de executá-lo

A inspeção do plano resolve o mesmo grafo imutável usado pela execução. O texto destina-se a pessoas; o JSON é um contrato versionado legível por máquina.

```bash
wails3 build release --plan
wails3 build release --plan --json
```

## Validar a configuração sem compilar

Use o comando dedicado de pré-verificação na CI, em uma tarefa do editor ou antes de fazer commit da configuração. Sem um perfil, ele valida o manifesto base e todos os perfis nomeados; com um perfil, verifica apenas essa solicitação completa:

```bash
wails3 config check
wails3 config check release
```

Falhas de sintaxe, esquema, semântica, caminho, destino, formato, assinatura, planejador e recursos do host definidos pela configuração usam o mesmo renderizador de diagnósticos de Build e Dev. Os diagnósticos identificam o campo e o intervalo no código-fonte, exibem a linha correspondente e um acento circunflexo e incluem uma sugestão de correção quando ela é útil. Os valores de ambiente são ocultados dos trechos do código-fonte exibidos nos diagnósticos. Falhas do host causadas apenas por opções da linha de comando continuam sendo erros comuns de comando, em vez de serem incorretamente relatadas como YAML inválido.

## Usar ativos pertencentes ao projeto

Os caminhos do manifesto fazem referência a arquivos que pertencem a você, como ícones, manifestos de plataforma, modelos de pacote e scripts de manutenção. O Wails nunca edita essas fontes. Ele as copia para um espaço de trabalho `.wails/` descartável antes que uma ferramenta as consuma. Os artefatos finais são armazenados em `build.output`, e um recibo verificado é gravado em `.wails/artifacts/receipt.json`.

## Implantar no Android

Os perfis de produção do Android geram AABs. Para implantação local, o Wails compila internamente um APK de depuração específico para a ABI, instala-o com o ADB e inicia o aplicativo:

```bash
wails3 android devices
wails3 android run --device SERIAL
wails3 android run --emulator AVD_NAME
```

Sem um seletor explícito, um terminal interativo apresenta os dispositivos conectados e os Android Virtual Devices configurados. A automação com vários destinos disponíveis deve fornecer `--device` ou `--emulator`.

Anexe um fluxo de logs restrito ao aplicativo e monitore o destino até que o aplicativo seja encerrado ou você pressione Ctrl+C:

```bash
wails3 android run --emulator AVD_NAME --logs
wails3 android run --emulator AVD_NAME --logs --stop-emulator
```

As execuções anexadas acompanham a reinicialização do aplicativo e informam desconexões do destino ou transições para o estado offline. Por padrão, os emuladores continuam em execução para agilizar as iterações. `--stop-emulator` é explícito, exige `--logs` e, ao término da sessão anexada, interrompe somente um destino que seja emulador. Se o Wails iniciar um novo emulador, mas a verificação de prontidão da inicialização falhar ou for cancelada, ele encerrará esse processo automaticamente.

## Executar hooks de ciclo de vida com escopo limitado

A interface experimental de hooks invoca um script pertencente ao projeto em cada fase declarada. Um hook não é um comando do shell: o Wails executa o arquivo diretamente, a partir da raiz do projeto, a menos que `directory` esteja definido.

Em hosts Unix, o script deve ser executável e começar com um shebang. Hosts Windows aceitam scripts `.cmd`, `.bat` e `.ps1` por meio de adaptadores fixos de interpretador; o Wails não analisa comandos de shell inline em nenhuma das plataformas.

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
  after_package:
    script: scripts/check-packages.sh
profiles:
  release:
    targets:
      linux/amd64: {}
```

As fases compatíveis são `before_build`, `after_build`, `before_package`, `after_package`, `before_sign` e `after_sign`. Conforme apropriado, os hooks de compilação são executados uma vez para o projeto compartilhado ou uma vez por destino. Os hooks de empacotamento e assinatura formam barreiras ao redor de todos os formatos solicitados para um destino, enquanto as operações de formato dentro da barreira permanecem simultâneas.

Os hooks recebem uma variável de ambiente pertencente ao Wails: `WAILS_HOOK_CONTEXT_FILE`. Ela indica um arquivo JSON protegido que contém o contexto versionado da invocação:

```json
{
  "version": 1,
  "phase": "after_build",
  "command": "build",
  "scope": "target",
  "project_dir": "/work/notes",
  "working_directory": "/work/notes",
  "profile": "release",
  "target": { "os": "linux", "arch": "amd64" },
  "output": "/work/notes/bin/notes",
  "declared_outputs": ["/work/notes/generated/version.go"]
}
```

O contexto contém caminhos resolvidos, mas nunca copia o ambiente herdado do processo nem segredos. O Wails o cria atomicamente em `.wails/hooks/context/`, com permissões de arquivo exclusivas para o proprietário. Ele é removido após a conclusão bem-sucedida ou o cancelamento. Em caso de falha, é mantido para diagnóstico, e seu caminho é incluído no erro. Um código de saída diferente de zero faz a etapa Plan normal falhar, e o cancelamento encerra o grupo de processos do hook.

Por padrão, os hooks são executados em todas as invocações. Habilite o cache somente quando o script tiver todas as entradas e saídas declaradas:

```yaml
version: 3
project:
  name: notes
  product_name: Notes
  identifier: com.example.notes
  version: 1.0.0
hooks:
  before_build:
    script: scripts/generate-version.sh
    cache: true
    inputs: [version.txt]
    outputs: [generated/version.go]
profiles:
  release:
    targets:
      linux/amd64: {}
```

O script e seu modo executável são entradas implícitas. Várias saídas devem compartilhar um único diretório que não seja a raiz; esse diretório é o artefato com escopo limitado que o Wails registra e restaura. A raiz de saída não pode conter o script nem uma entrada declarada, e as árvores de saída armazenadas em cache podem conter somente arquivos comuns e diretórios, nunca links simbólicos ou arquivos especiais. Um hook executado sempre impede o armazenamento em cache de suas operações posteriores, pois não é possível gerar com segurança uma impressão digital de efeitos colaterais não declarados. Um hook armazenável em cache preserva o cache posterior por meio da identidade de seu artefato declarado.

## Migrar um projeto existente

Crie e revise um rascunho inativo:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --json
```

As tarefas conservadoras de ciclo de vida sem argumentos que invocam um script pertencente ao projeto são convertidas em hooks. Outros comportamentos personalizados acessíveis do Taskfile são relatados como bloqueadores. Após revisar `wails.migrated.yaml`, ative-o explicitamente:

```bash
wails3 migrate --activate
```

A ativação cria `wails.yaml` de forma atômica. Ela não exclui, renomeia nem edita Taskfiles, `build/config.yml`, scripts ou ativos. Um YAML inválido falha como YAML; nunca recorre à rota legada.

## Materializar todos os valores padrão

Para inspecionar ou copiar a configuração resolvida completa, grave um snapshot inativo:

```bash
wails3 eject
```

Isso cria `wails.ejected.yaml`. A configuração ativa nunca é reescrita.

## Janela de compatibilidade

`wails3 package` e `wails3 sign` são aliases obsoletos do mesmo planejador; novas automações devem usar `wails3 build` com um perfil ou formatos. A implantação é condicionada às versões: a primeira versão torna o YAML o caminho para novos projetos, a versão seguinte emite um aviso sobre o roteamento pelo Taskfile legado e a próxima remove a execução do Taskfile junto com os aliases. Cada transição exige que a matriz de testes de aceitação dos hosts nativos dessa versão esteja totalmente aprovada.

Consulte a referência gerada de campos de [`wails.yaml`](/reference/wails-yaml-fields/) para ver todos os campos, tipos, valores padrão e regras de aplicabilidade.

A referência de configuração de [`wails3 run`](/reference/wails-yaml-run/) descreve os valores padrão de inicialização local, as substituições específicas de plataforma e o fallback de `go run .` quando não há manifesto.

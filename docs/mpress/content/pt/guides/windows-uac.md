---
title: "Configuração do UAC no Windows"
description: "Configure o Controle de Conta de Usuário (UAC) para seu aplicativo Wails no Windows"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

Plataformas relevantes: <span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

O Controle de Conta de Usuário (UAC) do Windows determina os privilégios de execução do seu aplicativo Wails. Por padrão, os aplicativos Wails v3 incluem uma configuração explícita do UAC no manifesto do Windows, garantindo um comportamento consistente em diferentes computadores.

## Níveis de execução do UAC

Os aplicativos do Windows podem solicitar diferentes níveis de execução por meio do arquivo de manifesto. O Wails v3 inclui automaticamente uma configuração do UAC com um nível de execução padrão, que você pode personalizar de acordo com as necessidades do aplicativo.

### Níveis de execução disponíveis

| Nível | Descrição | Caso de uso |
| --- | --- | --- |
| `asInvoker` | É executado com os mesmos privilégios do processo pai | Padrão para a maioria dos aplicativos |
| `highestAvailable` | É executado com os privilégios mais altos disponíveis para o usuário | Aplicativos que podem precisar de acesso elevado |
| `requireAdministrator` | Sempre exige privilégios de administrador | Utilitários do sistema e instaladores |

### Configuração padrão

Os aplicativos Wails v3 incluem uma configuração padrão do UAC no manifesto do Windows:

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

Essa configuração garante que seu aplicativo:

- Seja executado com os mesmos privilégios do processo que o iniciou
- Não exija elevação de privilégios por padrão
- Funcione de maneira consistente em diferentes computadores
- Não acione solicitações do UAC para usuários comuns

## Personalização da configuração do UAC

Como o Wails v3 incentiva os usuários a personalizar os recursos de compilação, você pode modificar a configuração do UAC editando diretamente o modelo de manifesto do Windows.

### Localização do modelo de manifesto

O modelo de manifesto do Windows está localizado em:

```
build/windows/wails.exe.manifest
```

### Modificação do nível de execução

Para alterar o nível de execução, edite o atributo `level` no elemento `requestedExecutionLevel`:

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### Exemplos

#### Aplicativo padrão (configuração padrão)

A maioria dos aplicativos deve usar o nível padrão `asInvoker`:

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### Utilitário do sistema

Aplicativos que precisam de acesso elevado quando disponível:

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### Ferramenta administrativa

Aplicativos que sempre exigem privilégios de administrador:

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## Acesso à interface do usuário

O atributo `uiAccess` controla se o aplicativo pode interagir com elementos de interface do usuário que tenham privilégios mais elevados. Na maioria dos casos, ele deve permanecer como `false`.

Defina-o como `true` somente se o aplicativo precisar:

- Enviar entradas para outros aplicativos
- Controlar a interface do usuário de outros aplicativos
- Acessar elementos da interface do usuário de processos com privilégios mais elevados

@note{type="caution" title="Requisitos de acesso à interface do usuário"}
Definir `uiAccess="true"` exige que o aplicativo esteja:

- Assinado digitalmente com um certificado emitido por uma autoridade certificadora confiável
- Instalado em um local seguro (Program Files ou Windows\System32)

@end

## Compilação com configurações personalizadas do UAC

Depois de modificar o modelo de manifesto, compile o aplicativo normalmente:

```bash
wails3 build
```

O processo de compilação incorporará automaticamente sua configuração personalizada do UAC ao executável.

## Verificação da configuração do UAC

Você pode verificar se as configurações do UAC foram incorporadas corretamente usando a ferramenta `go-winres`:

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

Em seguida, examine o arquivo de manifesto extraído para confirmar que a configuração do UAC está presente.

@note{type="tip" title="Persistência do manifesto"}
Ao contrário de alguns outros frameworks, a configuração do UAC do Wails v3 é incorporada diretamente ao executável durante a compilação, garantindo que ela seja preservada quando o aplicativo for copiado para outros computadores.

@end

## Solução de problemas

### As solicitações do UAC não aparecem

Se você definiu `requireAdministrator`, mas não vê solicitações do UAC:

- Verifique se o manifesto está incorporado corretamente ao executável
- Verifique se o aplicativo não está sendo executado a partir de um processo que já tenha privilégios elevados
- Certifique-se de que a sintaxe do manifesto seja um XML válido

### O aplicativo não inicia

Se o aplicativo não iniciar após as alterações no UAC:

- Verifique se há erros de XML na sintaxe do manifesto
- Verifique se o valor do nível de execução é válido
- Tente reverter para `asInvoker` para isolar o problema

### Comportamento inconsistente entre máquinas

Se o comportamento do UAC variar entre máquinas:

- Verifique se o manifesto está incorporado ao executável (e não é externo)
- Verifique se o executável não foi modificado após a compilação
- Verifique se as configurações do UAC do Windows estão habilitadas na máquina de destino

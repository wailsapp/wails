---
title: "Corrija a documentação"
description: "Envie um PR de correção para a documentação do Wails v3 usando o M-Press."
sourcePath: "contributing/documentation.md"
---

PRs de correção são bem-vindos. Corrija erros de digitação, links quebrados, exemplos desatualizados,  
explicações pouco claras ou traduções. Para uma correção exclusiva da documentação, você não precisa abrir uma issue nem ter  
um teste de código com falha.

## Visualize localmente

Faça um fork de [wailsapp/wails](https://github.com/wailsapp/wails/fork), clone seu fork  
e crie uma branch a partir de `master`.

Instale o gerador de documentação na versão fixada:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

Você também pode baixar um binário verificado da  
[versão v1.0.17 do M-Press](https://github.com/leaanthony/mpress/releases/tag/v1.0.17).

A partir da raiz do repositório do Wails:

```sh
mpress version
mpress dev
```

Edite os arquivos-fonte `.md` em `docs/mpress/content/`. O inglês é o idioma padrão  
e fica diretamente nesse diretório. As traduções existentes ficam em  
pastas de idiomas, como `fr/` e `id/`. A visualização é recompilada conforme você salva as alterações.

Preserve o bloco de metadados no início de cada página e os componentes `@...` / `@end`  
correspondentes. Parágrafos comuns, títulos, listas e código delimitado podem ser editados  
como texto. Não edite os arquivos gerados em `docs/mpress/site/`.

## Verifique a correção

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Verifique a página alterada no navegador e execute todos os exemplos de código que você modificar.  
Para correções de tradução, compare todo o trecho alterado com o texto em inglês.  
Mantenha intactos os comandos, nomes de APIs, links, exemplos de código e conexões dos diagramas.  
Use linguagem técnica natural, preserve os requisitos e as ressalvas e traduza,  
além da prosa, os rótulos visíveis dos diagramas, os rótulos de navegação e as descrições das imagens.

Todos os idiomas publicados devem ter uma tradução completa de cada página em inglês.  
Não use marcadores de posição em inglês nem páginas alternativas de fallback. Uma correção em uma tradução  
pode alterar apenas esse idioma. Se você alterar o significado do texto em inglês, atualize as páginas correspondentes  
nos outros idiomas publicados; páginas não relacionadas não precisam ser geradas novamente.

## Envie um pull request

Abra um PR direcionado a `master`. Descreva o que estava errado, explique sua correção  
e liste as verificações que você executou. Inclua capturas de tela para alterações visíveis no layout  
e detalhes da plataforma e da versão para exemplos de código alterados.

Não são necessárias credenciais da Cloudflare nem acesso a serviços privados. As verificações públicas do PR  
compilam e validam o site estático sem credenciais de implantação.

Para alterações de código e propostas de funcionalidades, consulte [Como contribuir com o Wails](/contributing/).  
Para conhecer os detalhes internos, consulte a [Visão geral técnica](/contributing/overview/).

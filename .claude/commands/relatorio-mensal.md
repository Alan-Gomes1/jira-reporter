# Gerar Relatório Mensal de Atividades

Gere o relatório mensal de prestação de serviços buscando as atividades no Jira e preenchendo o template HTML.

## Fluxo

### 1. Preparar o binário

Verifique se o binário `jira-reporter` existe na raiz do projeto. Se não existir, compile com:

```bash
cd /home/alangomes/projetos/jira-reporter && go build -o jira-reporter .
```

### 2. Determinar o mês do relatório

- Se o usuário especificou um mês/ano (ex: "03/2026"), use esse valor com a flag `-d`.
- Caso contrário, o CLI automaticamente usa o mês anterior - não precisa passar a flag `-d`.

### 3. Perguntar preferências ao usuário

Antes de executar, pergunte:
- **Formato**: HTML ou DOCX? (padrão: HTML)
- **Incluir tarefas de QA?** (flag `-q` inclui cards onde o usuário é QA assignee)

Se o usuário não quiser ser perguntado (ex: "gera logo", "gera rápido"), use os padrões: HTML sem QA.

### 4. Executar o CLI

Monte o comando baseado nas respostas:

```bash
cd /home/alangomes/projetos/jira-reporter && ./jira-reporter -n "Relatório de Prestação de Serviços" -p "./reports/" [-f html|docx] [-d MM/YYYY] [-q]
```

Flags disponíveis:
- `-n "nome"` - Nome do arquivo (sem extensão)
- `-p "./reports/"` - Diretório de saída
- `-f html` ou `-f docx` - Formato (padrão: html)
- `-d "MM/YYYY"` - Mês/ano específico (padrão: mês anterior)
- `-q` - Incluir tarefas de QA

### 5. Enriquecer descrições das atividades via API do Jira (padrão Problema / Solução)

O CLI gera o relatório, mas as descrições na seção "RESUMO DAS ATIVIDADES" frequentemente ficam truncadas ou genéricas (ex: "Descrição detalhada"). Este passo corrige isso buscando as descrições completas na API do Jira e reescrevendo cada item como um **resumo curto no padrão Problema / Solução** — descrevendo qual era o problema e o que foi feito para corrigir.

1. Leia o arquivo HTML gerado e extraia todos os IDs de issues (ex: `PSD-2788`, `PPRDJ-456`), **preservando a ordem em que aparecem** na tabela de atividades
2. Leia o `.env` para obter `EMAIL`, `API_KEY` e `URL`
3. Faça uma chamada à API do Jira para buscar `summary` e `description` completas:

```bash
source .env && curl -s -u "$EMAIL:$API_KEY" \
  -H "Content-Type: application/json" \
  "$URL/rest/api/3/search/jql" \
  -d '{
    "jql": "key in (ISSUE-1,ISSUE-2,...)",
    "fields": ["key","summary","description"],
    "maxResults": 50
  }'
```

**Importante**: A API antiga `/rest/api/3/search` foi descontinuada. Use `/rest/api/3/search/jql`.

4. Para cada issue, extraia o texto completo da descrição do campo `description` (formato Atlassian Document Format - ADF):
   - Percorra os nodes recursivamente extraindo o `text` de nodes do tipo `text` (trate `hardBreak` e fim de `paragraph`/`heading`/`listItem` como quebra de linha)
   - Use também o `summary` para entender o tema da tarefa
   - Salve o texto extraído em um arquivo temporário se for grande, para conseguir ler tudo

5. **Sintetize um resumo curto Problema / Solução para cada issue** (NÃO copie a descrição bruta do Jira):
   - **Problema**: 1 frase descrevendo o que estava errado / a necessidade (ex: bug, comportamento incorreto, lacuna funcional, requisito regulatório)
   - **Solução**: 1–2 frases objetivas descrevendo o que foi feito para resolver (a implementação/correção), mencionando arquivos/componentes-chave quando relevante
   - Escreva em português, em tom técnico e direto; cada resumo deve ter ~2 a 4 linhas no total
   - Para cards de bug: deixe claro o sintoma observado e a causa quando conhecida
   - Para cards de feature/épico: descreva a necessidade e a entrega; em cards-pai, sinalize "(card pai do épico)"
   - Faça escape de HTML nos textos

6. Monte cada parágrafo no formato exato abaixo e substitua todo o conteúdo da seção "RESUMO DAS ATIVIDADES", mantendo a mesma ordem das issues na tabela:

```html
<p><b>KEY:</b> <b>Problema:</b> ...uma frase... <b>Solução:</b> ...uma ou duas frases...</p>
```

   Dica de implementação: gere um script Python que (a) recebe o dicionário `{KEY: (problema, solucao)}`, (b) faz `html.escape` nos textos, (c) monta os `<p>` na ordem da tabela e (d) substitui via regex o conteúdo entre o cabeçalho `<td><b>RESUMO DAS ATIVIDADES</b></td>` e o `</td></tr></table>` seguinte.

7. Salve o arquivo HTML atualizado e (se o formato for DOCX) regenere/reconverta a partir do HTML enriquecido.

### 6. Verificar resultado

Após a execução e enriquecimento:
1. Verifique se o comando retornou sem erro
2. Confirme que o arquivo foi criado em `reports/`
3. Leia o arquivo HTML gerado para mostrar ao usuário um resumo:
   - Quantidade de atividades encontradas
   - Período do relatório (mês/ano)
   - Lista resumida das tarefas (key + summary)
4. Informe o caminho completo do arquivo gerado

### 7. Se houver erro

Erros comuns:
- **"missing required configuration"**: O arquivo `.env` está incompleto. Leia o `.env` e informe qual campo está faltando.
- **Erro de autenticação/401**: A API_KEY pode estar expirada. Peça ao usuário para gerar um novo token em https://id.atlassian.com/manage-profile/security/api-tokens
- **Nenhuma atividade encontrada**: O JQL pode não ter retornado resultados para o período. Sugira verificar se há cards atribuídos ao usuário no Jira para aquele mês, ou tentar com a flag `-q` para incluir tarefas de QA.
- **LibreOffice não encontrado** (só para DOCX): O formato DOCX requer LibreOffice instalado. Sugira usar HTML ou instalar o LibreOffice.

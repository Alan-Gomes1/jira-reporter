# Jira Reporter

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Aplicação CLI em Go para gerar relatórios mensais de prestação de serviços com base em tarefas do Jira. Busca automaticamente issues do mês anterior atribuídas ao usuário e gera relatórios em **HTML** ou **DOCX**.

## ✨ Funcionalidades

- 📊 Geração automática de relatórios mensais
- 📄 Suporte a múltiplos formatos: **HTML** e **DOCX**
- 🔗 Integração com Jira Cloud via API
- 📋 Template HTML personalizável
- ⚡ CLI simples e intuitiva
- 🤖 Slash command `/relatorio-mensal` para Claude Code, com enriquecimento das descrições via API do Jira

## 🏗️ Arquitetura

O projeto segue a arquitetura **MVC**

```
internal/
├── config/          # Configuração centralizada (Singleton)
├── model/           # Entidades de domínio (Issue, User, Report)
├── repository/      # Acesso a dados externos (Jira API)
├── service/         # Lógica de negócio e orquestração
└── view/            # Geradores de saída (HTML, DOCX)
```

| Camada         | Responsabilidade                                       |
| -------------- | ------------------------------------------------------ |
| **Config**     | Carregamento centralizado de variáveis de ambiente     |
| **Model**      | Estruturas de dados puras sem dependências externas    |
| **Repository** | Interface `JiraRepository` para acesso ao Jira (DIP)   |
| **Service**    | Orquestração da geração de relatórios                  |
| **View**       | Interface `ReportGenerator` para extensibilidade (OCP) |

## 🛠️ Tecnologias

| Tecnologia          | Descrição                              |
| ------------------- | -------------------------------------- |
| **Go 1.24+**        | Linguagem principal                    |
| **go-atlassian/v2** | Cliente para API do Jira               |
| **Cobra**           | Framework CLI                          |
| **godotenv**        | Gerenciamento de variáveis de ambiente |
| **LibreOffice**     | Conversão HTML → DOCX (opcional)       |

## 📥 Como Baixar e Usar (Para Usuários)

Se você quer apenas usar a ferramenta sem precisar instalar o Go, siga os passos abaixo.

### 1. Baixe a Versão Correta

1.  Vá para a [página de Releases](https://github.com/Alan-Gomes1/jira-reporter/releases) do projeto.
2.  Encontre a versão mais recente.
3.  Na seção **Assets**, baixe o arquivo `.zip` correspondente ao seu sistema operacional:
    - **Windows (64-bit):** `jira-reporter-windows-amd64.zip`
    - **Linux (64-bit):** `jira-reporter-linux-amd64.zip`
    - **macOS (Intel):** `jira-reporter-macos-amd64.zip`
    - **macOS (Apple Silicon M1/M2/M3):** `jira-reporter-macos-arm64.zip`

### 2. Extraia e Configure

1.  Descompacte o arquivo `.zip` em uma pasta de sua preferência.
2.  Dentro da pasta, renomeie o arquivo `env-example` para `.env`.
3.  Abra o arquivo `.env` com um editor de texto e preencha com suas credenciais.
4.  Use esse [link](https://id.atlassian.com/manage-profile/security/api-tokens) para gerar seu token de API no Jira.

### 3. Execute o Relatório

Abra seu terminal (CMD ou PowerShell no Windows) na pasta onde você extraiu os arquivos e execute o comando correspondente:

#### **No Windows**

```powershell
.\jira-reporter.exe
```

#### **No Linux ou macOS**

Primeiro, dê permissão de execução ao arquivo (só precisa fazer uma vez):

```bash
# Para Linux
chmod +x ./jira-reporter-linux-amd64

# Para macOS (Intel)
chmod +x ./jira-reporter-macos-amd64

# Para macOS (Apple Silicon)
chmod +x ./jira-reporter-macos-arm64
```

Depois, execute o programa:

```bash
# Exemplo para Linux
./jira-reporter-linux-amd64
```

O relatório será gerado na subpasta `reports/`.

---

## 🚀 Primeiros Passos (Para Desenvolvedores)

Siga estas instruções para configurar e executar o Jira Reporter a partir do código-fonte.

### Pré-requisitos

- Go 1.24 ou superior
- LibreOffice (apenas para geração de DOCX)

### Instalação e Configuração

1.  **Clone o repositório:**

    ```bash
    git clone https://github.com/Alan-Gomes1/jira-reporter.git
    cd jira-reporter
    ```

2.  **Configure as Variáveis de Ambiente:**

    Crie um arquivo `.env` na raiz do projeto copiando o arquivo `env-example`:

    ```bash
    cp env-example .env
    ```

    Edite o arquivo `.env` com suas credenciais:

    ```env
    # Credenciais Jira
    EMAIL="seu-email@exemplo.com"
    API_KEY="seu-token-api-jira"
    URL="https://seu-dominio.atlassian.net"

    # Dados do Relatório
    COMPANY_NAME="Nome da Empresa"
    CNPJ="00.000.000/0001-00"
    USER_NAME="Seu Nome Completo"
    ```

3.  **Instale as Dependências:**

    ```bash
    go mod tidy
    ```

### Executando a Aplicação

Para executar a aplicação e gerar um relatório:

```bash
# Gerar relatório HTML (padrão)
go run main.go

# Ou compile primeiro
go build -o jira-reporter .
./jira-reporter
```

### 📋 Opções de Linha de Comando

```bash
# Ver ajuda
./jira-reporter --help

# Especificar nome do relatório
./jira-reporter -n "meu-relatorio"

# Especificar caminho de saída
./jira-reporter -p "/caminho/para/saida"

# Gerar em formato DOCX (requer LibreOffice)
./jira-reporter -f docx

# Especificar mês/ano do relatório (formato MM/YYYY)
./jira-reporter -d "01/2025"

# Incluir cards onde o usuário está marcado como QA
./jira-reporter -q

# Gerar relatório de outubro/2024 em DOCX
./jira-reporter -d "10/2024" -f docx

# Gerar relatório com cards de QA em DOCX
./jira-reporter -d "10/2025" -f docx -q

# Combinando opções
./jira-reporter -n "relatorio-dezembro" -p "./relatorios" -f docx -d "12/2025" -q
```

| Flag           | Descrição                              | Padrão       |
| -------------- | -------------------------------------- | ------------ |
| `-n, --name`   | Nome do relatório                      | `report`     |
| `-p, --path`   | Diretório de saída                     | `reports/`   |
| `-f, --format` | Formato (`html` ou `docx`)             | `html`       |
| `-d, --date`   | Mês/ano do relatório (formato MM/YYYY) | mês anterior |
| `-q, --qa`     | Incluir cards onde o usuário é QA      | `false`      |

### 🔧 Build para Produção

```bash
# Build simples
go build -o jira-reporter .

# Build otimizado (menor tamanho)
go build -ldflags="-s -w" -o jira-reporter .

# Build para diferentes plataformas
GOOS=windows GOARCH=amd64 go build -o jira-reporter-windows-amd64.exe .
GOOS=linux GOARCH=amd64 go build -o jira-reporter-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o jira-reporter-macos-amd64 .
GOOS=darwin GOARCH=arm64 go build -o jira-reporter-macos-arm64 .
```

### 📝 Personalizando o Template

O arquivo `template.html` na raiz do projeto pode ser editado para personalizar a aparência do relatório. As variáveis disponíveis são:

| Variável                        | Descrição                     |
| ------------------------------- | ----------------------------- |
| `{{.User.CompanyName}}`         | Nome da empresa               |
| `{{.User.CNPJ}}`                | CNPJ da empresa               |
| `{{.User.Username}}`            | Nome do usuário               |
| `{{.DateWorked}}`               | Mês/Ano de competência        |
| `{{.Jira.Items}}`               | Lista de issues               |
| `{{.Jira.Items[].Key}}`         | Chave da issue (ex: PROJ-123) |
| `{{.Jira.Items[].Summary}}`     | Resumo da issue               |
| `{{.Jira.Items[].Description}}` | Descrição da issue            |
| `{{.Jira.Items[].Date}}`        | Data da issue                 |
| `{{.Jira.Items[].URL}}`         | URL da issue no Jira          |

---

## 🤖 Comando `/relatorio-mensal` (Claude Code)

O repositório traz um slash command para o [Claude Code](https://claude.com/claude-code), definido em `.claude/commands/relatorio-mensal.md`. Ele automatiza o fluxo completo do relatório e vai além do que a CLI faz sozinha: além de executar o binário, busca as descrições completas das issues na API do Jira e reescreve a seção **"Resumo das atividades"** no padrão **Problema / Solução**.

Isso resolve a principal limitação da geração automática: o campo de descrição vindo do Jira costuma chegar truncado ou genérico no template.

### Pré-requisitos

- Claude Code instalado e aberto na raiz do projeto
- Arquivo `.env` configurado (mesmas variáveis da seção de instalação)
- Binário compilado na raiz (`go build -o jira-reporter .`) — o comando compila automaticamente se não encontrar

### Uso

```
/relatorio-mensal
```

O comando pergunta o formato (HTML ou DOCX) e se deve incluir os cards em que você é QA assignee. Para pular as perguntas e usar os padrões (HTML, sem QA), basta pedir algo como "gera rápido". Para um mês específico, informe-o junto do comando:

```
/relatorio-mensal 03/2026
```

### O que o comando faz

| Etapa | Descrição                                                                         |
| ----- | --------------------------------------------------------------------------------- |
| 1     | Verifica o binário `jira-reporter` e compila se necessário                         |
| 2     | Determina o mês de competência (mês anterior, ou o informado via `-d`)             |
| 3     | Executa a CLI com as flags correspondentes às respostas                            |
| 4     | Extrai as chaves das issues do HTML gerado, preservando a ordem da tabela          |
| 5     | Consulta `/rest/api/3/search/jql` para obter `summary` e `description` completas   |
| 6     | Sintetiza um resumo curto Problema / Solução por issue e reescreve a seção no HTML |
| 7     | Informa o caminho do arquivo e um resumo do período e das atividades encontradas   |

> ℹ️ A API antiga `/rest/api/3/search` foi descontinuada pela Atlassian. O comando usa `/rest/api/3/search/jql`.

---

## ⏰ Automatizando com Cron

Você pode automatizar a geração do seu relatório Jira mensal usando `cron`.

1.  **Abra seu crontab para edição:**

    ```bash
    crontab -e
    ```

2.  **Adicione a seguinte linha ao seu crontab (usando o executável compilado):**

    ```cron
    0 10 1 * * cd /caminho/completo/para/seu/projeto/jira-reporter && ./jira-reporter >> ./jira-reporter.log 2>&1
    ```

3.  Salve e saia do seu editor crontab (geralmente pressionando Ctrl+X, depois S, depois Enter).

Seu relatório Jira será gerado automaticamente às 10h da manhã no primeiro dia de cada mês!

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👤 Autor

**Alan Gomes** - [GitHub](https://github.com/Alan-Gomes1)

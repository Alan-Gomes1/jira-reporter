# Jira Reporter

Este projeto é uma aplicação Go projetada para gerar um relatório de serviço mensal com base em tarefas do Jira. Ele busca issues do Jira do mês anterior que foram atribuídas ao usuário atual e então processa esses dados para gerar um relatório.

## Tecnologias Utilizadas

* **Go**: A linguagem de programação principal para a lógica da aplicação.
* **go-atlassian/v2**: Uma biblioteca Go usada para interagir com a API do Jira.
* **godotenv**: Uma biblioteca Go para gerenciar variáveis de ambiente.
* **Template HTML**: Usado para gerar a saída final do relatório.

## Como Baixar e Usar (Para Usuários)

Se você quer apenas usar a ferramenta sem precisar instalar o Go, siga os passos abaixo.

### 1. Baixe a Versão Correta

1.  Vá para a [página de Releases](https://github.com/Alan-Gomes1/jira-reporter/releases) do projeto.
2.  Encontre a versão mais recente.
3.  Na seção **Assets**, baixe o arquivo `.zip` correspondente ao seu sistema operacional:
    * **Windows (64-bit):** `jira-reporter-windows-amd64.zip`
    * **Linux (64-bit):** `jira-reporter-linux-amd64.zip`
    * **macOS (Intel):** `jira-reporter-macos-amd64.zip`
    * **macOS (Apple Silicon M1/M2/M3):** `jira-reporter-macos-arm64.zip`

### 2. Extraia e Configure

1.  Descompacte o arquivo `.zip` em uma pasta de sua preferência.
2.  Dentro da pasta, renomeie o arquivo `env-example` para `.env`.
3.  Abra o arquivo `.env` com um editor de texto e preencha com suas credenciais do Jira.

### 3. Execute o Relatório

Abra seu terminal (CMD ou PowerShell no Windows) na pasta onde você extraiu os arquivos e execute o comando correspondente:

#### **No Windows**
```powershell
.\jira-reporter.exe
````

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

-----

## Primeiros Passos (Para Desenvolvedores)

Siga estas instruções para configurar e executar o Jira Reporter a partir do código-fonte.

### Pré-requisitos

  * Go (versão 1.16 ou superior recomendada)

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

    Edite o arquivo `.env` e forneça os valores necessários para sua instância Jira.

3.  **Instale as Dependências:**

    ```bash
    go mod tidy
    ```

### Executando a Aplicação

Para executar a aplicação e gerar um relatório, execute o seguinte comando a partir da raiz do projeto:

```bash
go run main.go
```

Isso buscará as issues do Jira e gerará um relatório HTML no diretório `reports/`.

## Automatizando com Cron

Você pode automatizar a geração do seu relatório Jira mensal usando `cron`.

1.  **Abra seu crontab para edição:**

    ```bash
    crontab -e
    ```

2.  **Adicione a seguinte linha ao seu crontab (usando o executável compilado):**

    ```cron
    0 10 1 * * cd /caminho/completo/para/seu/projeto/jira-reporter && ./jira-reporter-linux-amd64 >> /caminho/completo/para/seu/projeto/jira-reporter.log 2>&1
    ```
3. Salve e saia do seu editor crontab (geralmente pressionando Ctrl+X, depois S, depois Enter).

Seu relatório Jira será gerado automaticamente às 10h da manhã no primeiro dia de cada mês!
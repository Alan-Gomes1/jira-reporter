// Package config fornece configuração centralizada para a aplicação.
// Implementa o padrão Singleton para carregar variáveis de ambiente uma única vez.
package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Config contém todas as configurações da aplicação carregadas do ambiente.
type Config struct {
	// Jira API configuration
	JiraURL   string
	JiraEmail string
	JiraToken string

	// User/Company configuration
	CompanyName string
	CNPJ        string
	Username    string
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

// Load carrega as configurações do arquivo .env e variáveis de ambiente.
// Utiliza sync.Once para garantir que o carregamento ocorra apenas uma vez.
func Load() (*Config, error) {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			// Não é um erro fatal se .env não existir,
			// as variáveis podem estar no ambiente
			if !os.IsNotExist(err) {
				loadErr = fmt.Errorf("erro ao carregar arquivo .env: %w", err)
				return
			}
		}

		instance = &Config{
			JiraURL:     getEnvOrDefault("URL", ""),
			JiraEmail:   getEnvOrDefault("EMAIL", ""),
			JiraToken:   getEnvOrDefault("API_KEY", ""),
			CompanyName: getEnvOrDefault("COMPANY_NAME", ""),
			CNPJ:        getEnvOrDefault("CNPJ", ""),
			Username:    getEnvOrDefault("USER_NAME", ""),
		}

		if err := instance.Validate(); err != nil {
			loadErr = err
			instance = nil
		}
	})

	return instance, loadErr
}

// Validate verifica se todas as configurações obrigatórias estão presentes.
func (c *Config) Validate() error {
	if c.JiraURL == "" {
		return fmt.Errorf("configuração obrigatória ausente: URL")
	}
	if c.JiraEmail == "" {
		return fmt.Errorf("configuração obrigatória ausente: EMAIL")
	}
	if c.JiraToken == "" {
		return fmt.Errorf("configuração obrigatória ausente: API_KEY")
	}
	if c.CompanyName == "" {
		return fmt.Errorf("configuração obrigatória ausente: COMPANY_NAME")
	}
	if c.CNPJ == "" {
		return fmt.Errorf("configuração obrigatória ausente: CNPJ")
	}
	if c.Username == "" {
		return fmt.Errorf("configuração obrigatória ausente: USER_NAME")
	}
	return nil
}

// getEnvOrDefault retorna o valor da variável de ambiente ou um valor padrão.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Reset limpa a instância singleton (útil para os testes).
func Reset() {
	once = sync.Once{}
	instance = nil
	loadErr = nil
}

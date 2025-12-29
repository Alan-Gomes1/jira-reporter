package view

import (
	"fmt"
	"html/template"
	"io"

	"github.com/alan-gomes1/jira-reporter/internal/model"
)

// htmlGenerator implementa ReportGenerator para formato HTML.
type htmlGenerator struct {
	templatePath string
}

// NewHTMLGenerator cria um novo gerador HTML.
func NewHTMLGenerator(templatePath string) ReportGenerator {
	return &htmlGenerator{
		templatePath: templatePath,
	}
}

// Generate gera o relatório em formato HTML.
func (g *htmlGenerator) Generate(
	writer io.Writer, data *model.ReportData, args ...string,
) error {
	if writer == nil {
		return fmt.Errorf("writer não pode ser nil para geração HTML")
	}

	templatePath := g.templatePath
	if len(args) > 0 && args[0] != "" {
		templatePath = args[0]
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("erro ao parsear template: %w", err)
	}

	if err := tmpl.Execute(writer, data); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}

	return nil
}

// Format retorna o formato suportado.
func (g *htmlGenerator) Format() model.ReportFormat {
	return model.FormatHTML
}

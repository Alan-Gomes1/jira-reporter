package view

import (
	"io"

	"github.com/alan-gomes1/jira-reporter/internal/model"
)

type ReportGenerator interface {
	// Generate gera o relatório no formato específico.
	// writer: destino da saída (pode ser nil para formatos que não usam)
	// data: dados do relatório
	// args: argumentos adicionais específicos do formato (ex: caminhos de arquivo)
	Generate(writer io.Writer, data *model.ReportData, args ...string) error

	// Format retorna o formato suportado pelo gerador.
	Format() model.ReportFormat
}

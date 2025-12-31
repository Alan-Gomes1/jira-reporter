/*
Copyright © 2025 Alan Gomes <alan.gomes.ag28@gmail.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/alan-gomes1/jira-reporter/internal/config"
	"github.com/alan-gomes1/jira-reporter/internal/model"
	"github.com/alan-gomes1/jira-reporter/internal/repository"
	"github.com/alan-gomes1/jira-reporter/internal/service"
	"github.com/alan-gomes1/jira-reporter/internal/view"
	"github.com/spf13/cobra"
)

const defaultTemplatePath = "template.html"

var rootCmd = &cobra.Command{
	Use:   "jira-reporter",
	Short: "Gera um relatório com dados do Jira",
	Long: `Gera um relatório mensal da Govone com dados do Jira, especificando
 quais tarefas foram realizadas no mês anterior.`,
	Run: runReport,
}

// runReport é o handler principal que orquestra a geração do relatório.
func runReport(cmd *cobra.Command, args []string) {
	// Obtem as flags
	reportName, _ := cmd.Flags().GetString("name")
	reportPath, _ := cmd.Flags().GetString("path")
	reportFormat, _ := cmd.Flags().GetString("format")
	reportDate, _ := cmd.Flags().GetString("date")
	includeQA, _ := cmd.Flags().GetBool("qa")

	// Carrega as configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Cria as dependências (Dependency Injection)
	reportService, err := buildReportService(cfg)
	if err != nil {
		log.Fatalf("Erro ao inicializar serviços: %v", err)
	}

	// Monta as opções do relatório
	opts := model.ReportOptions{
		Name:      reportName,
		Path:      reportPath,
		Format:    model.ReportFormat(reportFormat),
		Date:      reportDate,
		IncludeQA: includeQA,
	}

	// Gera o relatório
	if err := reportService.Generate(opts); err != nil {
		log.Fatalf("Erro ao gerar relatório: %v", err)
	}
}

// buildReportService constrói o ReportService com todas as dependências.
func buildReportService(cfg *config.Config) (service.ReportService, error) {
	// Repository
	jiraRepo, err := repository.NewJiraRepository(cfg)
	if err != nil {
		return nil, err
	}

	// Services
	dateService := service.NewDateService()
	fileService := service.NewFileService()

	// Generators (Views)
	generators := map[model.ReportFormat]view.ReportGenerator{
		model.FormatHTML: view.NewHTMLGenerator(defaultTemplatePath),
		model.FormatDOCX: view.NewDOCXGenerator(fileService),
	}

	reportService := service.NewReportService(
		cfg, jiraRepo, dateService, fileService, generators,
	)
	return reportService, nil
}

// Execute executa o comando raiz.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("name", "n", "", "Nome do relatório")
	rootCmd.Flags().StringP(
		"path", "p", "", "Caminho onde será salvo o relatório",
	)
	rootCmd.Flags().StringP(
		"format", "f", "html", "Formato do relatório (html ou docx)",
	)
	rootCmd.Flags().StringP(
		"date", "d", "",
		"Mês/ano do relatório no formato MM/YYYY (ex: 01/2025). "+
			"Padrão: mês anterior",
	)
	rootCmd.Flags().BoolP(
		"qa", "q", false,
		"Incluir cards onde o usuário está marcado como QA",
	)
}

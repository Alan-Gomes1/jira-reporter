/*
Copyright © 2025 Alan Gomes <alan.gomes.ag28@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/alan-gomes1/jira-reporter/internal/report"
	"github.com/spf13/cobra"
)

var (
	reportName string
	reportPath string
)

var rootCmd = &cobra.Command{
	Use:   "jira-reporter",
	Short: "Gera um relatório com dados do Jira",
	Long: `Gera um relatório mensal da Govone com dados do Jira, especificando
 quais tarefas foram realizadas no mês anterior.`,
	Run: func(cmd *cobra.Command, args []string) {
		report.Generate(reportName, reportPath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(
		&reportName, "name", "n", "", "Nome do relatório",
	)

	rootCmd.Flags().StringVarP(
		&reportPath, "path", "p", "", "Caminho onde será salvo o relatório",
	)
}

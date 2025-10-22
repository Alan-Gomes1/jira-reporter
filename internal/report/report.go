package report

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alan-gomes1/jira-reporter/internal/jira"
	"github.com/joho/godotenv"
)

func Generate() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	var (
		companyName = os.Getenv("COMPANY_NAME")
		cnpj        = os.Getenv("CNPJ")
		username    = os.Getenv("USERNAME")
	)

	userData := NewUserData(companyName, cnpj, username)
	jiraData, err := jira.GetJiraData()
	if err != nil {
		log.Fatalf("Error to generate report: %v", err)
	}

	now := time.Now()
	layout := "01/2006"
	dateWorked := now.Format(layout)
	reportData := ReportData{
		User:       *userData,
		Jira:       *jiraData,
		DateWorked: dateWorked,
	}

	day := time.Now().Day()
	monthAndYear := strings.Replace(dateWorked, "/", "_", 1)
	reportName := fmt.Sprintf("report_%d_%s.html", day, monthAndYear)
	reportPath := fmt.Sprintf("reports/%s", reportName)
	report, err := os.Create(reportPath)
	templatePath := "template.html"
	err = generateHTML(report, reportData, templatePath)
	if err != nil {
		log.Fatalf("Error to create report file: %v", err)
		return
	}
	defer report.Close()

	fmt.Printf("Relatório %s gerado com sucesso!\n", reportName)
}

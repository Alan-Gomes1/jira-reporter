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

func Generate(reportName, reportPath string) {
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
	previousMonth := now.AddDate(0, -1, 0)
	layout := "01/2006"
	dateWorked := previousMonth.Format(layout)
	reportData := ReportData{
		User:       *userData,
		Jira:       *jiraData,
		DateWorked: dateWorked,
	}

	day := time.Now().Day()
	monthAndYear := strings.Replace(dateWorked, "/", "_", 1)
	if reportName == "" {
		reportName = fmt.Sprintf("report_%d_%s.html", day, monthAndYear)
	} else {
		reportName = fmt.Sprintf(
			"%s_%d_%s.html", reportName, day, monthAndYear,
		)
	}

	if reportPath == "" {
		reportPath = fmt.Sprintf("reports/%s", reportName)
	} else {
		reportPath = fmt.Sprintf("%s/%s", reportPath, reportName)
	}

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

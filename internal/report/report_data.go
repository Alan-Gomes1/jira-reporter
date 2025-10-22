package report

import "github.com/alan-gomes1/jira-reporter/internal/jira"

type UserData struct {
	CompanyName string
	CNPJ        string
	Username    string
}

func NewUserData(companyName, cnpj, username string) *UserData {
	return &UserData{
		CompanyName: companyName,
		CNPJ:        cnpj,
		Username:    username,
	}
}

type ReportData struct {
	User       UserData
	Jira       jira.JiraData
	DateWorked string
}


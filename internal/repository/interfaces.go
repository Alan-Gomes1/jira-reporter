// Package repository define as interfaces e implementações para acesso a dados.
package repository

import (
	"time"

	"github.com/alan-gomes1/jira-reporter/internal/model"
)

// JiraRepository define a interface para acesso aos dados do Jira.
type JiraRepository interface {
	// FetchIssues busca issues do Jira no período especificado.
	FetchIssues(startDate, endDate time.Time) (*model.IssueCollection, error)
}

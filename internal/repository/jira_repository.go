package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/alan-gomes1/jira-reporter/internal/config"
	"github.com/alan-gomes1/jira-reporter/internal/model"
	jira "github.com/ctreminiom/go-atlassian/v2/jira/v3"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// Constantes de formatação de data
const (
	jiraDateFormat    = "2006-01-02"
	jiraTimeFormat    = "2006-01-02T15:04:05.999-0700"
	displayDateFormat = "02/01"
)

// jiraAPIRepository implementa JiraRepository usando a API do Jira.
type jiraAPIRepository struct {
	client *jira.Client
	config *config.Config
}

// NewJiraRepository cria uma nova instância do repositório Jira.
func NewJiraRepository(cfg *config.Config) (JiraRepository, error) {
	client, err := jira.New(nil, cfg.JiraURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente Jira: %w", err)
	}

	client.Auth.SetBasicAuth(cfg.JiraEmail, cfg.JiraToken)

	return &jiraAPIRepository{
		client: client,
		config: cfg,
	}, nil
}

// FetchIssues busca issues do Jira no período especificado.
// Se includeQA for true, também busca issues onde o usuário é QA.
func (r *jiraAPIRepository) FetchIssues(
	startDate, endDate time.Time, includeQA bool,
) (*model.IssueCollection, error) {
	jql := r.buildJQL(startDate, endDate, includeQA)
	fields := r.getRequiredFields()
	expand := []string{"changelog"}

	issues, response, err := r.client.Issue.Search.SearchJQL(
		context.Background(), jql, fields, expand, 100, "",
	)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf(
				"erro na busca de issues: %w - status: %s",
				err,
				response.Status,
			)
		}
		return nil, fmt.Errorf("erro na busca de issues: %w", err)
	}

	if issues == nil || len(issues.Issues) == 0 {
		return nil, fmt.Errorf("nenhuma issue encontrada no período")
	}

	collection := r.processIssues(issues)
	r.sortByDate(collection)

	return collection, nil
}

// buildJQL constrói a query JQL para buscar issues.
func (r *jiraAPIRepository) buildJQL(
	startDate, endDate time.Time, includeQA bool,
) string {
	firstDay := startDate.Format(jiraDateFormat)
	lastDay := endDate.Format(jiraDateFormat)

	// Condição base: assignee é o usuário atual
	baseCondition := fmt.Sprintf(
		"assignee = %s AND (%s during ('%s', '%s') OR %s >= '%s' AND %s <= '%s')",
		"currentUser()",
		"status changed to 'In Progress'",
		firstDay, lastDay,
		"created", firstDay,
		"created", lastDay,
	)

	// Se includeQA for true, adiciona condição para cards onde o usuário é QA
	if includeQA {
		qaCondition := fmt.Sprintf(
			"'QA[User Picker (single user)]' = %s AND (%s during ('%s', '%s') OR %s >= '%s' AND %s <= '%s')",
			"currentUser()",
			"status changed to 'In Progress'",
			firstDay, lastDay,
			"created", firstDay,
			"created", lastDay,
		)
		return fmt.Sprintf("(%s) OR (%s)", baseCondition, qaCondition)
	}

	return baseCondition
}

// getRequiredFields retorna os campos necessários para a busca.
func (r *jiraAPIRepository) getRequiredFields() []string {
	return []string{
		"key", "summary", "description", "status", "created", "assignee",
	}
}

// processIssues converte as issues da API para o modelo de domínio.
func (r *jiraAPIRepository) processIssues(
	issues *models.IssueSearchJQLScheme,
) *model.IssueCollection {
	collection := model.NewIssueCollection()

	for _, issue := range issues.Issues {
		issueDate := r.extractIssueDate(issue)
		description := r.extractDescription(issue)
		url := r.buildIssueURL(issue.Key)

		item := model.NewIssue(
			issue.Key,
			issue.Fields.Summary,
			description,
			issueDate,
			url,
		)
		collection.Add(*item)
	}

	return collection
}

// extractIssueDate extrai a data relevante da issue (In Progress, atribuição ou criação).
func (r *jiraAPIRepository) extractIssueDate(issue *models.IssueScheme) string {
	inProgressDate := r.findInProgressDate(issue)
	if inProgressDate != "" {
		return inProgressDate
	}

	assigneeDate := r.findAssigneeDate(issue)
	if assigneeDate != "" {
		return assigneeDate
	}

	return r.parseCreatedDate(issue)
}

// findInProgressDate busca a data em que a issue entrou em "In Progress".
func (r *jiraAPIRepository) findInProgressDate(
	issue *models.IssueScheme,
) string {
	if issue.Changelog == nil {
		return ""
	}

	for _, history := range issue.Changelog.Histories {
		for _, item := range history.Items {
			if item.Field == "status" && item.ToString == "In Progress" {
				if date := r.parseJiraTime(history.Created); date != "" {
					return date
				}
			}
		}
	}
	return ""
}

// findAssigneeDate busca a data em que a issue foi atribuída ao usuário atual.
func (r *jiraAPIRepository) findAssigneeDate(issue *models.IssueScheme) string {
	noChangelog := issue.Changelog == nil
	noFields := issue.Fields == nil
	noAssignee := issue.Fields.Assignee == nil
	if noChangelog || noFields || noAssignee {
		return ""
	}

	for _, history := range issue.Changelog.Histories {
		for _, item := range history.Items {
			assignee := item.Field == "assignee"
			if assignee && item.To == issue.Fields.Assignee.AccountID {
				if date := r.parseJiraTime(history.Created); date != "" {
					return date
				}
			}
		}
	}
	return ""
}

// parseCreatedDate extrai a data de criação da issue.
func (r *jiraAPIRepository) parseCreatedDate(issue *models.IssueScheme) string {
	if issue.Fields == nil || issue.Fields.Created == nil {
		return ""
	}
	created := fmt.Sprintf("%v", issue.Fields.Created)
	return r.parseJiraTime(created)
}

// parseJiraTime converte uma string de data do Jira para o formato de exibição.
func (r *jiraAPIRepository) parseJiraTime(dateStr string) string {
	parsedDate, err := time.Parse(jiraTimeFormat, dateStr)
	if err != nil {
		return ""
	}
	return parsedDate.Format(displayDateFormat)
}

// extractDescription extrai o texto da descrição da issue.
func (r *jiraAPIRepository) extractDescription(
	issue *models.IssueScheme,
) string {
	if issue.Fields == nil {
		return ""
	}
	desc := issue.Fields.Description
	if desc == nil || len(desc.Content) == 0 {
		return ""
	}
	if desc.Content[0] == nil || len(desc.Content[0].Content) == 0 {
		return ""
	}
	if desc.Content[0].Content[0] == nil {
		return ""
	}
	return desc.Content[0].Content[0].Text
}

// buildIssueURL constrói a URL da issue.
func (r *jiraAPIRepository) buildIssueURL(key string) string {
	return fmt.Sprintf("%s/browse/%s", r.config.JiraURL, key)
}

// sortByDate ordena a coleção de issues por data.
func (r *jiraAPIRepository) sortByDate(collection *model.IssueCollection) {
	sort.Slice(collection.Items, func(i, j int) bool {
		return collection.Items[i].Date < collection.Items[j].Date
	})
}

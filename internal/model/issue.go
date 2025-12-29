package model

// Issue representa uma issue do Jira com os dados necessários para o relatório.
type Issue struct {
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Date        string `json:"date"`
	URL         string `json:"url"`
}

// NewIssue cria uma nova instância de Issue.
func NewIssue(key, summary, description, date, url string) *Issue {
	return &Issue{
		Key:         key,
		Summary:     summary,
		Description: description,
		Date:        date,
		URL:         url,
	}
}

// IssueCollection representa uma coleção de issues.
type IssueCollection struct {
	Items []Issue
}

// NewIssueCollection cria uma nova coleção vazia de issues.
func NewIssueCollection() *IssueCollection {
	return &IssueCollection{
		Items: make([]Issue, 0),
	}
}

// Add adiciona uma issue à coleção.
func (c *IssueCollection) Add(issue Issue) {
	c.Items = append(c.Items, issue)
}

// Count retorna o número de issues na coleção.
func (c *IssueCollection) Count() int {
	return len(c.Items)
}

// IsEmpty verifica se a coleção está vazia.
func (c *IssueCollection) IsEmpty() bool {
	return len(c.Items) == 0
}

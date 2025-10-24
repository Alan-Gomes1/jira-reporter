package jira

type JiraItem struct {
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Date        string `json:"date"`
	URL         string `json:"url"`
}

func NewJiraItem(key, summary, description, date, url string) *JiraItem {
	return &JiraItem{
		Key:         key,
		Summary:     summary,
		Description: description,
		Date:        date,
		URL:         url,
	}
}

type JiraData struct {
	Items []JiraItem
}

func NewJiraData() *JiraData {
	return &JiraData{
		Items: make([]JiraItem, 0),
	}
}

func (rd *JiraData) AddItem(item JiraItem) {
	rd.Items = append(rd.Items, item)
}

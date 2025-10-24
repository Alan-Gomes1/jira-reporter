package jira

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	jira "github.com/ctreminiom/go-atlassian/v2/jira/v3"
	"github.com/joho/godotenv"
)

func GetJiraData() (*JiraData, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	var (
		host  = os.Getenv("URL")
		mail  = os.Getenv("EMAIL")
		token = os.Getenv("API_KEY")
	)

	client, err := jira.New(nil, host)
	if err != nil {
		log.Fatalf("Error to create new JIRA client: %v", err)
		return nil, err
	}

	client.Auth.SetBasicAuth(mail, token)
	dateFormat := "2006-01-02"
	year := time.Now().Year()
	month := time.Now().Month() - 1
	// Se o mês atual for Janeiro, o anterior é Dezembro do ano passado
	if month == 0 {
		month = 12
		year--
	}

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	firstDayOfMonth := firstDay.Format(dateFormat)
	lastDayOfMonth := lastDay.Format(dateFormat)
	jql := fmt.Sprintf(
		"assignee = %s AND (%s during ('%s', '%s') OR %s >= '%s' AND %s <= '%s')",
		"currentUser()", "status changed to 'In Progress'", firstDayOfMonth,
		lastDayOfMonth, "created", firstDayOfMonth, "created", lastDayOfMonth,
	)
	fields := []string{
		"key", "summary", "description", "status", "created", "assignee",
	}
	expand := []string{"changelog"}
	issues, response, err := client.Issue.Search.SearchJQL(
		context.Background(), jql, fields, expand, 100, "",
	)
	if err != nil {
		if response != nil {
			log.Fatalf("Error to search issues: %v - %s", err, response.Status)
			return nil, err
		}
		log.Fatalf("Error to search issues: %v", err)
		return nil, err
	}

	if issues == nil || len(issues.Issues) == 0 {
		log.Fatalf("Nenhum card encontrado.")
		return nil, err
	}

	reportData := NewJiraData()
	for _, issue := range issues.Issues {
		var descriptionText, assigneeDate, issueDate string
		description := issue.Fields.Description
		foundInProgressDate := false

		// Verifica se o changelog foi retornado
		if issue.Changelog != nil && len(issue.Changelog.Histories) > 0 {
			for _, history := range issue.Changelog.Histories {
				for _, item := range history.Items {
					if item.Field == "status" &&
						item.ToString == "In Progress" &&
						!foundInProgressDate {
						parsedDate, err := time.Parse(
							"2006-01-02T15:04:05.999-0700", history.Created,
						)
						if err == nil {
							issueDate = parsedDate.Format("02/01")
							foundInProgressDate = true
						}
					}

					assignee := issue.Fields.Assignee
					if item.Field == "assignee" &&
						assignee != nil &&
						item.To == assignee.AccountID {
						parsedDate, err := time.Parse(
							"2006-01-02T15:04:05.999-0700", history.Created,
						)
						if err == nil {
							assigneeDate = parsedDate.Format("02/01")
						}
					}
				}
			}
		}

		// Se não encontrou a data de "In Progress"
		// usa a data de atribuição do card
		if !foundInProgressDate {
			if assigneeDate != "" {
				issueDate = assigneeDate
			} else {
				created := fmt.Sprintf("%v", issue.Fields.Created)
				parsedDate, err := time.Parse(
					"2006-01-02T15:04:05.999-0700", created,
				)
				if err == nil {
					issueDate = parsedDate.Format("02/01")
				}
			}
		}

		// Tenta pegar a descrição do card
		if description != nil && len(description.Content) > 0 &&
			description.Content[0] != nil &&
			len(description.Content[0].Content) > 0 &&
			description.Content[0].Content[0] != nil {
			descriptionText = description.Content[0].Content[0].Text
		}

		issueURL := fmt.Sprintf("%s/browse/%s", host, issue.Key)
		item := NewJiraItem(
			issue.Key,
			issue.Fields.Summary,
			descriptionText,
			issueDate,
			issueURL,
		)
		reportData.AddItem(*item)
	}
	sort.Slice(reportData.Items, func(i, j int) bool {
		return reportData.Items[i].Date < reportData.Items[j].Date
	})
	return reportData, nil
}

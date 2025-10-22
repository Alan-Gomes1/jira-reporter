package report

import (
	"fmt"
	"html/template"
	"io"
)

func generateHTML(writer io.Writer, data ReportData, templatePath string) error {
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("Error to parse template: %v", err)
	}

	err = temp.Execute(writer, data)
	if err != nil {
		return fmt.Errorf("Error to execute template: %v", err)
	}

	return nil
}

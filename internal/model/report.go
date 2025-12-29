package model

// ReportData agrega todos os dados necessários para gerar um relatório.
type ReportData struct {
	User       User
	Jira       IssueCollection
	DateWorked string
}

// NewReportData cria uma nova instância de ReportData.
func NewReportData(
	user User, issues IssueCollection, dateWorked string,
) *ReportData {
	return &ReportData{
		User:       user,
		Jira:       issues,
		DateWorked: dateWorked,
	}
}

// ReportFormat representa os formatos de saída suportados.
type ReportFormat string

const (
	FormatHTML ReportFormat = "html"
	FormatDOCX ReportFormat = "docx"
)

// IsValid verifica se o formato é válido.
func (f ReportFormat) IsValid() bool {
	return f == FormatHTML || f == FormatDOCX
}

// String retorna a representação string do formato.
func (f ReportFormat) String() string {
	return string(f)
}

// Extension retorna a extensão de arquivo para o formato.
func (f ReportFormat) Extension() string {
	return string(f)
}

// ReportOptions contém as opções para geração de relatório.
type ReportOptions struct {
	Name   string
	Path   string
	Format ReportFormat
	Date   string // Mês/ano no formato MM/YYYY (opcional, padrão: mês anterior)
}

// NewReportOptions cria opções com valores padrão.
func NewReportOptions() *ReportOptions {
	return &ReportOptions{
		Name:   "",
		Path:   "reports",
		Format: FormatHTML,
		Date:   "", // Vazio significa mês anterior
	}
}

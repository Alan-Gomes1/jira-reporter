// Package service contém a lógica de negócio da aplicação.
package service

import "time"

// DateService fornece operações relacionadas a datas.
type DateService interface {
	// GetPreviousMonthRange retorna o primeiro e último dia do mês anterior.
	GetPreviousMonthRange() (firstDay, lastDay time.Time)
	// FormatDateWorked formata a data trabalhada para exibição (MM/YYYY).
	FormatDateWorked(date time.Time) string
}

// dateService implementa DateService.
type dateService struct{}

// NewDateService cria uma nova instância de DateService.
func NewDateService() DateService {
	return &dateService{}
}

// GetPreviousMonthRange retorna o primeiro e último dia do mês anterior.
func (s *dateService) GetPreviousMonthRange() (time.Time, time.Time) {
	now := time.Now()
	year := now.Year()
	month := now.Month() - 1

	// Se o mês atual for Janeiro, o anterior é Dezembro do ano passado
	if month == 0 {
		month = 12
		year--
	}

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay
}

// FormatDateWorked formata a data para o padrão MM/YYYY.
func (s *dateService) FormatDateWorked(date time.Time) string {
	return date.Format("01/2006")
}

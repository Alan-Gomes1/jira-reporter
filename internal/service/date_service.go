// Package service contém a lógica de negócio da aplicação.
package service

import (
	"fmt"
	"time"
)

// DateService fornece operações relacionadas a datas.
type DateService interface {
	// GetPreviousMonthRange retorna o primeiro e último dia do mês anterior.
	GetPreviousMonthRange() (firstDay, lastDay time.Time)
	// GetMonthRange retorna o primeiro e último dia de um mês/ano específico.
	GetMonthRange(month, year int) (firstDay, lastDay time.Time)
	// ParseMonthYear converte uma string MM/YYYY em mês e ano.
	ParseMonthYear(date string) (month, year int, err error)
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

// GetMonthRange retorna o primeiro e último dia de um mês/ano específico.
func (s *dateService) GetMonthRange(month, year int) (time.Time, time.Time) {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay
}

// ParseMonthYear converte uma string MM/YYYY em mês e ano.
func (s *dateService) ParseMonthYear(date string) (int, int, error) {
	parsedDate, err := time.Parse("01/2006", date)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"formato de data inválido '%s': use MM/YYYY (ex: 01/2025)", date,
		)
	}
	return int(parsedDate.Month()), parsedDate.Year(), nil
}

// FormatDateWorked formata a data para o padrão MM/YYYY.
func (s *dateService) FormatDateWorked(date time.Time) string {
	return date.Format("01/2006")
}

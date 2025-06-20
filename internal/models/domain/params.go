package domain

import (
	"avito_pvz/internal/http/gen"
	"time"
)

type Params struct {
	StartDate *time.Time

	// EndDate Конечная дата диапазона
	EndDate *time.Time

	// Page Номер страницы
	Page *int

	// Limit Количество элементов на странице
	Limit *int
}

func NewParamsFromDTO(p gen.GetPvzParams) *Params {
	return &Params{
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		Page:      p.Page,
		Limit:     p.Limit,
	}
}

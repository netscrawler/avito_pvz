package domain

import (
	"time"

	"avito_pvz/internal/http/gen"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

const (
	Kazan       PvzCity = "Казань"
	Moscow      PvzCity = "Москва"
	StPeterburg PvzCity = "Санкт-Петербург"
)

func NewPVZ(city PvzCity) *PVZ {
	uid := uuid.New()
	return &PVZ{
		ID:               (*PVZID)(&uid),
		City:             city,
		RegistrationDate: time.Now(),
	}
}

type PvzCity string

func (p PvzCity) IsValid() bool {
	return p == Kazan || p == Moscow || p == StPeterburg
}

type PVZ struct {
	ID               *PVZID
	RegistrationDate time.Time
	City             PvzCity
}

func (p *PVZ) ToDTO() gen.PVZ {
	return gen.PVZ{
		City:             gen.PVZCity(p.City),
		Id:               (*types.UUID)(p.ID),
		RegistrationDate: &p.RegistrationDate,
	}
}

type PVZID uuid.UUID

func (p *PVZID) String() string {
	u := uuid.UUID(*p)

	return u.String()
}

type PVZList []*PVZ

type PVZAgregate struct {
	Pvz        *PVZ
	Receptions *[]struct {
		Products  *[]Product
		Reception *Reception
	}
}

func AggregateToPvzResponse(aggregates []PVZAgregate) gen.GetPvz200JSONResponse {
	var response gen.GetPvz200JSONResponse

	for _, aggregate := range aggregates {
		var receptions []struct {
			Products  *[]gen.Product `json:"products,omitempty"`  // Исправлено на gen.Product
			Reception *gen.Reception `json:"reception,omitempty"` // Исправлено на gen.Reception
		}

		// Перебираем все приемки для текущего ПВЗ
		for _, receptionDetail := range *aggregate.Receptions {
			// Приводим продукты и приемку к нужному типу, если это необходимо
			products := convertToGenProducts(
				receptionDetail.Products,
			) // Конвертируем в *[]gen.Product
			reception := convertToGenReception(
				receptionDetail.Reception,
			) // Конвертируем в *gen.Reception

			receptions = append(receptions, struct {
				Products  *[]gen.Product `json:"products,omitempty"`  // Исправлено на gen.Product
				Reception *gen.Reception `json:"reception,omitempty"` // Исправлено на gen.Reception
			}{
				Products:  &products,
				Reception: reception,
			})
		}

		ag := aggregate.Pvz.ToDTO()

		// Добавляем объект для текущего ПВЗ в ответ
		response = append(response, struct {
			Pvz        *gen.PVZ `json:"pvz,omitempty"` // Исправлено на gen.PVZ
			Receptions *[]struct {
				Products  *[]gen.Product `json:"products,omitempty"`  // Исправлено на gen.Product
				Reception *gen.Reception `json:"reception,omitempty"` // Исправлено на gen.Reception
			} `json:"receptions,omitempty"` // Исправлено на gen.Reception
		}{
			Pvz:        &ag, // Это может быть типом *gen.PVZ, если вы используете gen.PVZ
			Receptions: &receptions,
		})
	}

	return response
}

func convertToGenProducts(products *[]Product) []gen.Product {
	var genProducts []gen.Product
	if products != nil {
		genProducts = make([]gen.Product, len(*products))
		for i, product := range *products {
			genProducts[i] = product.ToDto()
		}
	}
	return genProducts
}

func convertToGenReception(reception *Reception) *gen.Reception {
	if reception == nil {
		return nil
	}

	r := reception.ToDTO()

	return &r
}

func (p PVZList) ToDTO() []gen.PVZ {
	response := make([]gen.PVZ, 0, len(p))
	for _, pvz := range p {
		response = append(response, pvz.ToDTO())
	}

	return response
}

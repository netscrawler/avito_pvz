package domain

import (
	"avito_pvz/internal/http/gen"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

type ProductType string

func (p ProductType) IsValid() bool {
	return ((p == ProductTypeClothing || p == ProductTypeElectronics) || p == ProductTypeShoes)
}

const (
	ProductTypeElectronics ProductType = "электроника"
	ProductTypeClothing    ProductType = "одежда"
	ProductTypeShoes       ProductType = "обувь"
)

type Product struct {
	ID          uuid.UUID
	ReceptionID uuid.UUID
	Type        ProductType
	CreatedAt   time.Time
}

func (p *Product) ToDto() gen.Product {
	return gen.Product{
		DateTime:    &p.CreatedAt,
		Id:          &p.ID,
		ReceptionId: types.UUID(p.ReceptionID),
		Type:        gen.ProductType(p.Type),
	}
}

type ProductToAdd struct {
	UUID PVZID
	Type ProductType
}

func NewProduct(intake uuid.UUID, productType ProductType) *Product {
	return &Product{
		ID:          uuid.New(),
		ReceptionID: intake,
		Type:        productType,
		CreatedAt:   time.Now(),
	}
}

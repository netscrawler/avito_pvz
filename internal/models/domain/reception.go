package domain

import (
	"time"

	"avito_pvz/internal/http/gen"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

type ReceptionStatus string

const (
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
	ReceptionStatusClosed     ReceptionStatus = "close"
)

type Reception struct {
	ID        uuid.UUID
	PvzID     uuid.UUID
	Status    ReceptionStatus
	CreatedAt time.Time
}

func (r *Reception) Close() {
	r.Status = ReceptionStatusClosed
}

func NewReception(pvz uuid.UUID) *Reception {
	return &Reception{
		ID:        uuid.New(),
		PvzID:     pvz,
		Status:    ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}
}

func (r *Reception) IsActive() bool {
	return r.Status == ReceptionStatusInProgress
}

func (r Reception) ToDTO() gen.Reception {
	return gen.Reception{
		DateTime: r.CreatedAt,
		Id:       (*types.UUID)(&r.ID),
		PvzId:    (types.UUID)(r.PvzID),
		Status:   gen.ReceptionStatus(r.Status),
	}
}

type ReceptionID uuid.UUID

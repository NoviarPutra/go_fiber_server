package types

import "time"

type DivisionRow struct {
	ID        string     `json:"id"`
	CompanyID string     `json:"company_id"`
	Name      string     `json:"name"`
	Code      *string    `json:"code"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateDivisionRequest struct {
	CompanyID string  `json:"company_id" validate:"required,uuid4"`
	Name      string  `json:"name" validate:"required,min=1"`
	Code      *string `json:"code" validate:"omitempty"`
}

type UpdateDivisionRequest struct {
	Name *string `json:"name" validate:"omitempty,min=1"`
	Code *string `json:"code" validate:"omitempty"`
}

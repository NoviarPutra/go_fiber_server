package types

import "time"

type CompanyBranchRow struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Name      string    `json:"name"`
	Address   *string   `json:"address"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCompanyBranchRequest struct {
	CompanyID string  `json:"company_id" validate:"required,uuid"`
	Name      string  `json:"name" validate:"required,min=1"`
	Address   *string `json:"address"`
	Timezone  string  `json:"timezone" validate:"required"`
}

type UpdateCompanyBranchRequest struct {
	CompanyID *string `json:"company_id" validate:"omitempty,uuid"`
	Name      *string `json:"name" validate:"omitempty,min=1"`
	Address   *string `json:"address"`
	Timezone  *string `json:"timezone" validate:"omitempty"`
}

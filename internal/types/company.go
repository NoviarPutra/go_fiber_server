package types

import "time"

type CompanyRow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	LogoUrl   *string   `json:"logo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCompanyRequest struct {
	Name    string  `json:"name" validate:"required,min=1"`
	Code    string  `json:"code" validate:"required"`
	LogoUrl *string `json:"logo_url"`
}

type UpdateCompanyRequest struct {
	Name    *string `json:"name" validate:"omitempty,min=1"`
	Code    *string `json:"code" validate:"omitempty"`
	LogoUrl *string `json:"logo_url"`
}

package types

import (
	"time"

	"github.com/google/uuid"
)

type CompanyUser struct {
	ID           uuid.UUID  `json:"id"`
	CompanyID    uuid.UUID  `json:"company_id"`
	UserID       uuid.UUID  `json:"user_id"`
	BranchID     *uuid.UUID `json:"branch_id"`
	EmployeeCode *string    `json:"employee_code"`
	IsOwner      bool       `json:"is_owner"`
	IsActive     bool       `json:"is_active"`
	JoinedAt     time.Time  `json:"joined_at"`
	LeftAt       *time.Time `json:"left_at"`
	CreatedAt    time.Time  `json:"created_at"` // maps to joined_at conceptually or joined_at is separate
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	// Expanded fields for display
	UserName  *string `json:"user_name,omitempty"`
	UserEmail *string `json:"user_email,omitempty"`
}

type CreateCompanyUserRequest struct {
	UserID       uuid.UUID  `json:"user_id" validate:"required"`
	BranchID     *uuid.UUID `json:"branch_id"`
	EmployeeCode *string    `json:"employee_code"`
	IsOwner      bool       `json:"is_owner"`
}

type UpdateCompanyUserRequest struct {
	BranchID     *uuid.UUID `json:"branch_id"`
	EmployeeCode *string    `json:"employee_code"`
	IsOwner      *bool      `json:"is_owner"`
	IsActive     *bool      `json:"is_active"`
}

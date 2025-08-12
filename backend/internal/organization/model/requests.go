package organizationmodel

type CreateOrganizationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	Website     string `json:"website"`
}

type UpdateOrganizationRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Website     *string `json:"website"`
	IsActive    *bool   `json:"is_active"`
}
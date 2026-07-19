package request

// CreateTenant is the input for provisioning a new tenant.
type CreateTenant struct {
	Name   string  `json:"name" validate:"required,min=2,max=255"`
	Slug   string  `json:"slug" validate:"required,min=2,max=100,alphanum"`
	Domain *string `json:"domain" validate:"omitempty,fqdn"`
	Plan   string  `json:"plan" validate:"required,oneof=free starter pro enterprise"`
}

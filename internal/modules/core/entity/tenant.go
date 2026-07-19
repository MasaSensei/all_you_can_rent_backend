package entity

import "time"

// Tenant mirrors the tenants table exactly. It is persistence-shaped, not
// API-shaped — handlers map it to dto/response.Tenant before returning it.
type Tenant struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Slug      string     `db:"slug"`
	Domain    *string    `db:"domain"`
	Plan      string     `db:"plan"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

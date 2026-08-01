package entity

import "time"

// Expense mirrors the expenses table.
type Expense struct {
	ID          string     `db:"id"`
	TenantID    string     `db:"tenant_id"`
	AssetID     *string    `db:"asset_id"`
	Category    string     `db:"category"`
	Amount      float64    `db:"amount"`
	ExpenseDate time.Time  `db:"expense_date"`
	Description *string    `db:"description"`
	Vendor      *string    `db:"vendor"`
	Status      string     `db:"status"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedBy   *string    `db:"deleted_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}

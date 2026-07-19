package entity

import "time"

// SystemSetting mirrors the system_settings table (global, tenant-less).
type SystemSetting struct {
	ID        string     `db:"id"`
	Key       string     `db:"key"`
	Value     string     `db:"value"`
	Type      string     `db:"type"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

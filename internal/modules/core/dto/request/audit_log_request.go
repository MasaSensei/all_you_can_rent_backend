package request

// CreateAuditLog is the input used internally by other modules' services
// to record an audit entry. It is not exposed over HTTP — audit writes
// happen as a side effect of other modules' mutations, never directly
// from a client.
type CreateAuditLog struct {
	UserID     *string `json:"user_id"`
	EntityType string  `json:"entity_type" validate:"required,max=100"`
	EntityID   string  `json:"entity_id" validate:"required,uuid"`
	Action     string  `json:"action" validate:"required,oneof=create update delete login logout"`
	OldValues  *string `json:"old_values"`
	NewValues  *string `json:"new_values"`
	IPAddress  *string `json:"ip_address"`
	UserAgent  *string `json:"user_agent"`
}

// ListAuditLogsFilter holds the whitelisted query-param filters for
// listing audit logs. Built by the handler, consumed by the service.
type ListAuditLogsFilter struct {
	EntityType *string
	EntityID   *string
	Action     *string
	Page       int
	PerPage    int
}

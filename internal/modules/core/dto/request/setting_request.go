package request

// UpsertSetting is the input for creating/updating a per-tenant setting.
type UpsertSetting struct {
	Key   string `json:"key" validate:"required,max=100"`
	Value string `json:"value" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=string number boolean json"`
}

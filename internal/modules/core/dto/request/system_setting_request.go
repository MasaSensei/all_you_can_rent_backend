package request

// UpsertSystemSetting is the input for creating/updating a global setting.
type UpsertSystemSetting struct {
	Key   string `json:"key" validate:"required,max=100"`
	Value string `json:"value" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=string number boolean json"`
}

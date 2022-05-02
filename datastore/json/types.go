package json

type JsonVersion struct {
	Version string `json:"version"`
	Schema  string `json:"schema"`
}

// Implement common.Version
func (v *JsonVersion) String() string {
	return v.Version
}
func (v *JsonVersion) SchemaName() string {
	return v.Schema
}

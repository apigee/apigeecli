package types

type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ImportError struct {
	Err error
}

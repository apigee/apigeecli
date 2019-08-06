package types

//Attributes to used to hold custom attributes for entities
type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Arguments is the base struct to hold all command arguments
type Arguments struct {
	Org            string
	Env            string
	Token          string
	ServiceAccount string
}

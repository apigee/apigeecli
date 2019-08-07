package types

//Attributes to used to hold custom attributes for entities
type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Arguments is the base struct to hold all command arguments
type Arguments struct {
	Org            string //Apigee org
	Env            string //Apigee environment
	Token          string //Google OAuth access token
	ServiceAccount string //Google service account json
	AliasName      string //AliasName for the key store
	LogInfo        bool   //LogInfo controls the log level
	SkipCheck      bool   //skip checking access token expiry
	SkipCache      bool   //skip writing access token to file
}

// Structure to hold OAuth response
type OAuthAccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

type KeyAliasName string

func (a KeyAliasName) String() string {
	return string(a)
}

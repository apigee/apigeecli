// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bundlegen

var generateSetTarget bool

var securitySchemesList = securitySchemesListDef{}

var (
	quotaPolicyContent       = map[string]string{}
	spikeArrestPolicyContent = map[string]string{}
)

type pathDetailDef struct {
	OperationID    string
	Description    string
	Path           string
	Verb           string
	AssignMessage  string
	Backend        backendDef
	SecurityScheme securitySchemesDef
	SpikeArrest    []spikeArrestDef
	Quota          []quotaDef
}

type spikeArrestDef struct {
	SpikeArrestEnabled       bool
	SpikeArrestName          string
	SpikeArrestIdentifierRef string
	SpikeArrestRateRef       string
	SpikeArrestRateLiteral   string
}

type quotaDef struct {
	QuotaEnabled          bool
	QuotaName             string
	QuotaAllowRef         string
	QuotaAllowLiteral     string
	QuotaIntervalRef      string
	QuotaIntervalLiteral  string
	QuotaTimeUnitRef      string
	QuotaTimeUnitLiteral  string
	QuotaConfigStepName   string
	QuotaIdentifierRef    string
	QuotaIdentiferLiteral string
}

type securitySchemesListDef struct {
	SecuritySchemes []securitySchemesDef
}

type securitySchemesDef struct {
	SchemeName   string
	OAuthPolicy  oAuthPolicyDef
	APIKeyPolicy apiKeyPolicyDef
	JWTPolicy    jwtPolicyDef
}

type oAuthPolicyDef struct {
	OAuthPolicyEnabled bool
	Scope              string
}

type apiKeyPolicyDef struct {
	APIKeyPolicyEnabled bool
	APIKeyLocation      string
	APIKeyName          string
}

type jwtPolicyDef struct {
	JWTPolicyEnabled bool
	JwkUri           string
	Issuer           string
	Audience         string
	Source           string
	Location         map[string]string // only one location supported for now
}

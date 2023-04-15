// Copyright 2021 Google LLC
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

package developers

import (
	"encoding/json"
	"net/url"
	"path"

	"internal/apiclient"
)

type developerAdjustment struct {
	Adjustment money `json:"adjustment,omitempty"`
}

type money struct {
	CurrencyCode string `json:"currencyCode,omitempty"`
	Units        string `json:"units,omitempty"`
	Nanos        int64  `json:"nanos,omitempty"`
}

type transaction struct {
	TransactionAmount money  `json:"transactionAmount,omitempty"`
	TransactionId     string `json:"transactionId,omitempty"`
}

func Adjust(email string, adjust string) (respBody []byte, err error) {
	dAdjustment := developerAdjustment{}
	if err = json.Unmarshal([]byte(adjust), &dAdjustment); err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "balance:adjust")
	respBody, err = apiclient.HttpClient(u.String(), adjust)
	return respBody, err
}

func Credit(email string, transact string) (respBody []byte, err error) {
	txn := transaction{}
	if err = json.Unmarshal([]byte(transact), &txn); err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "developers", url.QueryEscape(email), "balance:credit")
	respBody, err = apiclient.HttpClient(u.String(), transact)
	return respBody, err
}

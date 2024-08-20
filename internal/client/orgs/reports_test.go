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

package orgs

import (
	"internal/client/clienttest"
	"testing"
	"time"
)

func TestTotalAPICallsInMonth(t *testing.T) {
	if err := clienttest.TestSetup(clienttest.ENV_NOT_REQD,
		clienttest.SITEID_NOT_REQD, clienttest.CLIPATH_NOT_REQD); err != nil {
		t.Fatalf("%v", err)
	}
	// Get the current time
	currentTime := time.Now()

	// Get the current month as an integer (1-12)
	currentMonth := int(currentTime.Month())
	currentYear := int(currentTime.Year())

	if _, _, _, err := TotalAPICallsInMonth(currentMonth, currentYear, true, false, "PAYG"); err != nil {
		t.Fatalf("%v", err)
	}
}

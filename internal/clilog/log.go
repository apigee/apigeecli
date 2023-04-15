// Copyright 2020 Google LLC
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

package clilog

import (
	"io"
	"log"
	"os"
)

// log levels, default is error
var (
	Debug        *log.Logger
	Info         *log.Logger
	Warning      *log.Logger
	Error        *log.Logger
	HttpResponse *log.Logger
	HttpError    *log.Logger
)

// Init function initializes the logger objects
func Init(debug bool, print bool, NoOutput bool) {
	debugHandle := io.Discard
	infoHandle := io.Discard
	var warningHandle, errorHandle, responseHandle io.Writer

	if debug {
		debugHandle = os.Stdout
	}

	if print {
		infoHandle = os.Stdout
	}

	if NoOutput {
		responseHandle = io.Discard
		infoHandle = io.Discard
		errorHandle = io.Discard
		warningHandle = io.Discard
	} else {
		responseHandle = os.Stdout
		warningHandle = os.Stdout
		errorHandle = os.Stderr
	}

	Debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"", 0)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	HttpResponse = log.New(responseHandle,
		"", 0)

	HttpError = log.New(errorHandle,
		"", 0)
}

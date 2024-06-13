# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

ADD ../internal /go/src/apigeecli/internal
ADD ../cmd /go/src/apigeecli/cmd

COPY ../go.mod ../go.sum /go/src/apigeecli/

WORKDIR /go/src/apigeecli

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod download


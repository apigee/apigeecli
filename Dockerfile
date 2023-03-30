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

FROM golang:1.20 as builder

ADD ./internal /go/src/apigeecli/internal
ADD ./cmd /go/src/apigeecli/cmd

COPY main.go /go/src/apigeecli/main.go
COPY go.mod go.sum /go/src/apigeecli/

WORKDIR /go/src/apigeecli

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -a -ldflags='-s -w -extldflags "-static"' -o /go/bin/apigeecli /go/src/apigeecli/main.go

FROM gcr.io/distroless/static-debian11
COPY --from=builder /go/bin/apigeecli /
COPY LICENSE.txt /
COPY third-party-licenses.txt /
CMD ["/apigeecli"]

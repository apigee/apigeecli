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

ARG TAG
ARG COMMIT

ADD ./internal /go/src/apigeecli/internal
ADD ./cmd /go/src/apigeecli/cmd

COPY main.go /go/src/apigeecli/main.go
COPY go.mod go.sum /go/src/apigeecli/

WORKDIR /go/src/apigeecli

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod download
RUN date +%FT%H:%I:%M+%Z > /tmp/date
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=true -a -gcflags='all="-l"' -ldflags='-s -w -extldflags "-static" -X main.version='${TAG}' -X main.commit='${COMMIT}' -X main.date='$(cat /tmp/date) -o /go/bin/apigeecli /go/src/apigeecli/main.go

FROM gcr.io/cloud-marketplace/google/debian11:latest
COPY --from=builder /go/bin/apigeecli /
COPY LICENSE.txt /
COPY third-party-licenses.txt /
CMD ["/apigeecli"]

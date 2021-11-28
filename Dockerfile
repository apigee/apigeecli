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

FROM golang:latest as builder
ADD ./apiclient /go/src/apigeecli/apiclient
ADD ./bundlegen /go/src/apigeecli/bundlegen
ADD ./client /go/src/apigeecli/client
ADD ./cmd /go/src/apigeecli/cmd
ADD ./clilog /go/src/apigeecli/clilog
COPY main.go /go/src/apigeecli/main.go
COPY go.mod go.sum /go/src/apigeecli/
WORKDIR /go/src/apigeecli
RUN groupadd -r -g 20000 apigee && useradd -M -u 20001 -g 0 -r -c "Default Apigee user" apigee && chown -R 20001:0 /go
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -a -ldflags='-s -w -extldflags "-static"' -o /go/bin/apigeecli /go/src/apigeecli/main.go

#without these certificates, we cannot verify the JWT token
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /etc/passwd /etc/group /etc/shadow /etc/
COPY --from=builder /go/bin/apigeecli .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt 
USER 20001
ENTRYPOINT ["/apigeecli"]
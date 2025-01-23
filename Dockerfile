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

FROM golang:1.23.2@sha256:a7f2fc9834049c1f5df787690026a53738e55fc097cd8a4a93faa3e06c67ee32 AS builder

ARG TAG
ARG COMMIT

ADD ./internal /go/src/apigeecli/internal
ADD ./cmd /go/src/apigeecli/cmd

COPY .github/workflows/licenses.tpl /go/src/apigeecli
COPY go.mod go.sum /go/src/apigeecli/

WORKDIR /go/src/apigeecli

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=true -a -gcflags='all="-l"' -ldflags='-s -w -extldflags "-static" -X main.version='${TAG}' -X main.commit='${COMMIT}' -X main.date='$(date +%FT%H:%I:%M+%Z) -o /go/bin/apigeecli /go/src/apigeecli/cmd/apigeecli/apigeecli.go
RUN GOBIN=/tmp/ go install github.com/google/go-licenses@v1.6.0
RUN /tmp/go-licenses report ./... --template /go/src/apigeecli/licenses.tpl --ignore internal > /tmp/third-party-licenses.txt 2> /dev/null || echo "Ignore warnings"

FROM ghcr.io/jqlang/jq:1.7.1@sha256:096b83865ad59b5b02841f103f83f45c51318394331bf1995e187ea3be937432 AS jq

# use debug because it includes busybox
FROM gcr.io/distroless/static-debian11:debug-nonroot@sha256:55716e80a7d4320ce9bc2dc8636fc193b418638041b817cf3306696bd0f975d1
LABEL org.opencontainers.image.url='https://github.com/apigee/apigeecli' \
      org.opencontainers.image.documentation='https://github.com/apigee/apigeecli' \
      org.opencontainers.image.source='https://github.com/apigee/apigeecli' \
      org.opencontainers.image.vendor='Google LLC' \
      org.opencontainers.image.licenses='Apache-2.0' \
      org.opencontainers.image.description='This is a tool to interact with Apigee APIs'

COPY --from=builder /go/bin/apigeecli /usr/local/bin/apigeecli
COPY --chown=nonroot:nonroot LICENSE.txt /
COPY --from=builder --chown=nonroot:nonroot /tmp/third-party-licenses.txt /
COPY --from=jq /jq /usr/local/bin/jq

ENTRYPOINT [ "apigeecli" ]

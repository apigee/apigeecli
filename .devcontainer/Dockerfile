# Copyright 2024 Google LLC
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

FROM gcr.io/projectsigstore/cosign:latest@sha256:b03690aa52bfe94054187142fba24dc54137650682810633901767d8a3e15b31 as cosign

FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm@@sha256:2e00578e7c526e76ebd51597260928d8542d27c17ffe6390d79200638c51a3e6

COPY --from=cosign /ko-app/cosign /usr/bin

#v1.1.1
RUN go install github.com/google/addlicense@dc31ac9ffcca99c9457226366135701794b128c0

# v1.6.0
RUN go install github.com/google/go-licenses@5348b744d0983d85713295ea08a20cca1654a45e

# v0.7.0
RUN go install mvdan.cc/gofumpt@86bffd62437a3c437c0b84d5d5ab244824e762fc

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/bff93a148ad87c64fafaa00f01701e0b4f26fd09/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.2

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
RUN sudo apt update
RUN sudo apt install -y goreleaser curl jq

# ref: https://stackoverflow.com/questions/28372328/how-to-install-the-google-cloud-sdk-in-a-docker-image
RUN curl https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz > /tmp/google-cloud-sdk.tar.gz

RUN mkdir -p /usr/local/gcloud \
    && tar -C /usr/local/gcloud -xvf /tmp/google-cloud-sdk.tar.gz \
    && /usr/local/gcloud/google-cloud-sdk/install.sh

ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin

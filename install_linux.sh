#!/bin/bash

# Copyright 2022 Google LLC
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

echo "Installing apigeecli to /usr/bin..."
APIGEECLI_VERSION=$(curl -s https://api.github.com/repos/apigee/apigeecli/releases/latest | jq .'name' -r)
wget https://github.com/apigee/apigeecli/releases/download/${APIGEECLI_VERSION}/apigeecli_${APIGEECLI_VERSION}_Linux_x86_64.zip
unzip apigeecli_${APIGEECLI_VERSION}_Linux_x86_64.zip
sudo mv -f apigeecli_${APIGEECLI_VERSION}_Linux_x86_64/apigeecli /usr/bin

rm -f apigeecli_${APIGEECLI_VERSION}_Linux_x86_64.zip
rm -rf apigeecli_${APIGEECLI_VERSION}_Linux_x86_64

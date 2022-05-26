#!/bin/bash -e
#
# Copyright 2017, 2019 the Velero contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

HACK_DIR=$(dirname "${BASH_SOURCE}")

echo "Updating plugin proto"

protoc pkg/plugin/proto/backupitemaction/v1/*.proto --go_out=plugins=grpc:pkg/plugin/generated/ -I pkg/plugin/proto/
protoc pkg/plugin/proto/*.proto --go_out=plugins=grpc:pkg/plugin/generated/ -I pkg/plugin/proto/
cp -rf pkg/plugin/generated/github.com/vmware-tanzu/velero/pkg/plugin/generated/* pkg/plugin/generated
rm -rf pkg/plugin/generated/github.com

echo "Updating plugin proto - done!"

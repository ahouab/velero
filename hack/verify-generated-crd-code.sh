#!/bin/bash -e
#
# Copyright 2017 the Velero contributors.
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

${HACK_DIR}/update-generated-crd-code.sh --verify-only

# ensure no changes to generated CRDs
if ! git diff --exit-code pkg/generated/crds/crds.go >/dev/null; then
  # revert changes to state before running CRD generation to stay consistent
  # with code-generator `--verify-only` option which discards generated changes
  git checkout pkg/generated/crds

  echo "run 'make update' to update generated CRDs"
  exit 1
fi

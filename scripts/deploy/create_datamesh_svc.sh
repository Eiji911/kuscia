#!/bin/bash
#
# Copyright 2023 Ant Group Co., Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

DOMAIN_ID=$1
DOMAIN_ENDPOINT=$2

usage="$(basename "$0") DOMAIN_ID DOMAIN_ENDPOINT"

if [[ ${DOMAIN_ID} == "" || ${DOMAIN_ENDPOINT} == "" ]]; then
  echo "missing argument: $usage"
  exit 1
fi

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)

DATA_MESH_SVC_TEMPLATE=$(sed "s/{{.DOMAIN_ID}}/${DOMAIN_ID}/g;
  s/{{.DATAMESH_ENDPOINT}}/${DOMAIN_ENDPOINT}/g" \
  < "${ROOT}/scripts/templates/datamesh_svc.yaml")

echo "${DATA_MESH_SVC_TEMPLATE}" | kubectl apply -f -
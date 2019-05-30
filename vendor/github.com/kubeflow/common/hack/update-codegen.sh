#!/bin/bash

# Copyright 2019 The Kubernetes Authors.
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

# This shell is used to auto generate some useful tools for k8s, such as lister,
# informer, deepcopy, defaulter and so on.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${GOPATH}/pkg/mod/k8s.io/code-generator@v0.0.0-20180621065459-6702109cc68e

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
cd ${SCRIPT_ROOT}

${CODEGEN_PKG}/generate-groups.sh "deepcopy" \
 github.com/kubeflow/common/client github.com/kubeflow/common \
 operator:v1 \
 --go-header-file hack/boilerplate/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "all" \
 github.com/kubeflow/common/client github.com/kubeflow/common \
 test_job:v1 \
 --go-header-file hack/boilerplate/boilerplate.go.txt

# Notice: The code in code-generator does not generate defaulter by default.
echo "Generating defaulters for operator/v1"
${GOPATH}/bin/defaulter-gen  --input-dirs github.com/kubeflow/common/operator/v1 -O zz_generated.defaults --go-header-file hack/boilerplate/boilerplate.go.txt "$@"
cd - > /dev/null

echo "Generating defaulters for test_job/v1"
${GOPATH}/bin/defaulter-gen  --input-dirs github.com/kubeflow/common/test_job/v1 -O zz_generated.defaults --go-header-file hack/boilerplate/boilerplate.go.txt "$@"
cd - > /dev/null

echo "Generating OpenAPI specification for operator/v1"
${GOPATH}/bin/openapi-gen --input-dirs github.com/kubeflow/common/operator/v1,k8s.io/api/core/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/version --output-package github.com/kubeflow/common/operator/v1 --go-header-file hack/boilerplate/boilerplate.go.txt "$@"
cd - > /dev/null

echo "Generating OpenAPI specification for test_job/v1"
${GOPATH}/bin/openapi-gen --input-dirs github.com/kubeflow/common/test_job/v1,k8s.io/api/core/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/version --output-package github.com/kubeflow/common/test_job/v1 --go-header-file hack/boilerplate/boilerplate.go.txt "$@"
cd - > /dev/null

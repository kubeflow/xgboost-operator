# Copyright 2019 kubeflow.org.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os

# XGBJob K8S constants
XGBJOB_GROUP = 'xgboostjob.kubeflow.org'
XGBJOB_KIND = 'XGBoostJob'
XGBJOB_PLURAL = 'xgboostjobs'
XGBJOB_VERSION = os.environ.get('XGBJOB_VERSION', 'v1alpha1')

XGBJOB_LOGLEVEL = os.environ.get('XGBJOB_LOGLEVEL', 'INFO').upper()

# How long to wait in seconds for requests to the ApiServer
APISERVER_TIMEOUT = 120

#XGBJob Labels Name
XGBJOB_CONTROLLER_LABEL = 'controller-name'
XGBJOB_GROUP_LABEL = 'group-name'
XGBJOB_NAME_LABEL = 'xgb-job-name'
XGBJOB_TYPE_LABEL = 'xgb-replica-type'
XGBJOB_INDEX_LABEL = 'xgb-replica-index'
XGBJOB_ROLE_LABEL = 'job-role'

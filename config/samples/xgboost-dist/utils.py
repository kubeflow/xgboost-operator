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


# Copyright 2018 Google Inc. All Rights Reserved.
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

import logging
from sklearn import datasets
import xgboost as xgb

logger = logging.getLogger(__name__)

def read_train_data(rank, path):
    """
    Read file based on the rank of worker. In this demo, we can use the sklearn.iris data for demonstration
    you can extend this function to read distributed data source like HDFS, HIVE etc
    :param rank:
    :param file_name:
    :return: dmatrix
    """
    iris = datasets.load_iris()
    x = iris.data
    y = iris.target
    dtrain = xgb.DMatrix(x, label=y)

    return dtrain

def dump_model(model, place):
    """
    dump the trained model into local place
    you can update this function to store the model into a remote place
    :param model:
    :param place:
    :return:
    """

    if model is None:
        raise  Exception("fail to get the xgboost train model")
    else:
        model.dump(place)
        logging.info("dump model into %s", place)

    return True
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
from sklearn.externals import joblib
import xgboost as xgb
import os
import tempfile
import oss2
import json

logger = logging.getLogger(__name__)

def extract_xgbooost_cluster_env():

    logger.info("starting to extract system env")

    master_addr = os.environ.get("MASTER_ADDR", "{}")
    master_port = int(os.environ.get("MASTER_PORT", "{}"))
    rank = int(os.environ.get("RANK", "{}"))
    world_size = int(os.environ.get("WORLD_SIZE", "{}"))

    logger.info("extract the rabit env from cluster : %s, port: %d, rank: %d, word_size: %d ",
                master_addr, master_port, rank, world_size)

    return master_addr, master_port, rank, world_size

def read_train_data(rank, num_workers, path):
    """
    Read file based on the rank of worker. In this demo, we can use the sklearn.iris data for demonstration
    you can extend this function to read distributed data source like HDFS, HIVE etc
    :param rank: the id of each worker
    :param file_name: the input file name or the place to read the data
    :return: dmatrix
    """
    iris = datasets.load_iris()
    x = iris.data
    y = iris.target

    start, end = get_range_data(len(x), rank, num_workers)
    x = x[start:end, :]
    y = y[start:end]

    logging.info("Read data from IRIS datasource with range from %d to %d", start, end)

    dtrain = xgb.DMatrix(x, label=y)

    return dtrain

def read_predict_data(rank, num_workers, path):
    """
    Read file based on the rank of worker. In this demo, we can use the sklearn.iris data for demonstration
    you can extend this function to read distributed data source like HDFS, HIVE etc
    :param rank: the id of each worker
    :param file_name: the input file name or the place to read the data
    :return: dmatrix, and real value
    """
    iris = datasets.load_iris()
    x = iris.data
    y = iris.target

    start, end = get_range_data(len(x), rank, num_workers)
    x = x[start:end, :]
    y = y[start:end]

    logging.info("Read data from IRIS datasource with range from %d to %d", start, end)

    dtrain = xgb.DMatrix(x, label=y)

    return dtrain, y

def get_range_data(num_row, rank, num_workers):
    """
    compute the data slicing range based on the input data size
    :param num_row:
    :param rank:
    :param num_workers:
    :return: begin and end range of input matrix
    """

    num_per_partition = int(num_row/num_workers)

    x_start = rank * num_per_partition
    x_end = (rank + 1 ) * num_per_partition

    if x_end > num_row:
        x_end = num_row

    return x_start, x_end

def dump_model(model, type, model_path, args):
    """
    dump the trained model into local place
    you can update this function to store the model into a remote place
    :param model:
    :param path:
    :return:
    """
    if model is None:
        raise Exception("fail to get the xgboost train model")
    else:
        if type == "local":
            joblib.dump(model, model_path)
            logging.info("Dump model into local place %s", model_path)

        elif type == "oss":
            oss_param = args.oss

            if oss_param is None:
                raise Exception("Please input the oss config to store trained model")
                return False

            oss_param['path'] = args.model_path
            dump_model_to_oss(oss_param, model)
            logging.info("Dump model into oss place %s", args.model_path)

        else:
            logging.warning("the other model storage is not supported")

    return True

def read_model(type, model_path, args):
    """
    read model from physical storage
    :param type:
    :param args:
    :return:
    """

    model = None

    if type == "local":
        model = joblib.load(model_path)
        logging.info("Read model from local place %s", model_path)

    elif type == "oss":
        oss_param = args.oss
        if oss_param is None:
            raise Exception("Please input the oss config to store trained model")
            return False

        model = read_model_from_oss(args)
        logging.info("read model from oss place %s", model_path)

    else:
        logging.warning("the other model storage is not supported")

    return model

def dump_model_to_oss(oss_parameters, booster):
    """
    dump the model to remote OSS disk
    :param oss_parameters:
    :param booster:
    :return:
    """
    """export model into oss"""
    model_fname = os.path.join(tempfile.mkdtemp(), 'model')
    text_model_fname = os.path.join(tempfile.mkdtemp(), 'model.text')
    feature_importance = os.path.join(tempfile.mkdtemp(), 'feature_importance.json')

    oss_path = oss_parameters['path']
    logger.info('---- export model ----')
    booster.save_model(model_fname)
    booster.dump_model(text_model_fname)  # format output model
    fscore_dict = booster.get_fscore()
    with open(feature_importance, 'w') as file:
        file.write(json.dumps(fscore_dict))
        logger.info('---- chief dump model successfully!')

    if os.path.exists(model_fname):
        logger.info('---- Upload Model start...')

        while oss_path[-1] == '/':
            oss_path = oss_path[:-1]

        upload_oss(oss_parameters, model_fname, oss_path)
        aux_path = oss_path + '_dir/'
        upload_oss(oss_parameters, model_fname, aux_path)
        upload_oss(oss_parameters, text_model_fname, aux_path)
        upload_oss(oss_parameters, feature_importance, aux_path)
    else:
        raise Exception("fail to generate model")
        return False

    return True

def upload_oss(kw, local_file, oss_path):
    """
    help function to upload a model to oss
    :param kw:
    :param local_file:
    :param oss_path:
    :return:
    """
    if oss_path[-1] == '/':
        oss_path = '%s%s' % (oss_path, os.path.basename(local_file))

    auth = oss2.Auth(kw['access_id'], kw['access_key'])
    bucket = kw['bucket']
    bkt = oss2.Bucket(auth=auth, endpoint=kw['endpoint'], bucket_name=bucket)

    try:
        bkt.put_object_from_file(key=oss_path, filename=local_file)
        logger.info("upload %s to %s successfully!" % (os.path.abspath(local_file), oss_path))
    except Exception():
        raise ValueError('upload %s to %s failed' % (os.path.abspath(local_file), oss_path))

def read_model_from_oss(kw):
    """
    read a model from oss
    :param kw:
    :return:
    """
    auth = oss2.Auth(kw['access_id'], kw['access_key'])
    bucket = kw['access_bucket']
    bkt = oss2.Bucket(auth=auth, endpoint=kw['endpoint'], bucket_name=bucket)
    oss_path = kw["path"]

    temp_model_fname = os.path.join(tempfile.mkdtemp(), 'local_model')
    try:
        bkt.get_object_to_file(key=oss_path, filename=temp_model_fname)
        logger.info("success to load model from oss %s", oss_path)
    except Exception():
        raise Exception("fail to load model from oss %s", oss_path)

    return xgb.Booster.load_model(temp_model_fname)

if __name__ == '__main__':
    rank = 1
    place = "/tmp/data"
    read_train_data(1, 10, place)


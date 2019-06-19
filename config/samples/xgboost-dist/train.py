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
import os
import xgboost as xgb
import traceback

from tracker import RabitTracker

logger = logging.getLogger(__name__)

def extract_xgbooost_cluster_env():
    master_addr = os.environ.get['MASTER_ADDR']
    master_port = int(os.environ.get['MASTER_PORT'])
    rank = int(os.environ.get['RANK'])
    world_size = int(os.environ.get['WORLD_SIZE'])

    return master_addr, master_port, rank, world_size

def read_train_data(rank, file_name):
    """
    Read file based on the rank of worker. we only consider the local file for demonstration.
    you can extend this function to read distributed data source like HDFS, HIVE etc
    :param rank:
    :param file_name:
    :return:
    """
    df = None

    return df

def train(args):
    """
    :param args:
    :return:
    """

    addr, port, world_size, rank = extract_xgbooost_cluster_env()

    rabit_tracker = None
    model = None

    try:
        """start to build the network"""
        if world_size > 1:
            if rank == 0:
                logger.info("start the master node")

                rabit = RabitTracker(hostIP="0.0.0.0", nslave=world_size,
                                     port=port, port_end=port + 1)
                rabit.start(world_size)
                rabit_tracker = rabit
                logger.info('########### RabitTracker Setup Finished #########')

            envs = [
                'DMLC_NUM_WORKER=%d' % world_size,
                'DMLC_TRACKER_URI=%s' % addr,
                'DMLC_TRACKER_PORT=%d' % port,
                'DMLC_TASK_ID=%d' % rank
            ]
            logger.info('##### Rabit rank setup with below envs #####')
            for i, env in enumerate(envs):
                logger.info(env)
                envs[i] = str.encode(env)

            xgb.rabit.init(envs)
            logger.info('##### Rabit rank = %d' % xgb.rabit.get_rank())

            rank = xgb.rabit.get_rank()

            df = read_train_data(rank, args.train_input)

            xgb.train()

        else:
            logging.info("only one node is started, start the train in a single node")


    except Exception as e:
        logger.error("something wrong happen: %s", traceback.format_exc())
        raise e
    finally:
        logger.info("xgboost training job finished!")
        if world_size > 1:
            xgb.rabit.finalize()
        if rabit_tracker:
            rabit_tracker.join()

    return model



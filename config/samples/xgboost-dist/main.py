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

import argparse
import logging

from train import train

def main(args):

    if args.job_type == "Predict":
        logging.info("starting the predict job")


    elif args.job_type == "Train":
        logging.info("starting the train job")
        train(args)

    elif args.job_type == "All":
        logging.info("starting the train and predict job")


    logging.info("finish distributed xgboost job")

if __name__ == '__main__':
  parser = argparse.ArgumentParser()

  parser.add_argument(
          '--job_type',
           help="Train or Predict job",
           required=True
          )
  parser.add_argument(
          '--train_input',
          help="Input training file",
          nargs='+',
          required=True
          )
  parser.add_argument(
          '--n_estimators',
          help='Number of trees in the model',
          type=int,
          default=1000
          )
  parser.add_argument(
          '--learning_rate',
          help='Learning rate for the model',
          default=0.1
          )
  parser.add_argument(
          '--model_file',
          help='Model file location for XGBoost',
          required=True
          )
  parser.add_argument(
          '--test_size',
          help='Fraction of training data to be reserved for test',
          default=0.25
          )
  parser.add_argument(
          '--early_stopping_rounds',
          help='XGBoost argument for stopping early',
          default=50
          )

  logging.basicConfig(format='%(message)s')
  logging.getLogger().setLevel(logging.INFO)
  main_args = parser.parse_args()
  main(main_args)
# Install python 2.7.
FROM python:2.7

RUN apt-get update
RUN apt-get install -y git make g++

RUN mkdir -p /opt/mlkube

# Download the rabit tracker and xgboost code.

COPY xgboost_smoke_test.py /opt/mlkube/
COPY tracker.py /opt/mlkube/

# Build XGBoost.
RUN git clone --recursive https://github.com/dmlc/xgboost && \
    cd xgboost && \
    make -j4 && \
    cd python-package; python setup.py install

ENTRYPOINT ["python", "/opt/mlkube/xgboost_smoke_test.py"]

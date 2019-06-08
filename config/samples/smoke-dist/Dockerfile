# Install python 2.7.
FROM python:2.7

RUN mkdir -p /opt/mlkube
COPY xgboost_smoke_test.py /opt/mlkube/
ENTRYPOINT ["python", "/opt/mlkube/xgboost_smoke_test.py"]

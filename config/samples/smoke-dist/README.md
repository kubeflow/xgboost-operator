### Distributed send/recv e2e test for xgboost rabit

This folder containers Dockerfile and distributed send/recv test.

**Build Image**

The default image name and tag is `kubeflow/xgboost-dist-rabit-test:1.0`.

```shell
docker build -f Dockerfile -t kubeflow/xgboost-dist-rabit-test:1.0 ./
```

**Start the XGBoost rabit tracker **

```
kubectl create -f ./xgboostjob_v1alpha1_rabit_test.yaml
```

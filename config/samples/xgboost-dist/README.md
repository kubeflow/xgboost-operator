### Distributed xgboost train and predicion

This folder containers Dockerfile and distributed xgboost model training and prediction.

**Build Image**

The default image name and tag is `kubeflow/xgboost-dist-iris-test:1.1`.

```shell
docker build -f Dockerfile -t kubeflow/xgboost-dist-iris-test:1.1 ./
```

We use the IRIS data to demonstration, 

**Start the distributed xgboost to train or predict **

```
kubectl create -f xgboostjob_v1alpha1_iris.yaml
```

### Distributed xgboost train and predicion

This folder containers Dockerfile and distributed XGBoost training and prediction.We use the [Iris Data Set](https://archive.ics.uci.edu/ml/datasets/iris) to demonstrate. 
You can extend provided data reader to read data from distributed data storage like HDFS, HBase or Hive etc.

**Build image**

The default image name and tag is `kubeflow/xgboost-dist-iris-test:1.1` respectiveily.

```shell
docker build -f Dockerfile -t kubeflow/xgboost-dist-iris-test:1.0 ./
```

Then you can push the docker image into repository
```shell
docker push kubeflow/xgboost-dist-iris-test:1.0 ./
```

**Start the distributed XGBoost train**
```
kubectl create -f xgboostjob_v1alpha1_iris.yaml 
```

**Start the distributed XGBoost predict**
```shell
kubectl create -f xgboostjob_v1alpha1_iris.yaml
```

**Look at the job status**
```
 kubectl get -o yaml XGBoostJob/xgboost-dist-iris-test
 ```

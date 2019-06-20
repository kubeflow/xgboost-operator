### Distributed send/recv e2e test for xgboost rabit

This folder containers Dockerfile and distributed send/recv test.

**Build Image**

The default image name and tag is `kubeflow/xgboost-dist-rabit-test:1.2`. 
You can build the image based on your requirement.

```shell
docker build -f Dockerfile -t kubeflow/xgboost-dist-rabit-test:1.2 ./
```

**Start and test XGBoost Rabit tracker**

```
kubectl create -f xgboostjob_v1alpha1_rabit_test.yaml
```

**Look at the job status**
```
 kubectl get -o yaml XGBoostJob/xgboost-dist-test
 ```
 
See the status section to monitor the job status. Here is sample output when the job is running.
```
 kubectl get -o yaml XGBoostJob/xgboost-dist-test
 ```
 

**Look at the pod status**



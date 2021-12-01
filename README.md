# XGBoost Operator

[![Build Status](https://travis-ci.com/kubeflow/xgboost-operator.svg?branch=master)](https://travis-ci.com/kubeflow/xgboost-operator/)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/xgboost-operator)](https://goreportcard.com/report/github.com/kubeflow/xgboost-operator)


## :warning: **kubeflow/xgboost-operator is not maintained**

This operator has been merged into [Kubeflow Training Operator](https://github.com/kubeflow/training-operator). This repository is not maintained and has been archived.

## Overview

Incubating project for [XGBoost](https://github.com/dmlc/xgboost) operator. The XGBoost operator makes it easy to run distributed XGBoost job training and batch prediction on Kubernetes cluster.

The overall design can be found [here]( https://github.com/kubeflow/community/issues/247).

This repository contains the specification and implementation of `XGBoostJob` custom resource definition.
 Using this custom resource, users can create and manage XGBoost jobs like other built-in resources in Kubernetes.
## Prerequisites
- Kubernetes >= 1.8
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)

## Install XGBoost Operator

You can deploy the operator with default settings by running the following commands using [kustomize](https://github.com/kubernetes-sigs/kustomize):

```bash
cd manifests
kubectl create namespace kubeflow
kustomize build base | kubectl apply -f -
```

Note that since Kubernetes v1.14, `kustomize` became a subcommand in `kubectl` so you can also run the following command instead:

```bash
kubectl kustomize base | kubectl apply -f -
```

## Build XGBoost Operator

XGBoost Operator is developed based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and [Kubeflow Common](https://github.com/kubeflow/common).

You can follow the [installation guide of Kubebuilder](https://book.kubebuilder.io/cronjob-tutorial/running.html) to install XGBoost operator into the Kubernetes cluster.

You can check whether the XGBoostJob custom resource has been installed via:
```
kubectl get crd
```
The output should include xgboostjobs.kubeflow.org like the following:
```
NAME                                                CREATED AT
xgboostjobs.xgboostjob.kubeflow.org                 2021-03-24T22:03:07Z
```
If it is not included you can add it as follows:
```
## setup the build enviroment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
cd $GOPATH
mkdir src/github.com/kubeflow
cd src/github.com/kubeflow

## clone the code
git clone git@github.com:kubeflow/xgboost-operator.git
cd xgboost-operator

## build and install xgboost operator
make
make install
make run
```
If the XGBoost Job operator can be installed into cluster, you can view the logs likes this

<details>
<summary>Logs</summary>

```
{"level":"info","ts":1589406873.090652,"logger":"entrypoint","msg":"setting up client for manager"}
{"level":"info","ts":1589406873.0991302,"logger":"entrypoint","msg":"setting up manager"}
{"level":"info","ts":1589406874.2192929,"logger":"entrypoint","msg":"Registering Components."}
{"level":"info","ts":1589406874.219318,"logger":"entrypoint","msg":"setting up scheme"}
{"level":"info","ts":1589406874.219448,"logger":"entrypoint","msg":"Setting up controller"}
{"level":"info","ts":1589406874.2194738,"logger":"controller","msg":"Running controller in local mode, using kubeconfig file"}
{"level":"info","ts":1589406874.224564,"logger":"controller","msg":"gang scheduling is set: ","gangscheduling":false}
{"level":"info","ts":1589406874.2247412,"logger":"kubebuilder.controller","msg":"Starting EventSource","controller":"xgboostjob-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1589406874.224958,"logger":"kubebuilder.controller","msg":"Starting EventSource","controller":"xgboostjob-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1589406874.2251048,"logger":"kubebuilder.controller","msg":"Starting EventSource","controller":"xgboostjob-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1589406874.225237,"logger":"entrypoint","msg":"setting up webhooks"}
{"level":"info","ts":1589406874.225247,"logger":"entrypoint","msg":"Starting the Cmd."}
{"level":"info","ts":1589406874.32791,"logger":"kubebuilder.controller","msg":"Starting Controller","controller":"xgboostjob-controller"}
{"level":"info","ts":1589406874.430336,"logger":"kubebuilder.controller","msg":"Starting workers","controller":"xgboostjob-controller","worker count":1}
```
</details>

## Creating a XGBoost Training/Prediction Job

You can create a XGBoost training or prediction (batch oriented) job by modifying the XGBoostJob config file.
See the distributed XGBoost Job training and prediction [example](https://github.com/kubeflow/xgboost-operator/tree/master/config/samples/xgboost-dist).
You can change the config file and related python file (i.e., train.py or predict.py)
based on your requirement.

Following the job configuration guild in the example, you can deploy a XGBoost Job to start training or prediction like:
```
## For training job
cat config/samples/xgboost-dist/xgboostjob_v1_iris_train.yaml
kubectl create -f  config/samples/xgboost-dist/xgboostjob_v1_iris_train.yaml

## For batch prediction job
cat config/samples/xgboost-dist/xgboostjob_v1_iris_predict.yaml
kubectl create -f  config/samples/xgboost-dist/xgboostjob_v1_iris_predict.yaml
```

## Monitor a distributed XGBoost Job

Once the XGBoost job is created, you should be able to watch how the related pod and service working.
Distributed XGBoost job is trained by synchronizing different worker status via tne Rabit of XGBoost.
You can also monitor the job status.

```
 kubectl get -o yaml XGBoostJob/xgboost-dist-iris-test-train
```

Here is the sample output when training job is finished.

<details>
<summary>XGBoost Job Details</summary>

```
apiVersion: xgboostjob.kubeflow.org/v1
kind: XGBoostJob
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"xgboostjob.kubeflow.org/v1","kind":"XGBoostJob","metadata":{"annotations":{},"name":"xgboost-dist-iris-test-train","namespace":"default"},"spec":{"xgbReplicaSpecs":{"Master":{"replicas":1,"restartPolicy":"Never","template":{"spec":{"containers":[{"args":["--job_type=Train","--xgboost_parameter=objective:multi:softprob,num_class:3","--n_estimators=10","--learning_rate=0.1","--model_path=/tmp/xgboost-model","--model_storage_type=local"],"image":"docker.io/merlintang/xgboost-dist-iris:1.1","imagePullPolicy":"Always","name":"xgboostjob","ports":[{"containerPort":9991,"name":"xgboostjob-port"}]}]}}},"Worker":{"replicas":2,"restartPolicy":"ExitCode","template":{"spec":{"containers":[{"args":["--job_type=Train","--xgboost_parameter=\"objective:multi:softprob,num_class:3\"","--n_estimators=10","--learning_rate=0.1"],"image":"docker.io/merlintang/xgboost-dist-iris:1.1","imagePullPolicy":"Always","name":"xgboostjob","ports":[{"containerPort":9991,"name":"xgboostjob-port"}]}]}}}}}}
  creationTimestamp: "2021-03-24T22:54:39Z"
  generation: 8
  name: xgboost-dist-iris-test-train
  namespace: default
  resourceVersion: "1060393"
  selfLink: /apis/xgboostjob.kubeflow.org/v1/namespaces/default/xgboostjobs/xgboost-dist-iris-test-train
  uid: 386c9851-7ef8-4928-9dba-2da8829bf048
spec:
  RunPolicy:
    cleanPodPolicy: None
  xgbReplicaSpecs:
    Master:
      replicas: 1
      restartPolicy: Never
      template:
        metadata:
          creationTimestamp: null
        spec:
          containers:
          - args:
            - --job_type=Train
            - --xgboost_parameter=objective:multi:softprob,num_class:3
            - --n_estimators=10
            - --learning_rate=0.1
            - --model_path=/tmp/xgboost-model
            - --model_storage_type=local
            image: docker.io/merlintang/xgboost-dist-iris:1.1
            imagePullPolicy: Always
            name: xgboostjob
            ports:
            - containerPort: 9991
              name: xgboostjob-port
            resources: {}
    Worker:
      replicas: 2
      restartPolicy: ExitCode
      template:
        metadata:
          creationTimestamp: null
        spec:
          containers:
          - args:
            - --job_type=Train
            - --xgboost_parameter="objective:multi:softprob,num_class:3"
            - --n_estimators=10
            - --learning_rate=0.1
            image: docker.io/merlintang/xgboost-dist-iris:1.1
            imagePullPolicy: Always
            name: xgboostjob
            ports:
            - containerPort: 9991
              name: xgboostjob-port
            resources: {}
status:
  completionTime: "2021-03-24T22:54:58Z"
  conditions:
  - lastTransitionTime: "2021-03-24T22:54:39Z"
    lastUpdateTime: "2021-03-24T22:54:39Z"
    message: xgboostJob xgboost-dist-iris-test-train is created.
    reason: XGBoostJobCreated
    status: "True"
    type: Created
  - lastTransitionTime: "2021-03-24T22:54:39Z"
    lastUpdateTime: "2021-03-24T22:54:39Z"
    message: XGBoostJob xgboost-dist-iris-test-train is running.
    reason: XGBoostJobRunning
    status: "False"
    type: Running
  - lastTransitionTime: "2021-03-24T22:54:58Z"
    lastUpdateTime: "2021-03-24T22:54:58Z"
    message: XGBoostJob xgboost-dist-iris-test-train is successfully completed.
    reason: XGBoostJobSucceeded
    status: "True"
    type: Succeeded
  replicaStatuses:
    Master:
      succeeded: 1
    Worker:
      succeeded: 2
 ```

</details>

## Docker Images

You can use [this Dockerfile](Dockerfile) to build the image yourself:

Alternatively, you can pull the existing image from Dockerhub [here](https://hub.docker.com/r/kubeflow/xgboost-operator/tags).

## Known Issues

XGBoost and `kubeflow/common` use pointer value in map like `map[commonv1.ReplicaType]*commonv1.ReplicaSpec`. However, `controller-gen` in [controller-tools](https://github.com/kubernetes-sigs/controller-tools) doesn't accept pointers as map values in latest version (v0.3.0), in order to generate crds and deepcopy files, we need to build custom `controller-gen`. You can follow steps below. Then `make generate` can work properly.

```shell
git clone --branch v0.2.2 git@github.com:kubernetes-sigs/controller-tools.git
git cherry-pick 71b6e91
go build -o controller-gen cmd/controller-gen/main.go
cp controller-gen /usr/local/bin
```

This can be removed once a newer `controller-gen` released and xgboost can upgrade to corresponding k8s version.

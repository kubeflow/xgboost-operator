# XGBoost Operator

[![Build Status](https://travis-ci.com/kubeflow/xgboost-operator.svg?branch=master)](https://travis-ci.com/kubeflow/xgboost-operator/)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/xgboost-operator)](https://goreportcard.com/report/github.com/kubeflow/xgboost-operator)

Incubating project for [XGBoost](https://github.com/dmlc/xgboost) operator. The XGBoost operator makes it easy to run distributed XGBoost job training and batch prediction on Kubernetes cluster.

The overall design can be found [here]( https://github.com/kubeflow/community/issues/247).

## Overview 
This repository contains the specification and implementation of `XGBoostJob` custom resource definition.
 Using this custom resource, users can create and manage XGBoost jobs like other built-in resources in Kubernetes. 
## Prerequisites
- Kubernetes >= 1.8
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)

## Installing XGBoost Operator
XGBoost Operator is developed based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and [Kubeflow Common](https://github.com/kubeflow/common). 

You can follow the [installation guide of Kubebuilder](https://book.kubebuilder.io/cronjob-tutorial/running.html) to install XGBoost operator into the Kubernetes cluster.

You can check whether the XGBoostJob custom resource has been installed via: 
```
kubectl get crd
``` 
The output should include xgboostjobs.kubeflow.org like the following:
```
NAME                                  CREATED AT
xgboostjobs.kubeflow.org   2019-06-14T06:49:45Z
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
``` 
INFO[0000] Update on create function xgboostjob-operator create object kubernetes 
INFO[0000] Update on create function xgboostjob-operator create object kube-dns 
INFO[0000] Update on create function xgboostjob-operator create object kubernetes-dashboard 
INFO[0000] Update on create function xgboostjob-operator create object storage-provisioner 
INFO[0000] Update on create function xgboostjob-operator create object kube-proxy-6l7mk 
INFO[0000] Update on create function xgboostjob-operator create object coredns-fb8b8dccf-mfsdw 
INFO[0000] Update on create function xgboostjob-operator create object coredns-fb8b8dccf-tx8nz 
INFO[0000] Update on create function xgboostjob-operator create object etcd-minikube 
INFO[0000] Update on create function xgboostjob-operator create object kube-addon-manager-minikube 
INFO[0000] Update on create function xgboostjob-operator create object kube-apiserver-minikube 
INFO[0000] Update on create function xgboostjob-operator create object kube-controller-manager-minikube 
INFO[0000] Update on create function xgboostjob-operator create object kube-scheduler-minikube 
INFO[0000] Update on create function xgboostjob-operator create object kubernetes-dashboard-79dd6bfc48-dvqzq 
{"level":"info","ts":1561676117.843403,"logger":"kubebuilder.controller","msg":"Starting Controller","controller":"xgboostjob-controller"}
{"level":"info","ts":1561676117.947829,"logger":"kubebuilder.controller","msg":"Starting workers","controller":"xgboostjob-controller","worker count":1}  
``` 
## Creating a XGBoost Training/Prediction Job

You can create a XGBoost training or prediction (batch oriented) job by modifying the XGBoostJob config file. 
See the distributed XGBoost Job training and prediction [example](https://github.com/kubeflow/xgboost-operator/tree/master/config/samples/xgboost-dist).    
You can change the config file and related python file (i.e., train.py or predict.py) 
based on your requirement. 

Following the job configuration guild in the example, you can deploy a XGBoost Job to start training or prediction like:
``` 
## For training job 
kubectl create -f xgboostjob_v1alpha1_iris_train.yaml 

## For bath prediction job 
kubectl create -f xgboostjob_v1alpha1_iris_predict.yaml
``` 

## Monitor a distributed XGBoost Job 

Once the XGBoost Job is created, you should be able to watch how th related pod and service working. 
Distributed XGBoost job is trained by synchronizing different worker status via tne Rabit of XGBoost.  
You can also monitor the job status.  

``` 
 kubectl get -o yaml XGBoostJob/xgboost-dist-iris-test-predict
``` 

Here is the sample output when training job is finished. 
```
Name:         xgboost-dist-iris-test
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  xgboostjob.kubeflow.org/v1alpha1
Kind:         XGBoostJob
Metadata:
  Creation Timestamp:  2019-06-27T01:16:09Z
  Generation:          9
  Resource Version:    385834
  Self Link:           /apis/xgboostjob.kubeflow.org/v1alpha1/namespaces/default/xgboostjobs/xgboost-dist-iris-test
  UID:                 2565e99a-9879-11e9-bbab-080027dfbfe2
Spec:
  Run Policy:
    Clean Pod Policy:  None
  Xgb Replica Specs:
    Master:
      Replicas:        1
      Restart Policy:  Never
      Template:
        Metadata:
          Creation Timestamp:  <nil>
        Spec:
          Containers:
            Args:
              --job_type=Train
              --xgboost_parameter=objective:multi:softprob,num_class:3
              --n_estimators=10
              --learning_rate=0.1
              --model_path=autoAI/xgb-opt/2
              --model_storage_type=oss
              --oss_param=unknown
            Image:              docker.io/merlintang/xgboost-dist-iris:1.1
            Image Pull Policy:  Always
            Name:               xgboostjob
            Ports:
              Container Port:  9991
              Name:            xgboostjob-port
            Resources:
    Worker:
      Replicas:        2
      Restart Policy:  ExitCode
      Template:
        Metadata:
          Creation Timestamp:  <nil>
        Spec:
          Containers:
            Args:
              --job_type=Train
              --xgboost_parameter="objective:multi:softprob,num_class:3"
              --n_estimators=10
              --learning_rate=0.1
              --model_path="/tmp/xgboost_model"
              --model_storage_type=oss
            Image:              docker.io/merlintang/xgboost-dist-iris:1.1
            Image Pull Policy:  Always
            Name:               xgboostjob
            Ports:
              Container Port:  9991
              Name:            xgboostjob-port
            Resources:
Status:
  Completion Time:  2019-06-27T01:17:04Z
  Conditions:
    Last Transition Time:  2019-06-27T01:16:09Z
    Last Update Time:      2019-06-27T01:16:09Z
    Message:               xgboostJob xgboost-dist-iris-test is created.
    Reason:                XGBoostJobCreated
    Status:                True
    Type:                  Created
    Last Transition Time:  2019-06-27T01:16:09Z
    Last Update Time:      2019-06-27T01:16:09Z
    Message:               XGBoostJob xgboost-dist-iris-test is running.
    Reason:                XGBoostJobRunning
    Status:                False
    Type:                  Running
    Last Transition Time:  2019-06-27T01:17:04Z
    Last Update Time:      2019-06-27T01:17:04Z
    Message:               XGBoostJob xgboost-dist-iris-test is successfully completed.
    Reason:                XGBoostJobSucceeded
    Status:                True
    Type:                  Succeeded
  Replica Statuses:
    Master:
      Succeeded:  1
    Worker:
      Succeeded:  2
Events:
  Type    Reason                   Age                From                 Message
  ----    ------                   ----               ----                 -------
  Normal  SuccessfulCreatePod      102s               xgboostjob-operator  Created pod: xgboost-dist-iris-test-master-0
  Normal  SuccessfulCreateService  102s               xgboostjob-operator  Created service: xgboost-dist-iris-test-master-0
  Normal  SuccessfulCreatePod      102s               xgboostjob-operator  Created pod: xgboost-dist-iris-test-worker-1
  Normal  SuccessfulCreateService  102s               xgboostjob-operator  Created service: xgboost-dist-iris-test-worker-0
  Normal  SuccessfulCreateService  102s               xgboostjob-operator  Created service: xgboost-dist-iris-test-worker-1
  Normal  SuccessfulCreatePod      64s                xgboostjob-operator  Created pod: xgboost-dist-iris-test-worker-0
  Normal  ExitedWithCode           47s (x3 over 49s)  xgboostjob-operator  Pod: default.xgboost-dist-iris-test-worker-1 exited with code 0
  Normal  ExitedWithCode           47s                xgboostjob-operator  Pod: default.xgboost-dist-iris-test-master-0 exited with code 0
  Normal  XGBoostJobSucceeded      47s                xgboostjob-operator  XGBoostJob xgboost-dist-iris-test is successfully completed.
 ```


# Docker Images

You can use the following Dockerfile to build the images yourself:

* [xgboost-operator](https://github.com/kubeflow/xgboost-operator/blob/master/Dockerfile)

For your convenience, you can pull the existing image from the following:

* [GCP-xgboost-operator-images](gcr.io/kubeflow-images-public/xgboost-operator)

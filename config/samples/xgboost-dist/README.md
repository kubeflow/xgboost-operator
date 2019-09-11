### Distributed XGBoost Job train and prediction

This folder containers related files for distributed XGBoost training and prediction. In this demo,  
[Iris Data Set](https://archive.ics.uci.edu/ml/datasets/iris) is a well known multi-class classification dataset. 
Thus, in this demo, distributed XGBoost job is able to do multi-class classification problem. Meanwhile,
User can extend provided data reader to read data from distributed data storage like HDFS, HBase or Hive etc. 


**Build image**

The default image name and tag is `kubeflow/xgboost-dist-iris-test:1.1` respectiveily.

```shell
docker build -f Dockerfile -t kubeflow/xgboost-dist-iris-test:1.0 ./
```

Then you can push the docker image into repository
```shell
docker push kubeflow/xgboost-dist-iris-test:1.0 ./
```

**Configure the job runtime via Yaml file**

There are two yaml files to setup distributed XGBoost computation runtime. 
For training job, you could configure xgboostjob_v1alpha1_iris_predict.yaml. 
Note, we use [OSS](https://www.alibabacloud.com/product/oss) to store the trained model, 
thus, you need to specify the OSS parameter in the yaml file. Therefore, remember to fill the OSS parameter in the yaml file. 
The oss parameter includes the account information such as access_id, access_key, access_bucket and endpoint. 
For Eg:
--oss_param=endpoint:http://oss-ap-south-1.aliyuncs.com,access_id:XXXXXXXXXXX,access_key:XXXXXXXXXXXXXXXXXXX,access_bucket:XXXXXX
Similarly, xgboostjob_v1alpha1_iris_predict.yaml is used to configure XGBoost job batch prediction. 

**Start the distributed XGBoost train**
```
kubectl create -f xgboostjob_v1alpha1_iris_train.yaml 
```

**Look at the train job status**
```
 kubectl get -o yaml XGBoostJob/xgboost-dist-iris-test-train
 ```
 Here is a sample output when the job is finished. The output log like this
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

**Start the distributed XGBoost job predict**
```shell
kubectl create -f xgboostjob_v1alpha1_iris_predict.yaml
```

**Look at the batch predict job status**
```
 kubectl get -o yaml XGBoostJob/xgboost-dist-iris-test-predict
 ```
 Here is a sample output when the job is finished. The output log like this
```
Name:         xgboost-dist-iris-test-predict
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  xgboostjob.kubeflow.org/v1alpha1
Kind:         XGBoostJob
Metadata:
  Creation Timestamp:  2019-06-27T06:06:53Z
  Generation:          8
  Resource Version:    394523
  Self Link:           /apis/xgboostjob.kubeflow.org/v1alpha1/namespaces/default/xgboostjobs/xgboost-dist-iris-test-predict
  UID:                 c2a04cbc-98a1-11e9-bbab-080027dfbfe2
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
              --job_type=Predict
              --model_path=autoAI/xgb-opt/3
              --model_storage_type=oss
              --oss_param=unkown
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
              --job_type=Predict
              --model_path=autoAI/xgb-opt/3
              --model_storage_type=oss
              --oss_param=unkown
            Image:              docker.io/merlintang/xgboost-dist-iris:1.1
            Image Pull Policy:  Always
            Name:               xgboostjob
            Ports:
              Container Port:  9991
              Name:            xgboostjob-port
            Resources:
Status:
  Completion Time:  2019-06-27T06:07:02Z
  Conditions:
    Last Transition Time:  2019-06-27T06:06:53Z
    Last Update Time:      2019-06-27T06:06:53Z
    Message:               xgboostJob xgboost-dist-iris-test-predict is created.
    Reason:                XGBoostJobCreated
    Status:                True
    Type:                  Created
    Last Transition Time:  2019-06-27T06:06:53Z
    Last Update Time:      2019-06-27T06:06:53Z
    Message:               XGBoostJob xgboost-dist-iris-test-predict is running.
    Reason:                XGBoostJobRunning
    Status:                False
    Type:                  Running
    Last Transition Time:  2019-06-27T06:07:02Z
    Last Update Time:      2019-06-27T06:07:02Z
    Message:               XGBoostJob xgboost-dist-iris-test-predict is successfully completed.
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
  Normal  SuccessfulCreatePod      47s                xgboostjob-operator  Created pod: xgboost-dist-iris-test-predict-worker-0
  Normal  SuccessfulCreatePod      47s                xgboostjob-operator  Created pod: xgboost-dist-iris-test-predict-worker-1
  Normal  SuccessfulCreateService  47s                xgboostjob-operator  Created service: xgboost-dist-iris-test-predict-worker-0
  Normal  SuccessfulCreateService  47s                xgboostjob-operator  Created service: xgboost-dist-iris-test-predict-worker-1
  Normal  SuccessfulCreatePod      47s                xgboostjob-operator  Created pod: xgboost-dist-iris-test-predict-master-0
  Normal  SuccessfulCreateService  47s                xgboostjob-operator  Created service: xgboost-dist-iris-test-predict-master-0
  Normal  ExitedWithCode           38s (x3 over 40s)  xgboostjob-operator  Pod: default.xgboost-dist-iris-test-predict-worker-0 exited with code 0
  Normal  ExitedWithCode           38s                xgboostjob-operator  Pod: default.xgboost-dist-iris-test-predict-master-0 exited with code 0
  Normal  XGBoostJobSucceeded      38s                xgboostjob-operator  XGBoostJob xgboost-dist-iris-test-predict is successfully completed.
```

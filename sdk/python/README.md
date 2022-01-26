# Kubeflow XGBJob SDK
Python SDK for XGB-Operator

## Requirements.

Python 2.7 and 3.5+

## Installation & Usage
### pip install

```sh
pip install kubeflow-xgbjob
```

Then import the package:
```python
from kubeflow import xgbjob 
```

### Setuptools

Install via [Setuptools](http://pypi.python.org/pypi/setuptools).

```sh
python setup.py install --user
```
(or `sudo python setup.py install` to install the package for all users)


## Getting Started

Please follow the [sample](examples/kubeflow-xgbjob-sdk.ipynb) to create, update and delete XGBJob.

## Documentation for API Endpoints

Class | Method | Description
------------ | -------------  | -------------
[XGBJobClient](docs/XGBJobClient.md) | [create](docs/XGBJobClient.md#create) | Create XGBJob|
[XGBJobClient](docs/XGBJobClient.md) | [get](docs/XGBJobClient.md#get)    | Get or watch the specified XGBJob or all XGBJob in the namespace |
[XGBJobClient](docs/XGBJobClient.md) | [patch](docs/XGBJobClient.md#patch)  | Patch the specified XGBJob|
[XGBJobClient](docs/XGBJobClient.md) | [delete](docs/XGBJobClient.md#delete) | Delete the specified XGBJob |
[XGBJobClient](docs/XGBJobClient.md) | [wait_for_job](docs/XGBJobClient.md#wait_for_job) | Wait for the specified job to finish |
[XGBJobClient](docs/XGBJobClient.md) | [wait_for_condition](docs/XGBJobClient.md#wait_for_condition) | Waits until any of the specified conditions occur |
[XGBJobClient](docs/XGBJobClient.md) | [get_job_status](docs/XGBJobClient.md#get_job_status) | Get the XGBJob status|
[XGBJobClient](docs/XGBJobClient.md) | [is_job_running](docs/XGBJobClient.md#is_job_running) | Check if the XGBJob status is Running |
[XGBJobClient](docs/XGBJobClient.md) | [is_job_succeeded](docs/XGBJobClient.md#is_job_succeeded) | Check if the XGBJob status is Succeeded |
[XGBJobClient](docs/XGBJobClient.md) | [get_pod_names](docs/XGBJobClient.md#get_pod_names) | Get pod names of XGBJob |
[XGBJobClient](docs/XGBJobClient.md) | [get_logs](docs/XGBJobClient.md#get_logs) | Get training logs of the XGBJob |

## Documentation For Models

 - [V1JobCondition](docs/V1JobCondition.md)
 - [V1JobStatus](docs/V1JobStatus.md)
 - [V1ReplicaSpec](docs/V1ReplicaSpec.md)
 - [V1ReplicaStatus](docs/V1ReplicaStatus.md)
 - [V1XGBJob](docs/V1XGBJob.md)
 - [V1XGBJobList](docs/V1XGBJobList.md)
 - [V1XGBJobSpec](docs/V1XGBJobSpec.md)


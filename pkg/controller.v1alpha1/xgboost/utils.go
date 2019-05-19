package xgboost

import (
	"testing"
	"time"
	"strings"
	"github.com/kubeflow/xgboost-operator/pkg/apis/xgboost/v1alpha1"
	jobcontroller "github.com/kubeflow/common/job_controller"
)

const (
	TestImageName      = "test-image-for-kubeflow-xgboost-operator:latest"
	TestXGBoostJobName = "test-xgboostjob"
	LabelWorker        = "worker"
	LabelMaster        = "master"

	SleepInterval = 500 * time.Millisecond
	ThreadCount   = 1

	LabelGroupName      = "group-name"
	LabelXGBoostJobName = "xgboost-job-name"

	GroupName = v1alpha1.GroupName
)

var (
	AlwaysReady = func() bool { return true }
	controllerKind = v1alpha1.SchemeGroupVersionKind

)

func GenLabels(jobName string) map[string]string {
	return map[string]string{
		LabelGroupName:      GroupName,
		LabelXGBoostJobName: strings.Replace(jobName, "/", "-", -1),
	}
}

func GetKey(job *v1alpha1.XGBoostJob, t *testing.T) string {
	key, err := jobcontroller.KeyFunc(job)
	if err != nil {
		t.Errorf("Unexpected error getting key for job %v: %v", job.Name, err)
		return ""
	}
	return key
}

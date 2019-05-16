package util

import (
	common "github.com/kubeflow/common/operator/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// JobCreatedReason is added in a job when it is created.
	JobCreatedReason = "JobCreated"
	// JobSucceededReason is added in a job when it is succeeded.
	JobSucceededReason = "JobSucceeded"
	// JobRunningReason is added in a job when it is running.
	JobRunningReason = "JobRunning"
	// JobFailedReason is added in a job when it is failed.
	JobFailedReason = "JobFailed"
	// JobRestarting is added in a job when it is restarting.
	JobRestartingReason = "JobRestarting"

	// labels for pods and servers.
	ReplicaTypeLabel  = "replica-type"
	ReplicaIndexLabel = "replica-index"
	LabelGroupName    = "group-name"
	LabelJobName      = "job-name"
	LabelJobRole      = "job-role"
)

// IsSucceeded checks if the job is succeeded
func IsSucceeded(status common.JobStatus) bool {
	return hasCondition(status, common.JobSucceeded)
}

// IsFailed checks if the job is failed
func IsFailed(status common.JobStatus) bool {
	return hasCondition(status, common.JobFailed)
}

// UpdateJobConditions adds to the jobStatus a new condition if needed, with the conditionType, reason, and message
func UpdateJobConditions(jobStatus *common.JobStatus, conditionType common.JobConditionType, reason, message string) error {
	condition := newCondition(conditionType, reason, message)
	setCondition(jobStatus, condition)
	return nil
}

func hasCondition(status common.JobStatus, condType common.JobConditionType) bool {
	for _, condition := range status.Conditions {
		if condition.Type == condType && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// newCondition creates a new job condition.
func newCondition(conditionType common.JobConditionType, reason, message string) common.JobCondition {
	return common.JobCondition{
		Type:               conditionType,
		Status:             v1.ConditionTrue,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

// getCondition returns the condition with the provided type.
func getCondition(status common.JobStatus, condType common.JobConditionType) *common.JobCondition {
	for _, condition := range status.Conditions {
		if condition.Type == condType {
			return &condition
		}
	}
	return nil
}

// setCondition updates the job to include the provided condition.
// If the condition that we are about to add already exists
// and has the same status and reason then we are not going to update.
func setCondition(status *common.JobStatus, condition common.JobCondition) {
	// Do nothing if JobStatus have failed condition
	if IsFailed(*status) {
		return
	}

	currentCond := getCondition(*status, condition.Type)

	// Do nothing if condition doesn't change
	if currentCond != nil && currentCond.Status == condition.Status && currentCond.Reason == condition.Reason {
		return
	}

	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}

	// Append the updated condition to the conditions
	newConditions := filterOutCondition(status.Conditions, condition.Type)
	status.Conditions = append(newConditions, condition)
}

// filterOutCondition returns a new slice of job conditions without conditions with the provided type.
func filterOutCondition(conditions []common.JobCondition, condType common.JobConditionType) []common.JobCondition {
	var newConditions []common.JobCondition
	for _, c := range conditions {
		if condType == common.JobRestarting && c.Type == common.JobRunning {
			continue
		}
		if condType == common.JobRunning && c.Type == common.JobRestarting {
			continue
		}

		if c.Type == condType {
			continue
		}

		// Set the running condition status to be false when current condition failed or succeeded
		if (condType == common.JobFailed || condType == common.JobSucceeded) && c.Type == common.JobRunning {
			c.Status = v1.ConditionFalse
		}

		newConditions = append(newConditions, c)
	}
	return newConditions
}

// +build !ignore_autogenerated

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// autogenerated by controller-gen object, do not modify manually

package v1alpha1

import (
	"github.com/kubeflow/common/operator/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XGBoostJob) DeepCopyInto(out *XGBoostJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XGBoostJob.
func (in *XGBoostJob) DeepCopy() *XGBoostJob {
	if in == nil {
		return nil
	}
	out := new(XGBoostJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *XGBoostJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XGBoostJobList) DeepCopyInto(out *XGBoostJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]XGBoostJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XGBoostJobList.
func (in *XGBoostJobList) DeepCopy() *XGBoostJobList {
	if in == nil {
		return nil
	}
	out := new(XGBoostJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *XGBoostJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XGBoostJobSpec) DeepCopyInto(out *XGBoostJobSpec) {
	*out = *in
	in.RunPolicy.DeepCopyInto(&out.RunPolicy)
	if in.XGBReplicaSpecs != nil {
		in, out := &in.XGBReplicaSpecs, &out.XGBReplicaSpecs
		*out = make(map[v1.ReplicaType]*v1.ReplicaSpec, len(*in))
		for key, val := range *in {
			var outVal *v1.ReplicaSpec
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(v1.ReplicaSpec)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XGBoostJobSpec.
func (in *XGBoostJobSpec) DeepCopy() *XGBoostJobSpec {
	if in == nil {
		return nil
	}
	out := new(XGBoostJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XGBoostJobStatus) DeepCopyInto(out *XGBoostJobStatus) {
	*out = *in
	in.JobStatus.DeepCopyInto(&out.JobStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XGBoostJobStatus.
func (in *XGBoostJobStatus) DeepCopy() *XGBoostJobStatus {
	if in == nil {
		return nil
	}
	out := new(XGBoostJobStatus)
	in.DeepCopyInto(out)
	return out
}

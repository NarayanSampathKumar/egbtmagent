// +build !ignore_autogenerated

/*
Copyright 2022.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BtmAgent) DeepCopyInto(out *BtmAgent) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BtmAgent.
func (in *BtmAgent) DeepCopy() *BtmAgent {
	if in == nil {
		return nil
	}
	out := new(BtmAgent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BtmAgent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BtmAgentList) DeepCopyInto(out *BtmAgentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BtmAgent, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BtmAgentList.
func (in *BtmAgentList) DeepCopy() *BtmAgentList {
	if in == nil {
		return nil
	}
	out := new(BtmAgentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BtmAgentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BtmAgentSpec) DeepCopyInto(out *BtmAgentSpec) {
	*out = *in
	if in.MonitoredNamespaces != nil {
		in, out := &in.MonitoredNamespaces, &out.MonitoredNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.UnMonitoredNamespaces != nil {
		in, out := &in.UnMonitoredNamespaces, &out.UnMonitoredNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MatchingLabels != nil {
		in, out := &in.MatchingLabels, &out.MatchingLabels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BtmAgentSpec.
func (in *BtmAgentSpec) DeepCopy() *BtmAgentSpec {
	if in == nil {
		return nil
	}
	out := new(BtmAgentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BtmAgentStatus) DeepCopyInto(out *BtmAgentStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BtmAgentStatus.
func (in *BtmAgentStatus) DeepCopy() *BtmAgentStatus {
	if in == nil {
		return nil
	}
	out := new(BtmAgentStatus)
	in.DeepCopyInto(out)
	return out
}

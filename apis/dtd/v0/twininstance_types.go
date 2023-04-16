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

package v0

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TwinInstancePhase string

const (
	TwinInstancePhasePending TwinInstancePhase = "Pending"
	TwinInstancePhaseUnknown TwinInstancePhase = "Unknown"
	TwinInstancePhaseRunning TwinInstancePhase = "Running"
	TwinInstancePhaseFailed  TwinInstancePhase = "Failed"
)

// TwinInstanceSpec defines the desired state of TwinInstance
type TwinInstanceSpec struct {
	Id             string                 `json:"id,omitempty"`
	ParentInstance string                 `json:"parentInstance,omitempty"`
	Interface      TwinInterfaceSpec      `json:"interface,omitempty"`
	Template       corev1.PodTemplateSpec `json:"template,omitempty"`
}

type TwinInstanceEvents struct {
	Filters TwinInstanceEventsFilters `json:"filters,omitempty"`
	Sink    TwinInterfaceEventsSink   `json:"sink,omitempty"`
}

// Based on CN Cloud Event Filters definitions: https://github.com/cloudevents/spec/blob/main/subscriptions/spec.md#324-filters
// TODO: build complex filtering criteria
type TwinInstanceEventsFilters struct {
	Exact  TwinInstanceEventsFiltersProperties `json:"exact,omitempty"`
	Prefix TwinInstanceEventsFiltersProperties `json:"prefix,omitempty"`
	Suffix TwinInstanceEventsFiltersProperties `json:"suffix,omitempty"`
	All    TwinInstanceEventsFiltersProperties `json:"all,omitempty"`
	Any    TwinInstanceEventsFiltersProperties `json:"any,omitempty"`
	Not    TwinInstanceEventsFiltersProperties `json:"not,omitempty"`
}

type TwinInstanceEventsFiltersProperties struct {
	Type    string `json:"type,omitempty"`
	Subject string `json:"subject,omitempty"`
}

type TwinInterfaceEventsSink struct {
	InstanceId string `json:"instanceId,omitempty"`
}

// TwinInstanceStatus defines the observed state of TwinInstance
type TwinInstanceStatus struct {
	Status TwinInstancePhase `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TwinInstance is the Schema for the twininstances API
type TwinInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TwinInstanceSpec   `json:"spec,omitempty"`
	Status TwinInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TwinInstanceList contains a list of TwinInstance
type TwinInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TwinInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TwinInstance{}, &TwinInstanceList{})
}

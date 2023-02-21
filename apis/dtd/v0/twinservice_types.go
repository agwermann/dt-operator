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

type TwinServicePhase string

const (
	TwinServicePhasePending string = "Pending"
	TwinServicePhaseUnknown string = "Unknown"
	TwinServicePhaseRunning string = "Running"
	TwinServicePhaseFailed  string = "Failed"
)

// TwinServiceSpec defines the desired state of TwinService
type TwinServiceSpec struct {
	Name       string                 `json:"name,omitempty"`
	DataSource string                 `json:"dataSource,omitempty"`
	DataTarget string                 `json:"dataTarget,omitempty"`
	Template   corev1.PodTemplateSpec `json:"template,omitempty"`
}

// TwinServiceStatus defines the observed state of TwinService
type TwinServiceStatus struct {
	Status TwinServicePhase `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TwinService is the Schema for the twinservices API
type TwinService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TwinServiceSpec   `json:"spec,omitempty"`
	Status TwinServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TwinServiceList contains a list of TwinService
type TwinServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TwinService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TwinService{}, &TwinServiceList{})
}

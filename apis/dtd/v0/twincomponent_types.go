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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TwinComponentSpec defines the desired state of TwinComponent
type TwinComponentSpec struct {
	Name       string          `json:"name,omitempty"`
	Classes    []TwinClassSpec `json:"classes,omitempty"`
	Service    TwinServiceSpec `json:"service,omitempty"`
	DataSource string          `json:"dataSource,omitempty"`
}

// TwinComponentStatus defines the observed state of TwinComponent
type TwinComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TwinComponent is the Schema for the twincomponents API
type TwinComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TwinComponentSpec   `json:"spec,omitempty"`
	Status TwinComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TwinComponentList contains a list of TwinComponent
type TwinComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TwinComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TwinComponent{}, &TwinComponentList{})
}

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

type PrimitiveTypes string
type Multiplicity string

const (
	Integer PrimitiveTypes = "integer"
	String  PrimitiveTypes = "string"
	Boolean PrimitiveTypes = "boolean"
	Double  PrimitiveTypes = "double"
)

const (
	ONE  Multiplicity = "one"
	MANY Multiplicity = "many"
)

// TwinClassSpec defines the desired state of TwinClass
type TwinClassSpec struct {
	Name          string                `json:"name"`
	Attributes    []TwinClassAttributes `json:"attributes,omitempty"`
	Relationships []TwinRelationship    `json:"relationships,omitempty"`
}

type TwinClassAttributes struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type TwinRelationship struct {
	Name         string       `json:"name,omitempty"`
	Multiplicity Multiplicity `json:"multiplicity,omitempty"`
	Reference    string       `json:"ref,omitempty"`
}

// TwinClassStatus defines the observed state of TwinClass
type TwinClassStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TwinClass is the Schema for the twinclasses API
type TwinClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TwinClassSpec   `json:"spec,omitempty"`
	Status TwinClassStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TwinClassList contains a list of TwinClass
type TwinClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TwinClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TwinClass{}, &TwinClassList{})
}

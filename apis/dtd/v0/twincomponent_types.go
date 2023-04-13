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

// TODO: review if component is a good name because DTDL has a component too

type TwinComponentPhase string

const (
	TwinComponentPhasePending TwinComponentPhase = "Pending"
	TwinComponentPhaseUnknown TwinComponentPhase = "Unknown"
	TwinComponentPhaseRunning TwinComponentPhase = "Running"
	TwinComponentPhaseFailed  TwinComponentPhase = "Failed"
)

type PrimitiveType string
type Multiplicity string

const (
	Integer PrimitiveType = "integer"
	String  PrimitiveType = "string"
	Boolean PrimitiveType = "boolean"
	Double  PrimitiveType = "double"
)

const (
	ONE  Multiplicity = "one"
	MANY Multiplicity = "many"
)

// TwinComponentSpec defines the desired state of TwinComponent
type TwinComponentSpec struct {
	Id            string              `json:"id,omitempty"`
	DisplayName   string              `json:"displayName,omitempty"`
	Description   string              `json:"description,omitempty"`
	Comment       string              `json:"comment,omitempty"`
	Properties    []TwinProperty      `json:"properties,omitempty"`
	Commands      []TwinCommand       `json:"commands,omitempty"`
	Relationships []TwinRelationship  `json:"relationships,omitempty"`
	Telemetries   []TwinTelemetry     `json:"telemetries,omitempty"`
	Extends       []TwinComponentSpec `json:"extends,omitempty"`
}

type TwinProperty struct {
	Id          string     `json:"id,omitempty"`
	Comment     string     `json:"comment,omitempty"`
	Description string     `json:"description,omitempty"`
	DisplayName string     `json:"displayName,omitempty"`
	Name        string     `json:"name,omitempty"`
	Schema      TwinSchema `json:"schema,omitempty"`
	Writeable   bool       `json:"writable,omitempty"`
}

type TwinCommand struct {
	Id          string `json:"id,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Name        string `json:"name,omitempty"`
	// Request     CommandRequest  `json:"request"`
	// Response    CommandResponse `json:"response"`
}

type TwinRelationship struct {
	Id              string         `json:"id,omitempty"`
	Comment         string         `json:"comment,omitempty"`
	Description     string         `json:"description,omitempty"`
	DisplayName     string         `json:"displayName,omitempty"`
	MaxMultiplicity int            `json:"maxMultiplicity,omitempty"`
	MinMultiplicity int            `json:"minMultiplicity,omitempty"`
	Name            string         `json:"name,omitempty"`
	Properties      []TwinProperty `json:"properties,omitempty"`
	Target          string         `json:"target,omitempty"`
	Schema          TwinSchema     `json:"schema,omitempty"`
	Writeable       bool           `json:"writeable,omitempty"`
}

type TwinTelemetry struct {
	Id          string     `json:"id,omitempty"`
	Comment     string     `json:"comment,omitempty"`
	Description string     `json:"description,omitempty"`
	DisplayName string     `json:"displayName,omitempty"`
	Name        string     `json:"name,omitempty"`
	Schema      TwinSchema `json:"schema,omitempty"`
}

type TwinSchema struct {
	PrimitiveType PrimitiveType  `json:"primitiveType,omitempty"`
	EnumType      TwinEnumSchema `json:"enumType,omitempty"`
}

type TwinEnumSchema struct {
	ValueSchema PrimitiveType          `json:"valueSchema,omitempty"`
	EnumValues  []TwinEnumSchemaValues `json:"enumValues,omitempty"`
}

type TwinEnumSchemaValues struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	EnumValue   string `json:"enumValue,omitempty"`
}

// TODO: review this definition: rename TwinComponent to TwinInterface
// TODO: TwinInstance instantiate the TwinInterface
// type Component struct {
// }

// TwinComponentStatus defines the observed state of TwinComponent
type TwinComponentStatus struct {
	Status TwinComponentPhase `json:"status,omitempty"`
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

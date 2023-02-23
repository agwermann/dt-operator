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

type MessageBrokerType string
type MessageBrokerPhase string

const (
	MESSAGE_BROKER_MQTT  MessageBrokerType = "mqtt"
	MESSAGE_BROKER_AMQP  MessageBrokerType = "rabbitmq"
	MESSAGE_BROKER_KAFKA MessageBrokerType = "kafka"
)

const (
	TwinServicePhasePending MessageBrokerPhase = "Pending"
	TwinServicePhaseUnknown MessageBrokerPhase = "Unknown"
	TwinServicePhaseRunning MessageBrokerPhase = "Running"
	TwinServicePhaseFailed  MessageBrokerPhase = "Failed"
)

type MessageBrokerSpec struct {
	Type MessageBrokerType `json:"type,omitempty"`
}

type MessageBrokerStatus struct {
	Status MessageBrokerPhase `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MessageBroker is the Schema for the messagebrokers API
type MessageBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MessageBrokerSpec   `json:"spec,omitempty"`
	Status MessageBrokerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MessageBrokerList contains a list of MessageBroker
type MessageBrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MessageBroker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MessageBroker{}, &MessageBrokerList{})
}

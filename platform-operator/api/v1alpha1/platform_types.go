package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProfileType string

const (
	Dev  ProfileType = "dev"
	Prod ProfileType = "prod"
)

type PlatformSpec struct {
	Profile ProfileType `json:"profile,omitempty"`
}

type StateType string

const (
	Ready        StateType = "Ready"
	Failed       StateType = "Failed"
	Reconciling  StateType = "Reconciling"
	Uninstalling StateType = "Uninstalling"
	Paused       StateType = "Paused"
	Upgrading    StateType = "Upgrading"
)

type PlatformStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type Platform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformSpec   `json:"spec,omitempty"`
	Status PlatformStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type PlatformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Platform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Platform{}, &PlatformList{})
}

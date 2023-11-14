package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ComponentSpec struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type PlatformSpec struct {
	Components map[string]*ComponentSpec `json:"components,omitempty"`
}

type ConditionType string

const (
	Ready ConditionType = "Ready"
)

type PlatformStatus struct {
	Conditions []metav1.Condition          `json:"conditions,omitempty"`
	Components map[string]*ComponentStatus `json:"components,omitempty"`
}

type ComponentStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Platform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformSpec   `json:"spec,omitempty"`
	Status PlatformStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type PlatformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Platform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Platform{}, &PlatformList{})
}

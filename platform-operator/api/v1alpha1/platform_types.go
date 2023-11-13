package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ArgoCDComponent struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type ComponentSpec struct {
	ArgoCD *ArgoCDComponent `json:"argoCD,omitempty"`
}

type PlatformSpec struct {
	Components ComponentSpec `json:"components,omitempty"`
}

type ConditionType string

const (
	Ready ConditionType = "Ready"
)

type PlatformStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Components ComponentStatusMap `json:"components,omitempty"`
}

type ComponentStatusMap map[string]*ComponentStatus

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

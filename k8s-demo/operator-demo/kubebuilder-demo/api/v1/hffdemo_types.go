/*
Copyright 2022 fangfenghuang.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HffDemoSpec defines the desired state of HffDemo
type HffDemoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HffDemo. Edit hffdemo_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	PVCName  string `json:"pvcName,omitempty"`
	Source   string `json:"source,omitempty"`
	DestPath string `json:"destPath,omitempty"`
}

type SourceType struct {
	URL      string `json:"url,omitempty"`
	HostPath string `json:"hostPath,omitempty"`
}

const (
	Running  = "Running"
	Pending  = "Pending"
	NotReady = "NotReady"
	Failed   = "Failed"
)

// HffDemoStatus defines the observed state of HffDemo
type HffDemoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase string `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HffDemo is the Schema for the hffdemoes API
type HffDemo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HffDemoSpec   `json:"spec,omitempty"`
	Status HffDemoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HffDemoList contains a list of HffDemo
type HffDemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HffDemo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HffDemo{}, &HffDemoList{})
}

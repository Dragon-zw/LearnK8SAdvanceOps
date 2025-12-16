/*
Copyright 2025.

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
// GuestbookSpec defines the desired state of Guestbook
type GuestbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Quantity of instances
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Size int32 `json:"size"`

	// Name of the ConfigMap for GuestbookSpec's configuration
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	ConfigMapName string `json:"configMapName"`

	// +kubebuilder:validation:Enum=Phone;Address;Name
	Type string `json:"type,omitempty"`
	// omitempty 表示可忽略的字段
}

// GuestbookStatus defines the observed state of Guestbook
type GuestbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// PodName of the active Guestbook node.
	Active string `json:"active"`

	// PodNames of the standby Guestbook nodes.
	Standby []string `json:"standby"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// CRD 为 Cluster 集群资源

// Guestbook is the Schema for the guestbooks API
type Guestbook struct {
	// 以 Kind 命名的结构体
	// 完整的 CRD 定义(代码层级的结构体)

	metav1.TypeMeta `json:",inline"` // 原生不必修改，默认需要(官方库自带)，metav1.TypeMeta 是 Group 和 Version

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"` // 原生不必修改，默认需要(官方库自带)，metav1.ObjectMeta 是 metadata 的对象

	// spec defines the desired state of Guestbook
	// +required
	Spec GuestbookSpec `json:"spec,omitempty"` // 是 Kind 中的 Spec 字段，需要开发者声明有哪一些的字段

	// status defines the observed state of Guestbook
	// +optional
	Status GuestbookStatus `json:"status,omitempty"` // 可选的状态，是资源的状态！是 Kind 中的 Status 字段(子资源)
}

// +kubebuilder:object:root=true

// GuestbookList contains a list of Guestbook
// 为 List 对象的结构体的要求定义，不必修改
type GuestbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Guestbook `json:"items"`
}

// 注册 CRD
func init() {
	SchemeBuilder.Register(&Guestbook{}, &GuestbookList{})
}

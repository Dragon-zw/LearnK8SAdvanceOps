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

// Package v1 contains API Schema definitions for the webapp v1 API group.
// +kubebuilder:object:generate=true
// +groupName=webapp.ksoperator.kubesphere
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// 定义的 Group 和 Version，内含注册的逻辑
var (
	// GroupVersion is group version used to register these objects.
    // webapp.ksoperator.kubesphere：定义 CRD 的完整资源路径，资源的 Groups
    // 即是 create api --group 指定的内容 [拼接(.)] init --domain 指定的内容
    // 在 Init 的时候需要特别注意！！！
	GroupVersion = schema.GroupVersion{Group: "webapp.ksoperator.kubesphere", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

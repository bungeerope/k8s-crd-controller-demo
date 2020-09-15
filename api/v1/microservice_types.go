/*


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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MicroServiceSpec defines the desired state of MicroService
type MicroServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Size      *int32                      `json:"size"`
	Image     string                      `json:"image"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Envs      []corev1.EnvVar             `json:"envs,omitempty"`
	Ports     []corev1.ServicePort        `json:"ports,omitempty"`
}

// MicroServiceStatus defines the observed state of MicroService
type MicroServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	appsv1.DeploymentStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// MicroService is the Schema for the microservices API
type MicroService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroServiceSpec   `json:"spec,omitempty"`
	Status MicroServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MicroServiceList contains a list of MicroService
type MicroServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroService{}, &MicroServiceList{})
}

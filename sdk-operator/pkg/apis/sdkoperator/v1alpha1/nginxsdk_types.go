package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NginxSDKSpec defines the desired state of NginxSDK
// +k8s:openapi-gen=true
type NginxSDKSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size          int32  `json:"size"`
	NodePort      int32  `json:"nodeport"`
	DefaultBody   string `json:"defaultbody"`
	ContainerPort int32  `json:"containerport"`
}

// NginxSDKStatus defines the observed state of NginxSDK
// +k8s:openapi-gen=true
type NginxSDKStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes map[string]string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NginxSDK is the Schema for the nginxsdks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nginxsdks,scope=Namespaced
type NginxSDK struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NginxSDKSpec   `json:"spec,omitempty"`
	Status NginxSDKStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NginxSDKList contains a list of NginxSDK
type NginxSDKList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NginxSDK `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NginxSDK{}, &NginxSDKList{})
}

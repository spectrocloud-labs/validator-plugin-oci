/*
Copyright 2023.

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

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OciValidatorSpec defines the desired state of OciValidator
type OciValidatorSpec struct {
	OciRegistryRules []OciRegistryRule `json:"ociRegistryRules,omitempty" yaml:"ociRegistryRules,omitempty"`
}

func (s OciValidatorSpec) ResultCount() int {
	return len(s.OciRegistryRules)
}

type OciRegistryRule struct {
	// Host is a reference to the host URL of an OCI compliant registry
	Host string `json:"host" yaml:"host"`

	// Artifacts is a slice of artifacts in the host registry that should be validated.
	// An individual artifact can take any of the following forms:
	// <repository-path>/<artifact-name>
	// <repository-path>/<artifact-name>:<tag>
	// <repository-path>/<artifact-name>@<digest>
	//
	// When no tag or digest are specified, the default tag "latest" is used.
	Artifacts []string `json:"artifacts" yaml:"artifacts"`
	Auth      Auth     `json:"auth,omitempty" yaml:"auth,omitempty"`

	// CaCert is the base64 encoded CA Certificate
	CaCert string `json:"caCert,omitempty" yaml:"caCert,omitempty"`
}

func (r OciRegistryRule) Name() string {
	return fmt.Sprintf("%s/%s", r.Host, r.Artifacts)
}

type Auth struct {
	Basic Basic `json:"basic,omitempty" yaml:"basic,omitempty"`
	Ecr   Ecr   `json:"ecr,omitempty" yaml:"ecr,omitempty"`
}

type Basic struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Ecr struct {
	AccessKey       string `json:"accessKey" yaml:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`
}

// OciValidatorStatus defines the observed state of OciValidator
type OciValidatorStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OciValidator is the Schema for the ocivalidators API
type OciValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OciValidatorSpec   `json:"spec,omitempty"`
	Status OciValidatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OciValidatorList contains a list of OciValidator
type OciValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OciValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OciValidator{}, &OciValidatorList{})
}

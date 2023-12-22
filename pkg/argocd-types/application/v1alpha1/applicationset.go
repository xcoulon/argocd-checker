package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplicationSet is a set of Application resources
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=applicationsets,shortName=appset;appsets
// +kubebuilder:subresource:status
type ApplicationSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ApplicationSetSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// Status            ApplicationSetStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ApplicationSetSpec represents a class of application set state.
type ApplicationSetSpec struct {
	GoTemplate bool                      `json:"goTemplate,omitempty" protobuf:"bytes,1,name=goTemplate"`
	Generators []ApplicationSetGenerator `json:"generators" protobuf:"bytes,2,name=generators"`
	Template   ApplicationSetTemplate    `json:"template" protobuf:"bytes,3,name=template"`
	// SyncPolicy        *ApplicationSetSyncPolicy   `json:"syncPolicy,omitempty" protobuf:"bytes,4,name=syncPolicy"`
	// Strategy          *ApplicationSetStrategy     `json:"strategy,omitempty" protobuf:"bytes,5,opt,name=strategy"`
	// PreservedFields   *ApplicationPreservedFields `json:"preservedFields,omitempty" protobuf:"bytes,6,opt,name=preservedFields"`
	// GoTemplateOptions []string                    `json:"goTemplateOptions,omitempty" protobuf:"bytes,7,opt,name=goTemplateOptions"`
	// // ApplyNestedSelectors enables selectors defined within the generators of two level-nested matrix or merge generators
	// ApplyNestedSelectors         bool                            `json:"applyNestedSelectors,omitempty" protobuf:"bytes,8,name=applyNestedSelectors"`
	// IgnoreApplicationDifferences ApplicationSetIgnoreDifferences `json:"ignoreApplicationDifferences,omitempty" protobuf:"bytes,9,name=ignoreApplicationDifferences"`
}

// ApplicationSetGenerator represents a generator at the top level of an ApplicationSet.
type ApplicationSetGenerator struct {
	// List                    *ListGenerator        `json:"list,omitempty" protobuf:"bytes,1,name=list"`
	Clusters *ClusterGenerator `json:"clusters,omitempty" protobuf:"bytes,2,name=clusters"`
	// Git                     *GitGenerator         `json:"git,omitempty" protobuf:"bytes,3,name=git"`
	// SCMProvider             *SCMProviderGenerator `json:"scmProvider,omitempty" protobuf:"bytes,4,name=scmProvider"`
	// ClusterDecisionResource *DuckTypeGenerator    `json:"clusterDecisionResource,omitempty" protobuf:"bytes,5,name=clusterDecisionResource"`
	// PullRequest             *PullRequestGenerator `json:"pullRequest,omitempty" protobuf:"bytes,6,name=pullRequest"`
	// Matrix                  *MatrixGenerator      `json:"matrix,omitempty" protobuf:"bytes,7,name=matrix"`
	// Merge                   *MergeGenerator       `json:"merge,omitempty" protobuf:"bytes,8,name=merge"`

	// Selector allows to post-filter all generator.
	Selector *metav1.LabelSelector `json:"selector,omitempty" protobuf:"bytes,9,name=selector"`

	// Plugin *PluginGenerator `json:"plugin,omitempty" protobuf:"bytes,10,name=plugin"`
}

// ApplicationSetTemplate represents argocd ApplicationSpec
type ApplicationSetTemplate struct {
	ApplicationSetTemplateMeta `json:"metadata" protobuf:"bytes,1,name=metadata"`
	Spec                       ApplicationSpec `json:"spec" protobuf:"bytes,2,name=spec"`
}

// ApplicationSetTemplateMeta represents the Argo CD application fields that may
// be used for Applications generated from the ApplicationSet (based on metav1.ObjectMeta)
type ApplicationSetTemplateMeta struct {
	Name        string            `json:"name,omitempty" protobuf:"bytes,1,name=name"`
	Namespace   string            `json:"namespace,omitempty" protobuf:"bytes,2,name=namespace"`
	Labels      map[string]string `json:"labels,omitempty" protobuf:"bytes,3,name=labels"`
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,4,name=annotations"`
	Finalizers  []string          `json:"finalizers,omitempty" protobuf:"bytes,5,name=finalizers"`
}

// ClusterGenerator defines a generator to match against clusters registered with ArgoCD.
type ClusterGenerator struct {
	// Selector defines a label selector to match against all clusters registered with ArgoCD.
	// Clusters today are stored as Kubernetes Secrets, thus the Secret labels will be used
	// for matching the selector.
	Selector metav1.LabelSelector   `json:"selector,omitempty" protobuf:"bytes,1,name=selector"`
	Template ApplicationSetTemplate `json:"template,omitempty" protobuf:"bytes,2,name=template"`

	// Values contains key/value pairs which are passed directly as parameters to the template
	Values map[string]string `json:"values,omitempty" protobuf:"bytes,3,name=values"`
}

// // ApplicationSetList contains a list of ApplicationSet
// // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// // +kubebuilder:object:root=true
// type ApplicationSetList struct {
// 	metav1.TypeMeta `json:",inline"`
// 	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
// 	Items           []ApplicationSet `json:"items" protobuf:"bytes,2,rep,name=items"`
// }

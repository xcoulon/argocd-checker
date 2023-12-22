package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Application is a definition of Application resource.
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// see https://github.com/argoproj/argo-cd/blob/master/pkg/apis/application/v1alpha1/types.go
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              ApplicationSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// Status            ApplicationStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	// Operation         *Operation        `json:"operation,omitempty" protobuf:"bytes,4,opt,name=operation"`
}

// ApplicationSpec represents desired application state. Contains link to repository with application definition and additional parameters link definition revision.
type ApplicationSpec struct {
	// Source is a reference to the location of the application's manifests or chart
	Source *ApplicationSource `json:"source,omitempty" protobuf:"bytes,1,opt,name=source"`
	// Destination is a reference to the target Kubernetes server and namespace
	Destination ApplicationDestination `json:"destination" protobuf:"bytes,2,name=destination"`
	// Project is a reference to the project this application belongs to.
	// The empty string means that application belongs to the 'default' project.
	Project string `json:"project" protobuf:"bytes,3,name=project"`
	// SyncPolicy controls when and how a sync will be performed
	SyncPolicy *SyncPolicy `json:"syncPolicy,omitempty" protobuf:"bytes,4,name=syncPolicy"`
	// IgnoreDifferences is a list of resources and their fields which should be ignored during comparison
	IgnoreDifferences IgnoreDifferences `json:"ignoreDifferences,omitempty" protobuf:"bytes,5,name=ignoreDifferences"`
	// Info contains a list of information (URLs, email addresses, and plain text) that relates to the application
	Info []Info `json:"info,omitempty" protobuf:"bytes,6,name=info"`
	// RevisionHistoryLimit limits the number of items kept in the application's revision history, which is used for informational purposes as well as for rollbacks to previous versions.
	// This should only be changed in exceptional circumstances.
	// Setting to zero will store no history. This will reduce storage used.
	// Increasing will increase the space used to store the history, so we do not recommend increasing it.
	// Default is 10.
	RevisionHistoryLimit *int64 `json:"revisionHistoryLimit,omitempty" protobuf:"bytes,7,name=revisionHistoryLimit"`

	// Sources is a reference to the location of the application's manifests or chart
	Sources ApplicationSources `json:"sources,omitempty" protobuf:"bytes,8,opt,name=sources"`
}

// ApplicationSource contains all required information about the source of an application
type ApplicationSource struct {
	// RepoURL is the URL to the repository (Git or Helm) that contains the application manifests
	RepoURL string `json:"repoURL" protobuf:"bytes,1,opt,name=repoURL"`
	// Path is a directory path within the Git repository, and is only valid for applications sourced from Git.
	Path string `json:"path,omitempty" protobuf:"bytes,2,opt,name=path"`
	// TargetRevision defines the revision of the source to sync the application to.
	// In case of Git, this can be commit, tag, or branch. If omitted, will equal to HEAD.
	// In case of Helm, this is a semver tag for the Chart's version.
	TargetRevision string `json:"targetRevision,omitempty" protobuf:"bytes,4,opt,name=targetRevision"`
	// // Helm holds helm specific options
	// Helm *ApplicationSourceHelm `json:"helm,omitempty" protobuf:"bytes,7,opt,name=helm"`
	// // Kustomize holds kustomize specific options
	// Kustomize *ApplicationSourceKustomize `json:"kustomize,omitempty" protobuf:"bytes,8,opt,name=kustomize"`
	// // Directory holds path/directory specific options
	// Directory *ApplicationSourceDirectory `json:"directory,omitempty" protobuf:"bytes,10,opt,name=directory"`
	// // Plugin holds config management plugin specific options
	// Plugin *ApplicationSourcePlugin `json:"plugin,omitempty" protobuf:"bytes,11,opt,name=plugin"`
	// // Chart is a Helm chart name, and must be specified for applications sourced from a Helm repo.
	// Chart string `json:"chart,omitempty" protobuf:"bytes,12,opt,name=chart"`
	// // Ref is reference to another source within sources field. This field will not be used if used with a `source` tag.
	// Ref string `json:"ref,omitempty" protobuf:"bytes,13,opt,name=ref"`
}

// ApplicationSources contains list of required information about the sources of an application
type ApplicationSources []ApplicationSource

// ApplicationDestination holds information about the application's destination
type ApplicationDestination struct {
	// Server specifies the URL of the target cluster's Kubernetes control plane API. This must be set if Name is not set.
	Server string `json:"server,omitempty" protobuf:"bytes,1,opt,name=server"`
	// Namespace specifies the target namespace for the application's resources.
	// The namespace will only be set for namespace-scoped resources that have not set a value for .metadata.namespace
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	// Name is an alternate way of specifying the target cluster by its symbolic name. This must be set if Server is not set.
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
}

// SyncPolicy controls when a sync will be performed in response to updates in git
type SyncPolicy struct {
	// Automated will keep an application synced to the target revision
	Automated *SyncPolicyAutomated `json:"automated,omitempty" protobuf:"bytes,1,opt,name=automated"`
	// Options allow you to specify whole app sync-options
	SyncOptions SyncOptions `json:"syncOptions,omitempty" protobuf:"bytes,2,opt,name=syncOptions"`
	// Retry controls failed sync retry behavior
	Retry *RetryStrategy `json:"retry,omitempty" protobuf:"bytes,3,opt,name=retry"`
	// ManagedNamespaceMetadata controls metadata in the given namespace (if CreateNamespace=true)
	ManagedNamespaceMetadata *ManagedNamespaceMetadata `json:"managedNamespaceMetadata,omitempty" protobuf:"bytes,4,opt,name=managedNamespaceMetadata"`
}

// SyncPolicyAutomated controls the behavior of an automated sync
type SyncPolicyAutomated struct {
	// Prune specifies whether to delete resources from the cluster that are not found in the sources anymore as part of automated sync (default: false)
	Prune bool `json:"prune,omitempty" protobuf:"bytes,1,opt,name=prune"`
	// SelfHeal specifies whether to revert resources back to their desired state upon modification in the cluster (default: false)
	SelfHeal bool `json:"selfHeal,omitempty" protobuf:"bytes,2,opt,name=selfHeal"`
	// AllowEmpty allows apps have zero live resources (default: false)
	AllowEmpty bool `json:"allowEmpty,omitempty" protobuf:"bytes,3,opt,name=allowEmpty"`
}

type SyncOptions []string

// RetryStrategy contains information about the strategy to apply when a sync failed
type RetryStrategy struct {
	// Limit is the maximum number of attempts for retrying a failed sync. If set to 0, no retries will be performed.
	Limit int64 `json:"limit,omitempty" protobuf:"bytes,1,opt,name=limit"`
	// Backoff controls how to backoff on subsequent retries of failed syncs
	Backoff *Backoff `json:"backoff,omitempty" protobuf:"bytes,2,opt,name=backoff,casttype=Backoff"`
}

// Backoff is the backoff strategy to use on subsequent retries for failing syncs
type Backoff struct {
	// Duration is the amount to back off. Default unit is seconds, but could also be a duration (e.g. "2m", "1h")
	Duration string `json:"duration,omitempty" protobuf:"bytes,1,opt,name=duration"`
	// Factor is a factor to multiply the base duration after each failed retry
	Factor *int64 `json:"factor,omitempty" protobuf:"bytes,2,name=factor"`
	// MaxDuration is the maximum amount of time allowed for the backoff strategy
	MaxDuration string `json:"maxDuration,omitempty" protobuf:"bytes,3,opt,name=maxDuration"`
}

type ManagedNamespaceMetadata struct {
	Labels      map[string]string `json:"labels,omitempty" protobuf:"bytes,1,opt,name=labels"`
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,2,opt,name=annotations"`
}

type IgnoreDifferences []ResourceIgnoreDifferences

// ResourceIgnoreDifferences contains resource filter and list of json paths which should be ignored during comparison with live state.
type ResourceIgnoreDifferences struct {
	Group             string   `json:"group,omitempty" protobuf:"bytes,1,opt,name=group"`
	Kind              string   `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	Name              string   `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	Namespace         string   `json:"namespace,omitempty" protobuf:"bytes,4,opt,name=namespace"`
	JSONPointers      []string `json:"jsonPointers,omitempty" protobuf:"bytes,5,opt,name=jsonPointers"`
	JQPathExpressions []string `json:"jqPathExpressions,omitempty" protobuf:"bytes,6,opt,name=jqPathExpressions"`
	// ManagedFieldsManagers is a list of trusted managers. Fields mutated by those managers will take precedence over the
	// desired state defined in the SCM and won't be displayed in diffs
	ManagedFieldsManagers []string `json:"managedFieldsManagers,omitempty" protobuf:"bytes,7,opt,name=managedFieldsManagers"`
}

type Info struct {
	Name  string `json:"name" protobuf:"bytes,1,name=name"`
	Value string `json:"value" protobuf:"bytes,2,name=value"`
}

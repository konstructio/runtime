/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package detokenization

type GitOpsDirectoryValues struct {
	AlertsEmail               string
	AtlantisAllowList         string
	CloudProvider             string
	CloudRegion               string
	ClusterId                 string
	ClusterName               string
	ClusterType               string
	DomainName                string
	Kubeconfig                string
	KubeconfigPath            string
	KubefirstArtifactsBucket  string
	KubefirstStateStoreBucket string
	KubefirstTeam             string
	KubefirstVersion          string
	StateStoreBucketHostname  string

	ArgoCDIngressURL               string
	ArgoCDIngressNoHTTPSURL        string
	ArgoWorkflowsIngressURL        string
	ArgoWorkflowsIngressNoHTTPSURL string
	ArgoWorkflowsDir               string
	AtlantisIngressURL             string
	AtlantisIngressNoHTTPSURL      string
	ChartMuseumIngressURL          string
	VaultDataBucketName            string
	VaultIngressURL                string
	VaultIngressNoHTTPSURL         string
	VouchIngressURL                string

	GitDescription       string
	GitNamespace         string
	GitProvider          string
	GitRunner            string
	GitRunnerDescription string
	GitRunnerNS          string
	GitURL               string

	GitHubHost  string
	GitHubOwner string
	GitHubUser  string

	GitlabHost         string
	GitlabOwner        string
	GitlabOwnerGroupID int
	GitlabUser         string

	GitOpsRepoAtlantisWebhookURL string
	GitOpsRepoGitURL             string
	GitOpsRepoNoHTTPSURL         string

	ContainerRegistryURL string
	UseTelemetry         string

	AtlantisWebhookURL   string
	AwsIamArnAccountRoot string
	AwsKmsKeyId          string
	AwsNodeCapacityType  string
	AwsAccountID         string

	GCPAuth    string
	GCPProject string

	// k3d compatibility
	MetaphorDevelopmentIngressURL string
	MetaphorStagingIngressURL     string
	MetaphorProductionIngressURL  string
}

type MetaphorTokenValues struct {
	CheckoutCWFTTemplate          string
	CloudRegion                   string
	ClusterName                   string
	CommitCWFTTemplate            string
	ContainerRegistryURL          string
	DomainName                    string
	MetaphorDevelopmentIngressURL string
	MetaphorProductionIngressURL  string
	MetaphorStagingIngressURL     string
}

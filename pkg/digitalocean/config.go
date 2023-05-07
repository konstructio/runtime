/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package digitalocean

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kubefirst/runtime/pkg"
	"github.com/rs/zerolog/log"
)

const (
	CloudProvider          = "digitalocean"
	GithubHost             = "github.com"
	GitlabHost             = "gitlab.com"
	KubectlClientVersion   = "v1.25.7"
	LocalhostOS            = runtime.GOOS
	LocalhostArch          = runtime.GOARCH
	TerraformClientVersion = "1.3.8"
	ArgocdHelmChartVersion = "4.10.5"

	ArgocdPortForwardURL = pkg.ArgocdPortForwardURL
	VaultPortForwardURL  = pkg.VaultPortForwardURL
)

type DigitaloceanConfig struct {
	DigitaloceanToken string
	GithubToken       string
	GitlabToken       string

	ArgoWorkflowsDir                string
	DestinationGitopsRepoHttpsURL   string
	DestinationGitopsRepoGitURL     string
	DestinationMetaphorRepoHttpsURL string
	DestinationMetaphorRepoGitURL   string
	GitopsDir                       string
	GitProvider                     string
	K1Dir                           string
	Kubeconfig                      string
	KubectlClient                   string
	KubefirstBotSSHPrivateKey       string
	KubefirstConfig                 string
	LogsDir                         string
	MetaphorDir                     string
	RegistryAppName                 string
	RegistryYaml                    string
	SSLBackupDir                    string
	TerraformClient                 string
	ToolsDir                        string
}

// GetConfig - load default values from kubefirst installer
func GetConfig(clusterName string, domainName string, gitProvider string, gitOwner string) *DigitaloceanConfig {
	config := DigitaloceanConfig{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Msgf("something went wrong getting home path: %s", err)
	}

	// cGitHost describes which git host to use depending on gitProvider
	var cGitHost string
	switch gitProvider {
	case "github":
		cGitHost = GithubHost
	case "gitlab":
		cGitHost = GitlabHost
	}

	config.DestinationGitopsRepoHttpsURL = fmt.Sprintf("https://%s/%s/gitops.git", cGitHost, gitOwner)
	config.DestinationGitopsRepoGitURL = fmt.Sprintf("git@%s:%s/gitops.git", cGitHost, gitOwner)
	config.DestinationMetaphorRepoHttpsURL = fmt.Sprintf("https://%s/%s/metaphor.git", cGitHost, gitOwner)
	config.DestinationMetaphorRepoGitURL = fmt.Sprintf("git@%s:%s/metaphor.git", cGitHost, gitOwner)

	config.ArgoWorkflowsDir = fmt.Sprintf("%s/.k1/%s/argo-workflows", homeDir, clusterName)
	config.GitopsDir = fmt.Sprintf("%s/.k1/%s/gitops", homeDir, clusterName)
	config.GitProvider = gitProvider
	config.Kubeconfig = fmt.Sprintf("%s/.k1/%s/kubeconfig", homeDir, clusterName)
	config.K1Dir = fmt.Sprintf("%s/.k1/%s", homeDir, clusterName)
	config.KubectlClient = fmt.Sprintf("%s/.k1/%s/tools/kubectl", homeDir, clusterName)
	config.KubefirstConfig = fmt.Sprintf("%s/.k1/%s/%s", homeDir, clusterName, ".kubefirst")
	config.LogsDir = fmt.Sprintf("%s/.k1/%s/logs", homeDir, clusterName)
	config.MetaphorDir = fmt.Sprintf("%s/.k1/%s/metaphor", homeDir, clusterName)
	config.RegistryAppName = "registry"
	config.RegistryYaml = fmt.Sprintf("%s/.k1/%s/gitops/registry/%s/registry.yaml", homeDir, clusterName, clusterName)
	config.SSLBackupDir = fmt.Sprintf("%s/.k1/%s/ssl/%s", homeDir, clusterName, domainName)
	config.TerraformClient = fmt.Sprintf("%s/.k1/%s/tools/terraform", homeDir, clusterName)
	config.ToolsDir = fmt.Sprintf("%s/.k1/%s/tools", homeDir, clusterName)

	return &config
}

type GitOpsDirectoryValues struct {
	AlertsEmail               string
	AtlantisAllowList         string
	CloudProvider             string
	CloudRegion               string
	ClusterId                 string
	ClusterName               string
	ClusterType               string
	DomainName                string
	KubeconfigPath            string
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

	UseTelemetry string
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

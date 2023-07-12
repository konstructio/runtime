/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package providerConfigs

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

type ProviderConfig struct {
	CivoToken         string
	DigitaloceanToken string
	GCPAuth           string
	GCPProject        string
	VultrToken        string

	GithubToken string
	GitlabToken string

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
func GetConfig(clusterName string, domainName string, gitProvider string, gitOwner string) *ProviderConfig {
	config := ProviderConfig{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msgf("something went wrong getting home path: %s", err)
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

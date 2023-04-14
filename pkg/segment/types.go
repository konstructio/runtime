/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package segment

import (
	"github.com/segmentio/analytics-go"
)

// SegmentIO constants
// SegmentIOWriteKey The write key is the unique identifier for a source that tells Segment which source data comes
// from, to which workspace the data belongs, and which destinations should receive the data.
const (
	SegmentIOWriteKey = "0gAYkX5RV3vt7s4pqCOOsDb6WHPLT30M"

	MetricInitStarted                 = "kubefirst.init.started"
	MetricInitCompleted               = "kubefirst.init.completed"
	MetricMgmtClusterInstallStarted   = "kubefirst.mgmt_cluster_install.started"
	MetricMgmtClusterInstallCompleted = "kubefirst.mgmt_cluster_install.completed"

	MetricCloudCredentialsCheckStarted   = "kubefirst.cloud_credentials_check.started"
	MetricCloudCredentialsCheckCompleted = "kubefirst.cloud_credentials_check.completed"
	MetricKbotSetupStarted               = "kubefirst.kbot_setup.started"
	MetricKbotSetupCompleted             = "kubefirst.kbot_setup.completed"
	MetricGitTerraformApplyStarted       = "kubefirst.git_terraform_apply.started"
	MetricGitTerraformApplyCompleted     = "kubefirst.git_terraform_apply.completed"
	MetricGitopsRepoPushStarted          = "kubefirst.gitops_repo_push.started"
	MetricGitopsRepoPushCompleted        = "kubefirst.gitops_repo_push.completed"
	MetricCloudTerraformApplyStarted     = "kubefirst.cloud_terraform_apply.started"
	MetricCloudTerraformApplyCompleted   = "kubefirst.cloud_terraform_apply.completed"
	MetricArgoCDInstallStarted           = "kubefirst.argocd_install.started"
	MetricArgoCDInstallCompleted         = "kubefirst.argocd_install.completed"
	MetricCreateRegistryStarted          = "kubefirst.create_registry.started"
	MetricCreateRegistryCompleted        = "kubefirst.create_registry.completed"
	MetricVaultInitializationStarted     = "kubefirst.vault_initialization.started"
	MetricVaultInitializationCompleted   = "kubefirst.vault_initialization.completed"
	MetricVaultTerraformApplyStarted     = "kubefirst.vault_terraform_apply.started"
	MetricVaultTerraformApplyCompleted   = "kubefirst.vault_terraform_apply.completed"
	MetricUsersTerraformApplyStarted     = "kubefirst.users_terraform_apply.started"
	MetricUsersTerraformApplyCompleted   = "kubefirst.users_terraform_apply.completed"
)

type SegmentClient struct {
	Client            analytics.Client
	CliVersion        string
	CloudProvider     string
	ClusterID         string
	ClusterType       string
	DomainName        string
	GitProvider       string
	KubefirstTeam     string
	KubefirstTeamInfo string
}

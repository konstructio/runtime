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
	SegmentIOWriteKey                 = "0gAYkX5RV3vt7s4pqCOOsDb6WHPLT30M"
	MetricInitStarted                 = "kubefirst.init.started"
	MetricInitCompleted               = "kubefirst.init.completed"
	MetricMgmtClusterInstallStarted   = "kubefirst.mgmt_cluster_install.started"
	MetricMgmtClusterInstallCompleted = "kubefirst.mgmt_cluster_install.completed"
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
	MachineID         string
}

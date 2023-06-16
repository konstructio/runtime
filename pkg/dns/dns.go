/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package dns

import (
	"fmt"
	"strings"

	"github.com/kubefirst/runtime/pkg"
	awsinternal "github.com/kubefirst/runtime/pkg/aws"
	"github.com/lixiangzhong/dnsutil"
	"github.com/rs/zerolog/log"
)

const (
	// Google
	dnsLookupHost string = "8.8.8.8"
)

var (
	CivoNameServers         []string = []string{"ns0.civo.com", "ns1.civo.com"}
	DigitalOceanNameServers []string = []string{"ns1.digitalocean.com", "ns2.digitalocean.com", "ns3.digitalocean.com"}
	VultrNameservers        []string = []string{"ns1.vultr.com", "ns2.vultr.com"}
)

// VerifyProviderDNS
func VerifyProviderDNS(cloudProvider string, cloudRegion string, domainName string) error {
	var nameServers []string

	switch cloudProvider {
	case "aws":
		awsClient := &awsinternal.AWSConfiguration{
			Config: awsinternal.NewAwsV2(cloudRegion),
		}
		hostedZoneID, err := awsClient.GetHostedZoneID(domainName)
		if err != nil {
			return err
		}
		hostedZone, err := awsClient.GetHostedZone(hostedZoneID)
		if err != nil {
			return err
		}
		nameServers = hostedZone.DelegationSet.NameServers
	case "civo":
		nameServers = CivoNameServers
	case "digitalocean":
		nameServers = DigitalOceanNameServers
	case "vultr":
		nameServers = VultrNameservers
	default:
		return fmt.Errorf("unsupported cloud provider for dns verification: %s", cloudProvider)
	}

	foundNSRecords, err := GetDomainNSRecords(domainName)
	if err != nil {
		return err
	}

	for _, reqrec := range nameServers {
		if pkg.FindStringInSlice(foundNSRecords, reqrec) {
			log.Info().Msgf("found NS record %s for domain %s", reqrec, domainName)
		} else {
			return fmt.Errorf("missing record %s for domain %s - please add the NS record", reqrec, domainName)
		}
	}

	return nil
}

// GetDomainNSRecords
func GetDomainNSRecords(domainName string) ([]string, error) {
	var dig dnsutil.Dig
	dig.SetDNS(dnsLookupHost)

	records, err := dig.NS(domainName)
	if err != nil {
		return []string{}, fmt.Errorf("error checking NS record for domain %s: %s", domainName, err)
	}

	var foundNSRecords []string
	for _, rec := range records {
		foundNSRecords = append(foundNSRecords, strings.TrimSuffix(rec.Ns, "."))
	}

	return foundNSRecords, nil
}

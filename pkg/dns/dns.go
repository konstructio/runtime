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
	"github.com/lixiangzhong/dnsutil"
	"github.com/rs/zerolog/log"
)

var (
	CivoNameServers         []string = []string{"ns0.civo.com", "ns1.civo.com"}
	DigitalOceanNameServers []string = []string{"ns1.digitalocean.com", "ns2.digitalocean.com", "ns3.digitalocean.com"}
	VultrNameservers        []string = []string{"ns1.vultr.com", "ns2.vultr.com"}
)

// VerifyProviderDNS
func VerifyProviderDNS(cloudProvider string, domainName string) error {
	var dig dnsutil.Dig
	var nameServers []string
	dig.SetDNS("8.8.8.8")

	switch cloudProvider {
	case "civo":
		nameServers = CivoNameServers
	case "digitalocean":
		nameServers = DigitalOceanNameServers
	case "vultr":
		nameServers = VultrNameservers
	default:
		return fmt.Errorf("unsupported cloud provider for dns verification: %s", cloudProvider)
	}

	records, err := dig.NS(domainName)
	if err != nil {
		return fmt.Errorf("error checking NS record for domain %s: %s", domainName, err)
	}

	var foundNSRecords []string
	for _, rec := range records {
		foundNSRecords = append(foundNSRecords, strings.TrimSuffix(rec.Ns, "."))
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

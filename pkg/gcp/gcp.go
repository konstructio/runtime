/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package gcp

import (
	"fmt"
	"net/http"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

// GetRegions lists all available regions
func (conf *GCPConfiguration) GetRegions() ([]string, error) {
	var regionList []string

	client, err := compute.NewRegionsRESTClient(conf.Context)
	if err != nil {
		return []string{}, fmt.Errorf("could not create google compute client: %s", err)
	}
	defer client.Close()

	req := &computepb.ListRegionsRequest{
		Project: conf.Project,
	}

	it := client.List(conf.Context, req)
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []string{}, err
		}
		regionList = append(regionList, *pair.Name)
	}

	return regionList, nil
}

// GetDomainApexContent determines whether or not a target domain features
// a host responding at zone apex
func GetDomainApexContent(domainName string) bool {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	exists := false
	for _, proto := range []string{"http", "https"} {
		fqdn := fmt.Sprintf("%s://%s", proto, domainName)
		_, err := client.Get(fqdn)
		if err != nil {
			log.Warn().Msgf("domain %s has no apex content", fqdn)
		} else {
			log.Info().Msgf("domain %s has apex content", fqdn)
			exists = true
		}
	}

	return exists
}

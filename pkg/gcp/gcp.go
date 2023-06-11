/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package gcp

import (
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
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

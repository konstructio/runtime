/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
)

// GetRegions lists all available regions
func (c *DigitaloceanConfiguration) GetRegions() ([]string, error) {
	var regionList []string

	regions, _, err := c.Client.Regions.List(c.Context, &godo.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	for _, region := range regions {
		regionList = append(regionList, region.Slug)
	}

	return regionList, nil
}

func (c *DigitaloceanConfiguration) ListInstances() ([]string, error) {

	instances, _, err := c.Client.Apps.ListInstanceSizes(context.Background())
	
	if err !=  nil {
		return nil, err
	}
	var instanceNames []string
	for _, instance := range instances {
		instanceNames = append(instanceNames, instance.Name)
	}

	return instanceNames, nil
}

/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package digitalocean

import (
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

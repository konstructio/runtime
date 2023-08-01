/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package cloudflare

// GetDNSDomains lists all available DNS domains
func (c *CloudflareConfiguration) GetDNSDomains() ([]string, error) {
	var domainList []string

	zones, err := c.Client.ListZones(c.Context)
	if err != nil {
		return []string{}, err
	}

	for _, domain := range zones {
		domainList = append(domainList, domain.Name)
	}

	return domainList, nil
}

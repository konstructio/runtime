/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package segment

import (
	"github.com/segmentio/analytics-go"
)

// SetupClient associates the Segment client with an instance of the local client
func (c *SegmentClient) SetupClient() {
	c.Client = analytics.New(SegmentIOWriteKey)
}

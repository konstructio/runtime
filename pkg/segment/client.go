package segment

import (
	"github.com/kubefirst/runtime/pkg/helpers"
	"github.com/segmentio/analytics-go"
)

var Client SegmentClient = SegmentClient{
	Client: newSegmentClient(),
}

func newSegmentClient() analytics.Client {

	client := analytics.New(helpers.SegmentIOWriteKey)

	return client
}

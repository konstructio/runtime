/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package gcp

import "context"

// GCPConfiguration stores session data to organize all GCP functions into a single struct
type GCPConfiguration struct {
	Context context.Context
	Project string
	Region  string
}

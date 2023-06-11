/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package gcp

import (
	"fmt"

	storage "cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

// CreateBucket creates a GCS bucket
func (conf *GCPConfiguration) CreateBucket(bucketName string) (*storage.BucketAttrs, error) {
	client, err := storage.NewClient(conf.Context)
	if err != nil {
		return nil, fmt.Errorf("could not create google storage client: %s", err)
	}
	defer client.Close()

	// Create bucket
	log.Info().Msgf("creating gcs bucket %s", bucketName)

	err = client.Bucket(bucketName).Create(conf.Context, conf.Project, &storage.BucketAttrs{})
	if err != nil {
		return nil, fmt.Errorf("error creating gcs bucket %s: %s", bucketName, err)
	}

	it := client.Buckets(conf.Context, conf.Project)
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			return nil, fmt.Errorf("error fetching created bucket: %s", err)
		}
		if err != nil {
			return nil, err
		}
		if pair.Name == bucketName {
			return pair, nil
		}
	}
}

// DeleteBucket deletes a GCS bucket
func (conf *GCPConfiguration) DeleteBucket(bucketName string) error {
	client, err := storage.NewClient(conf.Context)
	if err != nil {
		return fmt.Errorf("could not create google storage client: %s", err)
	}
	defer client.Close()

	// Create bucket
	log.Info().Msgf("deleting gcs bucket %s", bucketName)

	bucket := client.Bucket(bucketName)
	err = bucket.Delete(conf.Context)
	if err != nil {
		return fmt.Errorf("error deleting gcs bucket %s: %s", bucketName, err)
	}

	return nil
}

// ListBuckets lists all GCS buckets for a project
func (conf *GCPConfiguration) ListBuckets() ([]*storage.BucketAttrs, error) {
	client, err := storage.NewClient(conf.Context)
	if err != nil {
		return nil, fmt.Errorf("could not create google storage client: %s", err)
	}
	defer client.Close()

	var buckets []*storage.BucketAttrs

	it := client.Buckets(conf.Context, conf.Project)
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, pair)
	}

	return buckets, nil
}

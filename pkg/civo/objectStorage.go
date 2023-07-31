/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package civo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/civo/civogo"
	"github.com/kubefirst/runtime/pkg/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

// CreateStorageBucket creates an object storage bucket
func (c *CivoConfiguration) CreateStorageBucket(accessKeyId string, bucketName string, region string) (civogo.ObjectStore, error) {
	bucket, err := c.Client.NewObjectStore(&civogo.CreateObjectStoreRequest{
		Name:        bucketName,
		Region:      region,
		AccessKeyID: accessKeyId,
		MaxSizeGB:   500,
	})
	if err != nil {
		return civogo.ObjectStore{}, err
	}

	return *bucket, nil
}

// DeleteStorageBucket deletes an object storage bucket
func (c *CivoConfiguration) DeleteStorageBucket(bucketName string) error {
	objsts, err := c.Client.ListObjectStores()
	if err != nil {
		return err
	}

	var bucketID string

	for _, objst := range objsts.Items {
		if objst.Name == bucketName {
			bucketID = objst.ID
		}
	}

	if bucketID == "" {
		return fmt.Errorf("bucket %s not found", bucketName)
	}

	_, err = c.Client.DeleteObjectStore(bucketID)
	if err != nil {
		return fmt.Errorf("error deleting object store %s: %s", bucketName, err)
	}

	return nil
}

// GetAccessCredentials creates object store access credentials if they do not exist and returns them if they do
func (c *CivoConfiguration) GetAccessCredentials(credentialName string, region string) (civogo.ObjectStoreCredential, error) {
	creds, err := c.checkKubefirstCredentials(credentialName, region)
	if err != nil {
		log.Info().Msg(err.Error())
	}

	if creds == (civogo.ObjectStoreCredential{}) {
		log.Info().Msgf("credential name: %s not found, creating", credentialName)
		creds, err = c.createAccessCredentials(credentialName, region)
		if err != nil {
			return civogo.ObjectStoreCredential{}, err
		}

		for i := 0; i < 30; i++ {
			creds, err = c.getAccessCredentials(creds.ID, region)
			if err != nil {
				return civogo.ObjectStoreCredential{}, err
			}
			if creds.AccessKeyID != "" && creds.ID != "" && creds.Name != "" && creds.SecretAccessKeyID != "" {
				break
			}
			log.Warn().Msg("waiting for civo credentials creation")
			time.Sleep(time.Second * 10)
		}

		if creds.AccessKeyID == "" || creds.ID == "" || creds.Name == "" || creds.SecretAccessKeyID == "" {
			log.Error().Msg("Civo credentials for state bucket in object storage could not be fetched, please try to run again")
			os.Exit(1)
		}
		log.Info().Msgf("created object storage credential %s", credentialName)

		return creds, nil
	}

	return creds, nil
}

// DeleteAccessCredentials deletes object store credentials
func (c *CivoConfiguration) DeleteAccessCredentials(credentialName string, region string) error {
	creds, err := c.checkKubefirstCredentials(credentialName, region)
	if err != nil {
		log.Info().Msg(err.Error())
	}

	_, err = c.Client.DeleteObjectStoreCredential(creds.ID)
	if err != nil {
		return err
	}

	return nil
}

// checkKubefirstCredentials determines whether or not object store credentials exist
func (c *CivoConfiguration) checkKubefirstCredentials(credentialName string, region string) (civogo.ObjectStoreCredential, error) {
	log.Info().Msgf("looking for credential: %s", credentialName)
	remoteCredentials, err := c.Client.ListObjectStoreCredentials()
	if err != nil {
		log.Info().Msg(err.Error())
		return civogo.ObjectStoreCredential{}, err
	}

	var creds civogo.ObjectStoreCredential

	for i, cred := range remoteCredentials.Items {
		if cred.Name == credentialName {
			log.Info().Msgf("found credential: %s", credentialName)
			return remoteCredentials.Items[i], nil
		}
	}

	return creds, err
}

// createAccessCredentials creates access credentials for an object store
func (c *CivoConfiguration) createAccessCredentials(credentialName string, region string) (civogo.ObjectStoreCredential, error) {
	creds, err := c.Client.NewObjectStoreCredential(&civogo.CreateObjectStoreCredentialRequest{
		Name:   credentialName,
		Region: region,
	})
	if err != nil {
		log.Info().Msg(err.Error())
	}
	return *creds, nil
}

// getAccessCredentials retrieves an object store's access credentials
func (c *CivoConfiguration) getAccessCredentials(id string, region string) (civogo.ObjectStoreCredential, error) {
	creds, err := c.Client.GetObjectStoreCredential(id)
	if err != nil {
		return civogo.ObjectStoreCredential{}, err
	}
	return *creds, nil
}

// PutClusterObject exports a cluster definition as json and places it in the target object storage bucket
func PutClusterObject(cr *types.StateStoreCredentials, d *types.StateStoreDetails, obj *types.PushBucketObject) error {
	ctx := context.Background()

	// Initialize minio client
	minioClient, err := minio.New(d.Hostname, &minio.Options{
		Creds:  credentials.NewStaticV4(cr.AccessKeyID, cr.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return fmt.Errorf("error initializing minio client: %s", err)
	}

	// Reference for cluster object output file
	object, err := os.Open(obj.LocalFilePath)
	if err != nil {
		return fmt.Errorf("error during object local copy file lookup: %s", err)
	}
	defer object.Close()

	objectStat, err := object.Stat()
	if err != nil {
		return fmt.Errorf("error during object stat: %s", err)
	}

	// Put
	_, err = minioClient.PutObject(
		ctx,
		d.Name,
		obj.RemoteFilePath,
		object,
		objectStat.Size(),
		minio.PutObjectOptions{ContentType: obj.ContentType},
	)
	if err != nil {
		return fmt.Errorf("error during object put: %s", err)
	}
	log.Info().Msgf("uploaded cluster object %s to state store bucket %s successfully", obj.LocalFilePath, d.Name)

	return nil
}

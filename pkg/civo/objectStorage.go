/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package civo

import (
	"fmt"

	"github.com/civo/civogo"
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

		// Try 10 times to create credentials
		// If they don't create, something in Civo is broken and the operation will need
		// to be run again
		for i := 0; i < 10; i++ {
			// Verify all credentials fields are present
			var civoCredsFailureMessage string

			switch {
			case creds.AccessKeyID == "":
				civoCredsFailureMessage = "when retrieving civo access credentials, AccessKeyID was empty - the credentials will be recreated"
			case creds.ID == "":
				civoCredsFailureMessage = "when retrieving civo access credentials, ID was empty - the credentials will be recreated"
			case creds.Name == "":
				civoCredsFailureMessage = "when retrieving civo access credentials, Name was empty - the credentials will be recreated"
			case creds.SecretAccessKeyID == "":
				civoCredsFailureMessage = "when retrieving civo access credentials, SecretAccessKeyID was empty - the credentials will be recreated"
			}

			if civoCredsFailureMessage != "" {
				// Creds failed to properly parse, so remove them
				log.Warn().Msgf(civoCredsFailureMessage)

				err := c.DeleteAccessCredentials(credentialName, region)
				if err != nil {
					return civogo.ObjectStoreCredential{}, err
				}
			} else {
				// Continue if all fields were present
				break
			}
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

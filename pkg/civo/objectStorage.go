/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package civo

import (
	"github.com/civo/civogo"
	"github.com/rs/zerolog/log"
)

func CreateStorageBucket(civoToken string, accessKeyId string, bucketName string, region string) (civogo.ObjectStore, error) {
	client, err := civogo.NewClient(civoToken, region)
	if err != nil {
		log.Info().Msg(err.Error())
		return civogo.ObjectStore{}, err
	}

	bucket, err := client.NewObjectStore(&civogo.CreateObjectStoreRequest{
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

// todo refactor or remove this internal library and use the native client. functionality. see next todo client.
func GetAccessCredentials(civoToken string, credentialName string, region string) (civogo.ObjectStoreCredential, error) {
	creds, err := checkKubefirstCredentials(civoToken, credentialName, region)
	if err != nil {
		log.Info().Msg(err.Error())
	}

	if creds == (civogo.ObjectStoreCredential{}) {
		log.Info().Msgf("credential name: %s not found, creating", credentialName)
		creds, err = createAccessCredentials(civoToken, credentialName, region)
		if err != nil {
			return civogo.ObjectStoreCredential{}, err
		}

		creds, err = getAccessCredentials(civoToken, creds.ID, region)
		if err != nil {
			return civogo.ObjectStoreCredential{}, err
		}

		log.Info().Msgf("created object storage credential %s", credentialName)
		return creds, nil
	}

	return creds, nil
}

func DeleteAccessCredentials(civoToken string, credentialName string, region string) error {
	client, err := civogo.NewClient(civoToken, region)
	if err != nil {
		log.Info().Msg(err.Error())
		return err
	}

	creds, err := checkKubefirstCredentials(civoToken, credentialName, region)
	if err != nil {
		log.Info().Msg(err.Error())
	}

	_, err = client.DeleteObjectStoreCredential(creds.ID)
	if err != nil {
		return err
	}

	return nil
}

func checkKubefirstCredentials(civoToken string, credentialName string, region string) (civogo.ObjectStoreCredential, error) {
	client, err := civogo.NewClient(civoToken, region)
	if err != nil {
		log.Info().Msg(err.Error())
		return civogo.ObjectStoreCredential{}, err
	}

	// todo client.FindObjectStoreCredential()
	log.Info().Msgf("looking for credential: %s", credentialName)
	remoteCredentials, err := client.ListObjectStoreCredentials()
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

// todo client.NewObjectStoreCredential()
func createAccessCredentials(civoToken string, credentialName string, region string) (civogo.ObjectStoreCredential, error) {
	client, err := civogo.NewClient(civoToken, region)
	if err != nil {
		log.Info().Msg(err.Error())
		return civogo.ObjectStoreCredential{}, err
	}
	creds, err := client.NewObjectStoreCredential(&civogo.CreateObjectStoreCredentialRequest{
		Name:   credentialName,
		Region: region,
	})
	if err != nil {
		log.Info().Msg(err.Error())
	}
	return *creds, nil
}

func getAccessCredentials(civoToken string, id string, region string) (civogo.ObjectStoreCredential, error) {
	client, err := civogo.NewClient(civoToken, region)
	if err != nil {
		log.Info().Msg(err.Error())
		return civogo.ObjectStoreCredential{}, err
	}

	creds, err := client.GetObjectStoreCredential(id)
	if err != nil {
		return civogo.ObjectStoreCredential{}, err
	}
	return *creds, nil
}

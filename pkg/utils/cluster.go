package pkg

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kubefirst/runtime/pkg/k8s"
	"github.com/kubefirst/runtime/pkg/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

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

// ExportCluster port forward to the kubefirst-api and calls /cluster/import to restore the database
func ExportCluster(kcfg types.KubernetesClient, cl types.Cluster) error {
	//* kubefirst api port-forward
	kubefirstApiStopChannel := make(chan struct{}, 1)
	defer func() {
		close(kubefirstApiStopChannel)
	}()
	k8s.OpenPortForwardPodWrapper(
		kcfg.Clientset,
		kcfg.RestConfig,
		"kubefirst-console-kubefirst-api",
		"kubefirst",
		8081,
		8085,
		kubefirstApiStopChannel,
	)

	importUrl := "http://localhost:8085/api/v1/cluster/import"

	importObject := types.ImportClusterRequest{
		ClusterName:           cl.ClusterName,
		CloudRegion:           cl.CloudRegion,
		CloudProvider:         cl.CloudProvider,
		StateStoreCredentials: cl.StateStoreCredentials,
		StateStoreDetails:     cl.StateStoreDetails,
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient := http.Client{Transport: customTransport}

	payload, err := json.Marshal(importObject)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, importUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.Info().Msgf("unable to import cluster %s", res.Status)
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Info().Msgf("Import: %s", string(body))

	return nil
}

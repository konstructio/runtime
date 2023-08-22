/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package digitalocean

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/kubefirst/runtime/pkg/k8s"
	providerConfig "github.com/kubefirst/runtime/pkg/providerConfigs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BootstrapDigitaloceanMgmtCluster(
	digitalOceanToken string,
	kubeconfigPath string,
	gitProvider string,
	gitUser string,
	cloudflareAPIToken string,
	destinationGitopsRepoURL string,
	gitProtocol string,
	dnsProvider string,
) error {
	clientset, err := k8s.GetClientSet(kubeconfigPath)
	if err != nil {
		log.Info().Msg("error getting kubernetes clientset")
	}

	err = providerConfig.BootstrapMgmtCluster(
		clientset,
		gitProvider,
		gitUser,
		destinationGitopsRepoURL,
		gitProtocol,
		cloudflareAPIToken,
		digitalOceanToken, //AWS has no authentication method because we use roles
		dnsProvider,
	)

	//Create cloud specific secrets
	createSecrets := []*v1.Secret{}
	for _, secret := range createSecrets {
		_, err := clientset.CoreV1().Secrets(secret.ObjectMeta.Namespace).Get(context.TODO(), secret.ObjectMeta.Name, metav1.GetOptions{})
		if err == nil {
			log.Info().Msgf("kubernetes secret %s/%s already created - skipping", secret.Namespace, secret.Name)
		} else if strings.Contains(err.Error(), "not found") {
			_, err = clientset.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			if err != nil {
				log.Error().Msgf("error creating kubernetes secret %s/%s: %s", secret.Namespace, secret.Name, err)
				return err
			}
			log.Info().Msgf("created kubernetes secret: %s/%s", secret.Namespace, secret.Name)
		}
	}

	return err
}

/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/kubefirst/runtime/pkg/k8s"
	providerConfig "github.com/kubefirst/runtime/pkg/providerConfigs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (conf *AWSConfiguration) BootstrapAwsMgmtCluster(
	kubeconfigPath string,
	gitProvider string,
	gitUser string,
	cloudflareAPIToken string,
	destinationGitopsRepoURL string,
	gitProtocol string,
	clientset *kubernetes.Clientset,
	ecrFlag bool,
	containerRegistryURL string,
	dnsProvider string,
	cloudProvider string,
) error {

	err := providerConfig.BootstrapMgmtCluster(
		clientset,
		gitProvider,
		gitUser,
		destinationGitopsRepoURL,
		gitProtocol,
		cloudflareAPIToken,
		"", //AWS has no authentication method because we use roles
		dnsProvider,
		CloudProvider,
		(fmt.Sprintf(os.Getenv(fmt.Sprintf("%s_TOKEN", strings.ToUpper(gitProvider))))),
		viper.GetString("kbot.private-key"),
	)

	//Create cloud specific secrets
	createSecrets := []*v1.Secret{}

	for _, secret := range createSecrets {
		_, err := clientset.CoreV1().Secrets(secret.ObjectMeta.Namespace).Get(context.TODO(), secret.ObjectMeta.Name, metav1.GetOptions{})
		if err == nil {
			log.Info().Msgf("kubernetes secret %s/%s already created - skipping", secret.Namespace, secret.Name)
		} else if strings.Contains(err.Error(), "not found") {
			err := k8s.CreateSecretV2(clientset, secret)
			if err != nil {
				log.Error().Msgf("error creating kubernetes secret %s/%s: %s", secret.Namespace, secret.Name, err)
				return err
			}
			log.Info().Msgf("created kubernetes secret: %s/%s", secret.Namespace, secret.Name)
		}
	}

	log.Info().Msg("secret create for argocd to connect to gitops repo")

	//flag out the ecr token

	if ecrFlag {
		ecrToken, err := conf.GetECRAuthToken()
		if err != nil {
			return err
		}
		dockerConfigString := fmt.Sprintf(`{"auths": {"%s": {"auth": "%s"}}}`, containerRegistryURL, ecrToken)
		dockerCfgSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "docker-config", Namespace: "argo"},
			Data:       map[string][]byte{"config.json": []byte(dockerConfigString)},
			Type:       "Opaque",
		}
		_, err = clientset.CoreV1().Secrets(dockerCfgSecret.ObjectMeta.Namespace).Create(context.TODO(), dockerCfgSecret, metav1.CreateOptions{})
		if err != nil {
			log.Info().Msgf("error creating kubernetes secret %s/%s: %s", dockerCfgSecret.Namespace, dockerCfgSecret.Name, err)
			return err
		}
	}

	return err
}

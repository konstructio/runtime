/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package ssh

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/caarlos0/sshmarshal"
	goGitSsh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kubefirst/runtime/pkg/gitlab"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func CreateSshKeyPair() (string, string, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	ecdsaPublicKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return "", "", err
	}

	pemPrivateKey, err := sshmarshal.MarshalPrivateKey(privKey, "kubefirst key")
	if err != nil {
		return "", "", err
	}

	privateKey := string(pem.EncodeToMemory(pemPrivateKey))
	publicKey := string(ssh.MarshalAuthorizedKey(ecdsaPublicKey))

	return privateKey, publicKey, nil
}

func PublicKeyV2() (*goGitSsh.PublicKeys, error) {
	var publicKey *goGitSsh.PublicKeys
	publicKey, err := goGitSsh.NewPublicKeys("kube1st", []byte(viper.GetString("kubefirst.bot.private-key")), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

// EvalSSHKey
func EvalSSHKey(req *EvalSSHKeyRequest) error {
	// For GitLab, we currently need to add an ssh key to the authenticating user
	if req.GitProvider == "gitlab" {
		gitlabClient, err := gitlab.NewGitLabClient(req.GitToken, req.GitlabGroupFlag)
		if err != nil {
			return err
		}
		keys, err := gitlabClient.GetUserSSHKeys()
		if err != nil {
			log.Fatal().Msgf("unable to check for ssh keys in gitlab: %s", err.Error())
		}

		var keyName = "kbot-ssh-key"
		var createKey bool = false
		for _, key := range keys {
			if key.Title == keyName {
				if strings.Contains(key.Key, strings.TrimSuffix(viper.GetString("kbot.public-key"), "\n")) {
					log.Info().Msgf("ssh key %s already exists and key is up to date, continuing", keyName)
				} else {
					log.Warn().Msgf("ssh key %s already exists and key data has drifted - it will be recreated", keyName)
					err := gitlabClient.DeleteUserSSHKey(keyName)
					if err != nil {
						return fmt.Errorf("error deleting gitlab user ssh key %s: %s", keyName, err)
					}
					createKey = true
				}
			}
		}
		if createKey {
			log.Info().Msgf("creating ssh key %s...", keyName)
			err := gitlabClient.AddUserSSHKey(keyName, viper.GetString("kbot.public-key"))
			if err != nil {
				log.Fatal().Msgf("error adding ssh key %s: %s", keyName, err.Error())
			}
			viper.Set("kbot.gitlab-user-based-ssh-key-title", keyName)
			viper.WriteConfig()
		}
	}

	return nil
}

/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package terraform

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"

	"github.com/kubefirst/runtime/configs"
	"github.com/kubefirst/runtime/pkg"
)

func initActionAutoApprove(tfAction, tfEntrypoint string, tfEnvs map[string]string) error {
	config := configs.ReadConfig()
	log.Printf("initActionAutoApprove - action: %s entrypoint: %s", tfAction, tfEntrypoint)

	err := os.Chdir(tfEntrypoint)
	if err != nil {
		log.Info().Msg("error: could not change to directory " + tfEntrypoint)
		return err
	}
	err = pkg.ExecShellWithVars(tfEnvs, config.TerraformClientPath, "init", "-force-copy")
	if err != nil {
		log.Printf("error: terraform init for %s failed: %s", tfEntrypoint, err)
		return err
	}

	err = pkg.ExecShellWithVars(tfEnvs, config.TerraformClientPath, tfAction, "-auto-approve")
	if err != nil {
		log.Printf("error: terraform %s -auto-approve for %s failed %s", tfAction, tfEntrypoint, err)
		return err
	}
	os.RemoveAll(fmt.Sprintf("%s/.terraform/", tfEntrypoint))
	os.Remove(fmt.Sprintf("%s/.terraform.lock.hcl", tfEntrypoint))
	return nil
}

func InitApplyAutoApprove(tfEntrypoint string, tfEnvs map[string]string) error {
	tfAction := "apply"
	err := initActionAutoApprove(tfAction, tfEntrypoint, tfEnvs)
	if err != nil {
		return err
	}
	return nil
}

func InitDestroyAutoApprove(tfEntrypoint string, tfEnvs map[string]string) error {
	tfAction := "destroy"
	err := initActionAutoApprove(tfAction, tfEntrypoint, tfEnvs)
	if err != nil {
		return err
	}
	return nil
}

// todo need to write something that outputs -json type and can get multiple values
func OutputSingleValue(directory, tfEntrypoint, outputName string) {
	config := configs.ReadConfig()
	os.Chdir(directory)

	var tfOutput bytes.Buffer
	tfOutputCmd := exec.Command(config.TerraformClientPath, "output", outputName)
	tfOutputCmd.Stdout = &tfOutput
	tfOutputCmd.Stderr = os.Stderr
	err := tfOutputCmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to call tfOutputCmd.Run()")
	}

	log.Print("tfOutput is: ", tfOutput.String())
}

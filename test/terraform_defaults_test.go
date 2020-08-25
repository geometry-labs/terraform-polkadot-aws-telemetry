package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTerraformDefaults(t *testing.T) {
	t.Parallel()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/defaults")

	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.Destroy(t, terraformOptions)

		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)
		aws.DeleteEC2KeyPair(t, keyPair)
	})

	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions, keyPair := configureTerraformOptions(t, exampleFolder)
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)
		test_structure.SaveEc2KeyPair(t, exampleFolder, keyPair)

		terraform.InitAndApply(t, terraformOptions)
	})

	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)

		testEndpoint(t, terraformOptions)
	})
}


func configureTerraformOptions(t *testing.T, exampleFolder string) (*terraform.Options, *aws.Ec2Keypair) {

	uniqueID := random.UniqueId()
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	keyPairName := fmt.Sprintf("terratest-ssh-example-%s", uniqueID)
	keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, keyPairName)

	privateKeyPath := path.Join(exampleFolder, "id_rsa")
	err := ioutil.WriteFile(privateKeyPath, []byte(keyPair.PrivateKey), 0600)
	if err != nil {
		panic(err)
	}

	publicKeyPath := path.Join(exampleFolder, "id_rsa.pub")
	err = ioutil.WriteFile(publicKeyPath, []byte(keyPair.PublicKey), 0644)
	if err != nil {
		panic(err)
	}

	terraformOptions := &terraform.Options{
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"aws_region":    awsRegion,
			"public_key_path":    publicKeyPath,
			"private_key_path": privateKeyPath,
		},
	}

	return terraformOptions, keyPair
}

func testEndpoint(t *testing.T, terraformOptions *terraform.Options) {

	endpointIP := terraform.Output(t, terraformOptions, "dns_name")

	expectedStatus := "200"
	url := fmt.Sprintf("http://%s:3000", endpointIP)

	description := fmt.Sprintf("curl to LB %s with error command", endpointIP)
	maxRetries := 60
	timeBetweenRetries := 2 * time.Second

	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {

		outputStatus, _, err := http_helper.HTTPDoE(t, "GET", url, nil, nil,nil)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(strconv.Itoa(outputStatus)) != expectedStatus {
			return "", fmt.Errorf("expected '%s' but got '%s'", expectedStatus, strconv.Itoa(outputStatus))
		}

		return "", nil
	})
}

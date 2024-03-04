/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package aws

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// TXTRecord stores Route53 TXT record data
type TXTRecord struct {
	Name          string
	Value         string
	SetIdentifier *string
	Weight        *int64
	TTL           int64
}

// ARecord stores Route53 A record data
type ARecord struct {
	Name        string
	RecordType  string
	TTL         *int64
	AliasTarget *route53Types.AliasTarget
}

func NewAwsV2(region string) aws.Config {
	// todo these should also be supported flags
	profile := os.Getenv("AWS_PROFILE")

	awsClient, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		log.Error().Msg("unable to create aws client")
	}

	return awsClient
}

func NewAwsV3(region string, accessKeyID string, secretAccessKey string, sessionToken string) aws.Config {
	awsClient, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			sessionToken,
		)),
	)
	if err != nil {
		log.Error().Msg("unable to create aws client")
	}

	return awsClient
}

func NewAwsClientServiceAccountTokenV1(region string, accessKeyID string, secretAccessKey string, sessionToken string) aws.Config {

	region = os.Getenv("AWS_REGION")
	fmt.Println("Region: ", region)

	roleArn := os.Getenv("AWS_ROLE_ARN")
	fmt.Println("Role ARN: ", roleArn)

	tokenFilePath := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	fmt.Println("Token File Path: ", tokenFilePath)

	// Get the service account token from the pod's projected volume
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		panic(err.Error())
	}

	stsEndpoints := os.Getenv("AWS_STS_REGIONAL_ENDPOINTS")
	fmt.Println("STS Endpoints: ", stsEndpoints)
	awsClient, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			string(token),
			"",
			string(token),
		)),
	)
	if err != nil {
		log.Error().Msg("unable to create aws client")
	}

	return awsClient
}

// GetRegions lists all available regions
func (conf *AWSConfiguration) GetRegions(region string) ([]string, error) {
	var regionList []string

	ec2Client := ec2.NewFromConfig(conf.Config)

	regions, err := ec2Client.DescribeRegions(context.Background(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return []string{}, fmt.Errorf("error listing regions: %s", err)
	}

	for _, region := range regions.Regions {
		regionList = append(regionList, *region.RegionName)
	}

	return regionList, nil
}

func (conf *AWSConfiguration) ListInstanceSizesForRegion() ([]string, error) {

	ec2Client := ec2.NewFromConfig(conf.Config)

	sizes, err := ec2Client.DescribeInstanceTypeOfferings(context.Background(), &ec2.DescribeInstanceTypeOfferingsInput{})

	if err != nil {
		return nil, err
	}

	var instanceNames []string
	for _, size := range sizes.InstanceTypeOfferings {
		instanceNames = append(instanceNames, string(size.InstanceType))
	}

	return instanceNames, nil
}

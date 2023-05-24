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

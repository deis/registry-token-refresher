package credentials

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ecrDockerCredProvider struct {
	params Parameters
}

func (d *ecrDockerCredProvider) GetDockerCredentials() (DockerConfig, error) {
	accessKey := d.params["accesskey"]
	if accessKey == nil {
		accessKey = ""
	}
	secretKey := d.params["secretkey"]
	if secretKey == nil {
		secretKey = ""
	}
	regionName := d.params["region"]
	if regionName == nil || fmt.Sprint(regionName) == "" {
		return DockerConfig{}, fmt.Errorf("No region parameter provided")
	}
	region := fmt.Sprint(regionName)
	registryID := d.params["registryid"]
	if registryID == nil {
		registryID = ""
	}
	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     fmt.Sprint(accessKey),
				SecretAccessKey: fmt.Sprint(secretKey),
			},
		},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
		&ec2rolecreds.EC2RoleProvider{Client: ec2metadata.New(session.New())},
	})
	awsConfig := aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(region)
	svc := ecr.New(session.New(awsConfig))

	var authInput *ecr.GetAuthorizationTokenInput
	if registryID == "" {
		authInput = &ecr.GetAuthorizationTokenInput{}
	} else {
		authInput = &ecr.GetAuthorizationTokenInput{
			RegistryIds: []*string{
				aws.String(fmt.Sprint(registryID)),
			},
		}
	}

	resp, err := svc.GetAuthorizationToken(authInput)
	if err != nil {
		return DockerConfig{}, err
	}
	authData := resp.AuthorizationData[0]
	hostname := d.params["hostname"]
	if hostname == nil || fmt.Sprint(hostname) == "" {
		hostname = *authData.ProxyEndpoint
	}
	dockerConfig := DockerConfig{Token: *authData.AuthorizationToken,
		ExpiresAt: *authData.ExpiresAt,
		Hostname:  fmt.Sprint(hostname)}
	return dockerConfig, nil
}

func (d *ecrDockerCredProvider) GetRefreshTime() time.Duration {
	defaultRefreshTime := d.params["defaultRefreshTime"]
	if defaultRefreshTime != nil {
		return time.Minute * time.Duration(defaultRefreshTime.(int))
	}
	return time.Hour * 11
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

type InstanceIdentity struct {
	Ec2InstanceId        string
	ContainerInstanceArn string
	Cluster              string
	Region               string
	AvailabilityZone     string
	IAMCredentials       credentials.Credentials
}

type EcsAgentMetadata struct {
	Cluster              string `json:"Cluster"`
	ContainerInstanceArn string `json:"ContainerInstanceArn"`
	Version              string `json:"Version"`
}

func NewInstanceIdentity() (*InstanceIdentity, error) {
	// Create metadata session and instance
	metasess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to create metadata session: %s", err)
	}

	metaclient := ec2metadata.New(metasess)

	ec2, err := metaclient.GetInstanceIdentityDocument()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to get instance identity: %s", err)
	}

	// Gather instance placement info and IAM creds
	providers := []credentials.Provider{
		&credentials.EnvProvider{},
		&ec2rolecreds.EC2RoleProvider{
			Client: metaclient,
		},
	}

	creds := credentials.NewChainCredentials(providers)

	// Get ECS container instance ID
	ecs, err := NewEcsAgentMetadata()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to get container instance id: %s", err)
	}

	ident := InstanceIdentity{
		Ec2InstanceId:        ec2.InstanceID,
		ContainerInstanceArn: ecs.ContainerInstanceArn,
		Cluster:              ecs.Cluster,
		Region:               ec2.Region,
		AvailabilityZone:     ec2.AvailabilityZone,
		IAMCredentials:       *creds,
	}

	return &ident, nil
}

func NewEcsAgentMetadata() (*EcsAgentMetadata, error) {
	// Get container instance ID from ECS agent
	res, err := http.Get("http://localhost:51678/v1/metadata")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to connect to ECS agent: %s", err)
	}

	metadata := EcsAgentMetadata{}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to unmarshal ECS metadata: %s", err)
	}

	return &metadata, nil
}

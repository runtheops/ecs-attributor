package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	InstanceAttributes = flag.String("a", "", "comma-separated list of key=value attribute pairs")
)

func main() {
	if err := RunAttributor(); err != nil {
		log.Fatal(err)
	}
}

func RunAttributor() error {
	flag.Parse()

	ident, err := NewInstanceIdentity()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Instance identity: %s\n", ident.ContainerInstanceArn)

	attr, err := ParseAttributes(*InstanceAttributes, ident.ContainerInstanceArn)
	if err != nil {
		return err
	}

	log.Println("[INFO] Putting attributes...")
	err = PutAttributes(ident, attr)

	log.Println("[INFO] Done")

	return err
}

func ParseAttributes(a string, t string) ([]*ecs.Attribute, error) {
	ecsa := []*ecs.Attribute{}

	for _, attr := range strings.Split(a, ",") {
		p := strings.Split(attr, "=")
		if len(p) < 2 {
			return nil, fmt.Errorf("[ERROR] Attributes string is malformed.")
		}

		ecsa = append(ecsa, &ecs.Attribute{Name: &p[0], Value: &p[1], TargetId: &t})
	}

	return ecsa, nil
}

func PutAttributes(i *InstanceIdentity, a []*ecs.Attribute) error {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      &i.Region,
			Credentials: &i.IAMCredentials,
		},
	)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to create an API session:", err)
	}

	svc := ecs.New(sess)

	params := &ecs.PutAttributesInput{
		Attributes: a,
		Cluster:    aws.String(i.Cluster),
	}

	_, err = svc.PutAttributes(params)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to put instance attribute:", err)
	}

	return nil
}

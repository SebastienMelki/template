package vpc

import (
	"fmt"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Result contains the outputs of the VPC component.
type Result struct {
	VPCID pulumi.IDOutput
}

// New creates a VPC for the project.
func New(ctx *pulumi.Context, project, environment string) (*Result, error) {
	name := fmt.Sprintf("%s-%s-vpc", project, environment)

	v, err := digitalocean.NewVpc(ctx, name, &digitalocean.VpcArgs{
		Name:        pulumi.String(name),
		Region:      pulumi.String("ams3"), // TODO: Make configurable
		IpRange:     pulumi.String("10.0.0.0/16"),
		Description: pulumi.Sprintf("%s %s VPC", project, environment),
	})
	if err != nil {
		return nil, err
	}

	return &Result{VPCID: v.ID()}, nil
}

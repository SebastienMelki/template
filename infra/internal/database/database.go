package database

import (
	"fmt"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/vpc"
)

// Result contains the outputs of the database component.
type Result struct {
	Host     pulumi.StringOutput
	Port     pulumi.IntOutput
	Name     pulumi.StringOutput
	User     pulumi.StringOutput
	Password pulumi.StringOutput
	URI      pulumi.StringOutput
}

// New creates a managed PostgreSQL database.
func New(ctx *pulumi.Context, project, environment string, vpcResult *vpc.Result) (*Result, error) {
	name := fmt.Sprintf("%s-%s-db", project, environment)

	cluster, err := digitalocean.NewDatabaseCluster(ctx, name, &digitalocean.DatabaseClusterArgs{
		Name:                 pulumi.String(name),
		Engine:               pulumi.String("pg"),
		Version:              pulumi.String("17"),
		Size:                 pulumi.String("db-s-1vcpu-1gb"),
		Region:               pulumi.String("ams3"),
		NodeCount:            pulumi.Int(1),
		PrivateNetworkUuid:   vpcResult.VPCID.ToStringOutput(),
	})
	if err != nil {
		return nil, err
	}

	db, err := digitalocean.NewDatabaseDb(ctx, fmt.Sprintf("%s-db", name), &digitalocean.DatabaseDbArgs{
		ClusterId: cluster.ID(),
		Name:      pulumi.String(project),
	})
	if err != nil {
		return nil, err
	}

	return &Result{
		Host:     cluster.PrivateHost,
		Port:     cluster.Port,
		Name:     db.Name,
		User:     cluster.User,
		Password: cluster.Password,
		URI:      cluster.PrivateUri,
	}, nil
}

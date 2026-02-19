package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/backend"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/dns"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/ingress"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/kubernetes"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/registry"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/vpc"
	infradb "github.com/{{ORG}}/{{PROJECT}}/infra/internal/database"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		project := cfg.Require("project")
		environment := cfg.Require("environment")
		domain := cfg.Require("domain")

		// VPC
		vpcResult, err := vpc.New(ctx, project, environment)
		if err != nil {
			return err
		}

		// Database
		dbResult, err := infradb.New(ctx, project, environment, vpcResult)
		if err != nil {
			return err
		}

		// Kubernetes
		k8sResult, err := kubernetes.New(ctx, project, environment, vpcResult)
		if err != nil {
			return err
		}

		// Container Registry
		registryResult, err := registry.New(ctx, project, k8sResult)
		if err != nil {
			return err
		}

		// Ingress (Traefik)
		ingressResult, err := ingress.New(ctx, project, environment, k8sResult)
		if err != nil {
			return err
		}

		// DNS
		_, err = dns.New(ctx, domain, ingressResult)
		if err != nil {
			return err
		}

		// Backend Deployment
		_, err = backend.New(ctx, project, environment, k8sResult, registryResult, dbResult)
		if err != nil {
			return err
		}

		// Outputs
		ctx.Export("vpcId", vpcResult.VPCID)
		ctx.Export("clusterId", k8sResult.ClusterID)
		ctx.Export("clusterEndpoint", k8sResult.Endpoint)
		ctx.Export("databaseHost", dbResult.Host)
		ctx.Export("traefikIP", ingressResult.LoadBalancerIP)

		return nil
	})
}

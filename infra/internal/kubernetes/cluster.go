package kubernetes

import (
	"fmt"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/vpc"
)

// Result contains the outputs of the Kubernetes cluster.
type Result struct {
	ClusterID  pulumi.IDOutput
	Endpoint   pulumi.StringOutput
	Kubeconfig pulumi.StringOutput
}

// New creates a managed Kubernetes cluster.
func New(ctx *pulumi.Context, project, environment string, vpcResult *vpc.Result) (*Result, error) {
	name := fmt.Sprintf("%s-%s-k8s", project, environment)

	cluster, err := digitalocean.NewKubernetesCluster(ctx, name, &digitalocean.KubernetesClusterArgs{
		Name:    pulumi.String(name),
		Region:  pulumi.String("ams3"),
		Version: pulumi.String("1.31.1-do.5"), // TODO: Use latest
		VpcUuid: vpcResult.VPCID.ToStringOutput(),
		NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
			Name:      pulumi.String(fmt.Sprintf("%s-pool", name)),
			Size:      pulumi.String("s-2vcpu-2gb"),
			AutoScale: pulumi.Bool(true),
			MinNodes:  pulumi.Int(2),
			MaxNodes:  pulumi.Int(4),
		},
	})
	if err != nil {
		return nil, err
	}

	kubeconfig := cluster.KubeConfigs.Index(pulumi.Int(0)).RawConfig()

	return &Result{
		ClusterID:  cluster.ID(),
		Endpoint:   cluster.Endpoint,
		Kubeconfig: kubeconfig,
	}, nil
}

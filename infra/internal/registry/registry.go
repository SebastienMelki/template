package registry

import (
	"fmt"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/kubernetes"
)

// Result contains the outputs of the container registry.
type Result struct {
	Endpoint   pulumi.StringOutput
	ServerURL  pulumi.StringOutput
}

// New creates a container registry and connects it to the Kubernetes cluster.
func New(ctx *pulumi.Context, project string, k8sResult *kubernetes.Result) (*Result, error) {
	name := fmt.Sprintf("%s-registry", project)

	reg, err := digitalocean.NewContainerRegistry(ctx, name, &digitalocean.ContainerRegistryArgs{
		Name:                 pulumi.String(name),
		SubscriptionTierSlug: pulumi.String("starter"),
	})
	if err != nil {
		return nil, err
	}

	// Connect registry to Kubernetes cluster
	_, err = digitalocean.NewContainerRegistryDockerCredentials(ctx, fmt.Sprintf("%s-creds", name), &digitalocean.ContainerRegistryDockerCredentialsArgs{
		RegistryName: reg.Name,
	})
	if err != nil {
		return nil, err
	}

	return &Result{
		Endpoint:  reg.Endpoint,
		ServerURL: reg.ServerUrl,
	}, nil
}

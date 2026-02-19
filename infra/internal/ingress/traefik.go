package ingress

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/kubernetes"
)

// Result contains the outputs of the ingress controller.
type Result struct {
	LoadBalancerIP pulumi.StringOutput
}

// New deploys Traefik as an ingress controller with Let's Encrypt TLS.
func New(ctx *pulumi.Context, project, environment string, k8sResult *kubernetes.Result) (*Result, error) {
	name := fmt.Sprintf("%s-%s-traefik", project, environment)

	// Determine Let's Encrypt server based on environment
	certServer := "https://acme-v02.api.letsencrypt.org/directory"
	if environment == "dev" {
		certServer = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	release, err := helm.NewRelease(ctx, name, &helm.ReleaseArgs{
		Chart:     pulumi.String("traefik"),
		Version:   pulumi.String("28.0.0"),
		Namespace: pulumi.String("traefik"),
		CreateNamespace: pulumi.Bool(true),
		RepositoryOpts: &helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://traefik.github.io/charts"),
		},
		Values: pulumi.Map{
			"ports": pulumi.Map{
				"web": pulumi.Map{
					"redirectTo": pulumi.Map{
						"port": pulumi.String("websecure"),
					},
				},
			},
			"certificatesResolvers": pulumi.Map{
				"letsencrypt": pulumi.Map{
					"acme": pulumi.Map{
						"email":   pulumi.String("admin@{{DOMAIN}}"),
						"storage": pulumi.String("/data/acme.json"),
						"caServer": pulumi.String(certServer),
						"httpChallenge": pulumi.Map{
							"entryPoint": pulumi.String("web"),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_ = release

	// TODO: Extract LoadBalancer IP from Traefik service
	// This requires querying the Kubernetes service after deployment
	return &Result{
		LoadBalancerIP: pulumi.String("pending").ToStringOutput(),
	}, nil
}

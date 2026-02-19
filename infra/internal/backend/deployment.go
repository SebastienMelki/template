package backend

import (
	"fmt"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	infradb "github.com/{{ORG}}/{{PROJECT}}/infra/internal/database"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/kubernetes"
	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/registry"
)

// Result contains the outputs of the backend deployment.
type Result struct {
	ServiceName pulumi.StringOutput
}

// New deploys the backend application to Kubernetes.
func New(
	ctx *pulumi.Context,
	project, environment string,
	k8sResult *kubernetes.Result,
	registryResult *registry.Result,
	dbResult *infradb.Result,
) (*Result, error) {
	name := fmt.Sprintf("%s-%s-backend", project, environment)
	namespace := fmt.Sprintf("%s-%s", project, environment)
	labels := pulumi.StringMap{"app": pulumi.String(name)}

	// Namespace
	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
		},
	})
	if err != nil {
		return nil, err
	}

	// Deployment
	_, err = appsv1.NewDeployment(ctx, name, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(name),
			Namespace: ns.Metadata.Name(),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(2),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: labels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: labels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name: pulumi.String("backend"),
							Image: pulumi.Sprintf("%s/%s:latest",
								registryResult.ServerURL,
								project,
							),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8080),
									Name:          pulumi.String("http"),
								},
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8081),
									Name:          pulumi.String("probe"),
								},
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name:  pulumi.String("DATABASE_MODE"),
									Value: pulumi.String("remote"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("DATABASE_HOST"),
									Value: dbResult.Host,
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("DATABASE_NAME"),
									Value: dbResult.Name,
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("DATABASE_USER"),
									Value: dbResult.User,
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("DATABASE_PASSWORD"),
									Value: dbResult.Password,
								},
							},
							LivenessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/healthz"),
									Port: pulumi.Int(8081),
								},
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/readyz"),
									Port: pulumi.Int(8081),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Service
	svc, err := corev1.NewService(ctx, fmt.Sprintf("%s-svc", name), &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(name),
			Namespace: ns.Metadata.Name(),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: labels,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(80),
					TargetPort: pulumi.Int(8080),
					Name:       pulumi.String("http"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &Result{
		ServiceName: svc.Metadata.Name().Elem(),
	}, nil
}

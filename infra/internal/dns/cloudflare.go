package dns

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/{{ORG}}/{{PROJECT}}/infra/internal/ingress"
)

// Result contains the outputs of the DNS component.
type Result struct {
	RecordID pulumi.IDOutput
}

// New creates DNS records pointing to the ingress controller.
func New(ctx *pulumi.Context, domain string, ingressResult *ingress.Result) (*Result, error) {
	// TODO: Look up zone ID from domain name
	// zone, err := cloudflare.LookupZone(ctx, &cloudflare.LookupZoneArgs{Name: &domain})

	record, err := cloudflare.NewRecord(ctx, "wildcard-dns", &cloudflare.RecordArgs{
		ZoneId:  pulumi.String("REPLACE_WITH_ZONE_ID"), // TODO: Use zone lookup
		Name:    pulumi.String("*"),
		Content: ingressResult.LoadBalancerIP,
		Type:    pulumi.String("A"),
		Ttl:     pulumi.Int(300),
		Proxied: pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return &Result{RecordID: record.ID()}, nil
}

package rule

import (
	"fmt"
	"strings"

	"github.com/aquasecurity/defsec/infra"
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/result"
	"github.com/aquasecurity/defsec/severity"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
)

// Rule is a targeted security test which can be applied to terraform templates. It includes the types to run on e.g.
// "resource", and the labels to run on e.g. "aws_s3_bucket".
type Rule struct {
	LegacyID string

	Service   string // EC2
	ShortCode string // ebs-volume-encrypted

	Documentation       RuleDocumentation
	Provider            provider.Provider
	RequiredTypes       []string
	RequiredLabels      []string
	RequiredSources     []string
	DefaultSeverity     severity.Severity
	CheckTerraform      func(result.Set, block.Block, block.Module)
	CheckInfrastructure func(*infra.Context) []*result.Result
}

func (r Rule) ID() string {
	return strings.ToLower(fmt.Sprintf("%s-%s-%s", r.Provider, r.Service, r.ShortCode))
}

func (r Rule) MatchesID(id string) bool {
	return r.LegacyID == id || r.ID() == id
}

type RuleDocumentation struct {

	// Summary is a brief description of the check, e.g. "Unencrypted S3 Bucket"
	Summary string

	// Explanation (markdown) contains reasoning for the check, details on it's value, and remediation info
	Explanation string

	// Impact contains a brief summary of the impact of failing the check
	Impact string

	// Resolution contains a brief summary of the resolution for the failing check
	Resolution string

	// BadExample (hcl) contains Terraform code which would cause the check to fail
	BadExample []string

	// GoodExample (hcl) modifies the BadExample content to cause the check to pass
	GoodExample []string

	// Links are URLs which contain further reading related to the check
	Links []string
}

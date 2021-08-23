package apigateway

// ATTENTION!
// This rule was autogenerated!
// Before making changes, consider updating the generator.

// generator-locked
import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/result"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Provider:  provider.AWSProvider,
		Service:   "api-gateway",
		ShortCode: "enable-tracing",
		Documentation: rule.RuleDocumentation{
			Summary:     "API Gateway must have X-Ray tracing enabled",
			Explanation: `X-Ray tracing enables end-to-end debugging and analysis of all API Gateway HTTP requests.`,
			Impact:      "WIthout full tracing enabled it is difficult to trace the flow of logs",
			Resolution:  "Enable tracing",
			BadExample: []string{`
resource "aws_api_gateway_stage" "bad_example" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.test.id
  deployment_id = aws_api_gateway_deployment.test.id
  xray_tracing_enabled = false
}
`},
			GoodExample: []string{`
resource "aws_api_gateway_stage" "good_example" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.test.id
  deployment_id = aws_api_gateway_deployment.test.id
  xray_tracing_enabled = true
}
`},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage#xray_tracing_enabled",
			},
		},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"aws_api_gateway_stage",
		},
		DefaultSeverity: severity.Low,
		CheckTerraform: func(set result.Set, resourceBlock block.Block, _ block.Module) {
			if xrayTracingEnabledAttr := resourceBlock.GetAttribute("xray_tracing_enabled"); xrayTracingEnabledAttr.IsNil() { // alert on use of default value
				set.AddResult().
					WithDescription("Resource '%s' uses default value for xray_tracing_enabled", resourceBlock.FullName())
			} else if xrayTracingEnabledAttr.IsFalse() {
				set.AddResult().
					WithDescription("Resource '%s' does not have xray_tracing_enabled set to true", resourceBlock.FullName()).
					WithAttribute("")
			}
		},
	})
}

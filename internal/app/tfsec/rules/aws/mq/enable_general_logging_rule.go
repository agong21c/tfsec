package mq

// ATTENTION!
// This rule was autogenerated!
// Before making changes, consider updating the generator.

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
		Service:   "mq",
		ShortCode: "enable-general-logging",
		Documentation: rule.RuleDocumentation{
			Summary:     "MQ Broker should have general logging enabled",
			Explanation: `Logging should be enabled to allow tracing of issues and activity to be investigated more fully. Logs provide additional information and context which is often invalauble during investigation`,
			Impact:      "Without logging it is difficult to trace issues",
			Resolution:  "Enable general logging",
			BadExample: []string{`
resource "aws_mq_broker" "bad_example" {
  broker_name = "example"

  configuration {
    id       = aws_mq_configuration.test.id
    revision = aws_mq_configuration.test.latest_revision
  }

  engine_type        = "ActiveMQ"
  engine_version     = "5.15.0"
  host_instance_type = "mq.t2.micro"
  security_groups    = [aws_security_group.test.id]

  user {
    username = "ExampleUser"
    password = "MindTheGap"
  }
  logs {
    general = false
  }
}
`},
			GoodExample: []string{`
resource "aws_mq_broker" "good_example" {
  broker_name = "example"

  configuration {
    id       = aws_mq_configuration.test.id
    revision = aws_mq_configuration.test.latest_revision
  }

  engine_type        = "ActiveMQ"
  engine_version     = "5.15.0"
  host_instance_type = "mq.t2.micro"
  security_groups    = [aws_security_group.test.id]

  user {
    username = "ExampleUser"
    password = "MindTheGap"
  }
  logs {
    general = true
  }
}
`},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/mq_broker#general",
			},
		},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"aws_mq_broker",
		},
		DefaultSeverity: severity.Low,
		CheckTerraform: func(set result.Set, resourceBlock block.Block, _ block.Module) {
			if generalAttr := resourceBlock.GetBlock("logs").GetAttribute("general"); generalAttr.IsNil() { // alert on use of default value
				set.AddResult().
					WithDescription("Resource '%s' uses default value for logs.general", resourceBlock.FullName())
			} else if generalAttr.IsFalse() {
				set.AddResult().
					WithDescription("Resource '%s' does not have logs.general set to true", resourceBlock.FullName()).
					WithAttribute("")
			}
		},
	})
}

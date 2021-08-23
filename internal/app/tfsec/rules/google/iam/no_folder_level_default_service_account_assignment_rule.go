package iam

import (
	"github.com/aquasecurity/defsec/result"
	"github.com/aquasecurity/defsec/severity"

	"github.com/aquasecurity/defsec/provider"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"

	"github.com/aquasecurity/tfsec/pkg/rule"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Service:   "iam",
		ShortCode: "no-folder-level-default-service-account-assignment",
		Documentation: rule.RuleDocumentation{
			Summary:     "Roles should not be assigned to default service accounts",
			Explanation: `Default service accounts should not be used - consider creating specialised service accounts for individual purposes.`,
			Impact:      "Violation of principal of least privilege",
			Resolution:  "Use specialised service accounts for specific purposes.",
			BadExample: []string{`
resource "google_folder_iam_member" "folder-123" {
	folder = "folder-123"
	role    = "roles/whatever"
	member  = "123-compute@developer.gserviceaccount.com"
}
`,
				`
resource "google_folder_iam_member" "folder-123" {
	folder = "folder-123"
	role    = "roles/whatever"
	member  = "123@appspot.gserviceaccount.com"
}
`, `
data "google_compute_default_service_account" "default" {
}

resource "google_folder_iam_member" "folder-123" {
	folder = "folder-123"
	role    = "roles/whatever"
	member  = data.google_compute_default_service_account.default.id
}
`,
			},
			GoodExample: []string{`
resource "google_service_account" "test" {
	account_id   = "account123"
	display_name = "account123"
}
			  
resource "google_folder_iam_member" "folder-123" {
	folder = "folder-123"
	role    = "roles/whatever"
	member  = "serviceAccount:${google_service_account.test.email}"
}
`,
			},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_folder_iam",
				"",
			},
		},
		Provider:        provider.GoogleProvider,
		RequiredTypes:   []string{"resource"},
		RequiredLabels:  []string{"google_folder_iam_binding", "google_folder_iam_member"},
		DefaultSeverity: severity.Medium,
		CheckTerraform: func(set result.Set, resourceBlock block.Block, module block.Module) {

			if memberAttr := resourceBlock.GetAttribute("member"); memberAttr.IsNotNil() {
				if memberAttr.IsString() {
					if isMemberDefaultServiceAccount(memberAttr.Value().AsString()) {
						set.AddResult().
							WithAttribute("").
							WithDescription("Resource '%s' assigns a role to a default service account.", resourceBlock.FullName())
					}
				} else {
					computeServiceAccounts := module.GetDatasByType("google_compute_default_service_account")
					serviceAccounts := append(computeServiceAccounts, module.GetResourcesByType("google_app_engine_default_service_account")...)
					for _, serviceAccount := range serviceAccounts {
						if memberAttr.References(serviceAccount.Reference()) {
							set.AddResult().
								WithAttribute("").
								WithDescription("Resource '%s' assigns a role to a default service account.", resourceBlock.FullName())
						}
					}
				}
			}

			if membersAttr := resourceBlock.GetAttribute("members"); membersAttr.IsNotNil() {
				for _, member := range membersAttr.ValueAsStrings() {
					if isMemberDefaultServiceAccount(member) {
						set.AddResult().
							WithAttribute("").
							WithDescription("Resource '%s' assigns a role to a default service account.", resourceBlock.FullName())
					}
				}
				computeServiceAccounts := module.GetDatasByType("google_compute_default_service_account")
				serviceAccounts := append(computeServiceAccounts, module.GetResourcesByType("google_app_engine_default_service_account")...)
				for _, serviceAccount := range serviceAccounts {
					if membersAttr.References(serviceAccount.Reference()) {
						set.AddResult().
							WithAttribute("").
							WithDescription("Resource '%s' assigns a role to a default service account.", resourceBlock.FullName())
					}
				}
			}

		},
	})
}

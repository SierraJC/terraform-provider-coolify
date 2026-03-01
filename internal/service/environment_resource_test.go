package service_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccEnvironmentResource(t *testing.T) {
	resName := "coolify_environment.test"
	
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_project" "test" {
					name        = "TerraformAccTestEnvProject"
					description = "Test project for environment"
				}

				resource "coolify_environment" "test" {
					name         = "terraform-test-env"
					project_uuid = coolify_project.test.uuid
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "terraform-test-env"),
					resource.TestCheckResourceAttrSet(resName, "project_uuid"),
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{ // ImportState testing
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resName]
					projectUuid := rs.Primary.Attributes["project_uuid"]
					envUuid := rs.Primary.Attributes["uuid"]
					return projectUuid + "/" + envUuid, nil
				},
			},
			{ // Import by name testing
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resName]
					projectUuid := rs.Primary.Attributes["project_uuid"]
					envName := rs.Primary.Attributes["name"]
					return projectUuid + "/" + envName, nil
				},
			},
			{ // Update and Read testing
				Config: `
				resource "coolify_project" "test" {
					name        = "TerraformAccTestEnvProject"
					description = "Test project for environment"
				}

				resource "coolify_environment" "test" {
					name         = "terraform-test-env-updated"
					project_uuid = coolify_project.test.uuid
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Name changes require replacement, not update
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionReplace),
						plancheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact("terraform-test-env-updated")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "terraform-test-env-updated"),
				),
			},
			// Delete testing is implicit and handled automatically
		},
	})
}

func TestAccEnvironmentResource_InvalidImportId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "coolify_environment" "test" {
					name         = "test"
					project_uuid = "fake-uuid"
				}
				`,
				ResourceName:  "coolify_environment.test",
				ImportState:   true,
				ImportStateId: "invalid-format",
				ExpectError:   regexp.MustCompile(`import format must be project_uuid/environment_name_or_uuid`),
			},
		},
	})
}

func TestAccEnvironmentResource_MinimalConfig(t *testing.T) {
	resName := "coolify_environment.minimal"
	
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "coolify_project" "test" {
					name = "MinimalEnvTest"
				}

				resource "coolify_environment" "minimal" {
					name         = "minimal"
					project_uuid = coolify_project.test.uuid
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "minimal"),
					resource.TestCheckResourceAttrSet(resName, "project_uuid"),
					resource.TestCheckResourceAttrSet(resName, "uuid"),
				),
			},
		},
	})
}

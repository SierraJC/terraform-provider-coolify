package service_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccServiceResource(t *testing.T) {
	resName := "coolify_service.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_service" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"

					server_uuid = "` + acctest.ServerUUID + `"
					project_uuid = "` + acctest.ProjectUUID + `"
					environment_name = "` + acctest.EnvironmentName + `"
					destination_uuid = "` + acctest.DestinationUUID + `"

					instant_deploy = false
  				compose = <<EOF
services:
  whoami:
    image: "containous/whoami"
    container_name: "simple-service"
EOF
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "server_uuid", acctest.ServerUUID),
					resource.TestCheckResourceAttr(resName, "project_uuid", acctest.ProjectUUID),
					resource.TestCheckResourceAttr(resName, "environment_name", acctest.EnvironmentName),
					resource.TestCheckResourceAttr(resName, "instant_deploy", "false"),

					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "compose"),
				),
			},
			{ // ImportState testing
				ResourceName:                         resName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ExpectError: regexp.MustCompile(
					`("instant_deploy")`,
				),
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r := s.RootModule().Resources[resName].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s/%s",
						r["server_uuid"],
						r["project_uuid"],
						r["environment_name"],
						r["uuid"],
					), nil
				},
			},
			{ // Update and Read testing
				Config: `
				resource "coolify_service" "test" {
					name        = "TerraformAccTestUpdated"
					description = "Terraform acceptance testing"

					server_uuid = "` + acctest.ServerUUID + `"
					project_uuid = "` + acctest.ProjectUUID + `"
					environment_name = "` + acctest.EnvironmentName + `"
					destination_uuid = "` + acctest.DestinationUUID + `"

					instant_deploy = false

  				compose = <<EOF
services:
  whoami2:
    image: "containous/whoami"
    container_name: "simple-service"
EOF

				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact("TerraformAccTestUpdated")),
						plancheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.StringExact("Terraform acceptance testing")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTestUpdated"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "server_uuid", acctest.ServerUUID),
				),
			},
		},
	})
}

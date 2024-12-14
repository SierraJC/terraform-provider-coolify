package service_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccMysqlDatabaseResource(t *testing.T) {
	randomName := acctest.GetRandomResourceName("mysql-db")
	resName := "coolify_mysql_database." + randomName
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: testAccMysqlDatabaseResourceConfig(randomName, "test_db", "user", "password", "root_password"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", randomName),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "server_uuid", acctest.ServerUUID),
					resource.TestCheckResourceAttr(resName, "project_uuid", acctest.ProjectUUID),
					resource.TestCheckResourceAttr(resName, "environment_name", acctest.EnvironmentName),
					resource.TestCheckResourceAttr(resName, "instant_deploy", "false"),

					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "internal_db_url"),
					resource.TestCheckResourceAttrSet(resName, "mysql_password"),
					resource.TestCheckResourceAttrSet(resName, "mysql_root_password"),
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
				Config: testAccMysqlDatabaseResourceConfig(randomName, "test_db2", "user2", "password2", "root_password2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(resName, tfjsonpath.New("internal_db_url")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "internal_db_url"),
					resource.TestCheckResourceAttr(resName, "name", randomName),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "server_uuid", acctest.ServerUUID),
					resource.TestCheckResourceAttrSet(resName, "mysql_password"),
					resource.TestCheckResourceAttrSet(resName, "mysql_root_password"),
				),
			},
		},
	})
}

func testAccMysqlDatabaseResourceConfig(name, db, user, password, rootPassword string) string {
	return fmt.Sprintf(`
		resource "coolify_mysql_database" "%[1]s" {
			name        = "%[1]s"
			description = "Terraform acceptance testing"

			server_uuid = "`+acctest.ServerUUID+`"
			project_uuid = "`+acctest.ProjectUUID+`"
			environment_name = "`+acctest.EnvironmentName+`"

			mysql_database = "%[2]s"
			mysql_user = "%[3]s"
			mysql_password = "%[4]s"
			mysql_root_password = "%[5]s"
		}
	`,
		name, db, user, password, rootPassword,
	)
}

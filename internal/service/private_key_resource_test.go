package service_test

import (
	"fmt"
	"strings"
	"testing"

	tf_acctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccPrivateKeyResource(t *testing.T) {
	randomName := acctest.GetRandomResourceName("pk")
	resName := "coolify_private_key." + randomName

	_, privateKey1, err := tf_acctest.RandSSHKeyPair(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	_, privateKey2, err := tf_acctest.RandSSHKeyPair(t.Name())
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: testAccPrivateKeyResourceConfig(randomName, privateKey1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", randomName),
					resource.TestCheckResourceAttr(resName, "private_key", privateKey1),
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{ // ImportState testing
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources[resName].Primary.Attributes["uuid"], nil
				},
			},
			{ // Update and Read testing
				Config: testAccPrivateKeyResourceConfig(randomName, privateKey2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(resName, tfjsonpath.New("fingerprint")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttr(resName, "name", randomName),
					resource.TestCheckResourceAttr(resName, "private_key", privateKey2),
				),
			},
		},
	})
}

func testAccPrivateKeyResourceConfig(name, privateKey string) string {
	return fmt.Sprintf(`
		resource "coolify_private_key" "%[1]s" {
			name        = "%[1]s"
			description = "Terraform acceptance testing"
			private_key = "%[2]s"
		}
	`,
		name,
		strings.ReplaceAll(privateKey, "\n", "\\n"),
	)
}

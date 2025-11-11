package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccApplicationResource_Public(t *testing.T) {
	resName := "coolify_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_application" "test" {
					source_type = "public"
					project_uuid = "` + acctest.ProjectUUID + `"
					server_uuid = "` + acctest.ServerUUID + `"
					environment_name = "` + acctest.EnvironmentName + `"
					git_repository = "https://github.com/coollabsio/coolify"
					git_branch = "main"
					build_pack = "nixpacks"
					ports_exposes = "80"
					name = "TerraformAccTest Public App"
					description = "Terraform acceptance testing"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "source_type", "public"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest Public App"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "git_repository", "https://github.com/coollabsio/coolify"),
					resource.TestCheckResourceAttr(resName, "git_branch", "main"),
					resource.TestCheckResourceAttr(resName, "build_pack", "nixpacks"),
					resource.TestCheckResourceAttr(resName, "ports_exposes", "80"),
					// Verify dynamic values
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{ // ImportState testing
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{ // Update and Read testing
				Config: `
				resource "coolify_application" "test" {
					source_type = "public"
					project_uuid = "` + acctest.ProjectUUID + `"
					server_uuid = "` + acctest.ServerUUID + `"
					environment_name = "` + acctest.EnvironmentName + `"
					git_repository = "https://github.com/coollabsio/coolify"
					git_branch = "main"
					build_pack = "nixpacks"
					ports_exposes = "80"
					name = "TerraformAccTest Public App Updated"
					description = "Terraform acceptance testing updated"
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact("TerraformAccTest Public App Updated")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest Public App Updated"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing updated"),
				),
			},
		},
	})
}

func TestAccApplicationResource_Dockerfile(t *testing.T) {
	resName := "coolify_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_application" "test" {
					source_type = "dockerfile"
					project_uuid = "` + acctest.ProjectUUID + `"
					server_uuid = "` + acctest.ServerUUID + `"
					environment_name = "` + acctest.EnvironmentName + `"
					dockerfile = <<EOF
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
EOF
					name = "TerraformAccTest Dockerfile App"
					description = "Terraform acceptance testing"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "source_type", "dockerfile"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest Dockerfile App"),
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
		},
	})
}

func TestAccApplicationResource_Dockerimage(t *testing.T) {
	resName := "coolify_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_application" "test" {
					source_type = "dockerimage"
					project_uuid = "` + acctest.ProjectUUID + `"
					server_uuid = "` + acctest.ServerUUID + `"
					environment_name = "` + acctest.EnvironmentName + `"
					docker_registry_image_name = "nginx"
					docker_registry_image_tag = "alpine"
					ports_exposes = "80"
					name = "TerraformAccTest Docker Image App"
					description = "Terraform acceptance testing"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "source_type", "dockerimage"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest Docker Image App"),
					resource.TestCheckResourceAttr(resName, "docker_registry_image_name", "nginx"),
					resource.TestCheckResourceAttr(resName, "docker_registry_image_tag", "alpine"),
					resource.TestCheckResourceAttr(resName, "ports_exposes", "80"),
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
		},
	})
}


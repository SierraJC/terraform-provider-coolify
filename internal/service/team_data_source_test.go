package service_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
	"terraform-provider-coolify/internal/service"
)

func TestAccTeamDataSource(t *testing.T) {
	resName := "data.coolify_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// With ID
			{
				Config: `data "coolify_team" "test" {
					id = 5
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "5"),
					resource.TestCheckResourceAttr(resName, "name", "Test"),
					resource.TestCheckResourceAttr(resName, "personal_team", "false"),
					resource.TestCheckResourceAttr(resName, "members.#", "1"),
				),
			},
			// Without ID, current API Key authenticated team
			{
				Config: `data "coolify_team" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "0"),
					resource.TestCheckResourceAttr(resName, "name", "Root Team"),
					resource.TestCheckResourceAttr(resName, "members.#", "1"),
					resource.TestCheckResourceAttr(resName, "personal_team", "true"),
				),
			},
		},
	})
}

func TestTeamDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := service.NewTeamDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test id is optional
	idAttr := resp.Schema.Attributes["id"].(schema.Int64Attribute)
	if idAttr.Required {
		t.Error("id field should not be marked as required in schema")
	}
	if !idAttr.Optional {
		t.Error("id field should be marked as optional in schema")
	}
}

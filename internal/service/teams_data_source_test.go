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

func TestAccTeamsDataSource(t *testing.T) {
	resName := "data.coolify_teams.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_teams" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "teams.#", "5"),
					// Check the first team in the list
					resource.TestCheckResourceAttr(resName, "teams.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "teams.0.name", "Root Team"),
					resource.TestCheckResourceAttrSet(resName, "teams.0.created_at"),
				),
			},
			// Single filter by name
			{
				Config: `
					data "coolify_teams" "test" {
						filter {
							name = "name"
							values = ["Root Team", "Test"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "teams.#", "3"),
					resource.TestCheckResourceAttr(resName, "teams.0.name", "Root Team"),
					resource.TestCheckResourceAttr(resName, "teams.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "teams.1.id", "1"),
					resource.TestCheckResourceAttr(resName, "teams.2.id", "2"),
					// todo: fix acceptance test server, multiple teams with the same name
				),
			},
			// Multiple filters
			{
				Config: `
					data "coolify_teams" "test" {
						filter {
							name = "id"
							values = ["0"]
						}
						filter {
							name = "personal_team"
							values = ["true"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "teams.#", "1"),
					resource.TestCheckResourceAttr(resName, "teams.0.name", "Root Team"),
					resource.TestCheckResourceAttrSet(resName, "teams.0.personal_team"),
				),
			},
			// Filter by non-string fields
			{
				Config: `
					data "coolify_teams" "test" {
						filter {
							name = "id"
							values = ["0"]
						}
						filter {
							name = "personal_team"
							values = ["true"]
						}
						filter {
							name = "name"
							values = ["Root Team"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "teams.#", "1"),
					resource.TestCheckResourceAttr(resName, "teams.0.name", "Root Team"),
					resource.TestCheckResourceAttr(resName, "teams.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "teams.0.personal_team", "true"),
				),
			},
			// Test with_members=true
			{
				Config: `
					data "coolify_teams" "test" {
						with_members = true
						filter {
							name = "id"
							values = ["0"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "teams.#", "1"),
					resource.TestCheckResourceAttr(resName, "teams.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "teams.0.name", "Root Team"),
					// Verify members array exists and has entries
					resource.TestCheckResourceAttr(resName, "teams.0.members.#", "1"),
					// Check first member has expected fields
					resource.TestCheckResourceAttrSet(resName, "teams.0.members.0.id"),
					resource.TestCheckResourceAttrSet(resName, "teams.0.members.0.name"),
					resource.TestCheckResourceAttrSet(resName, "teams.0.members.0.email"),
				),
			},
			// Test with_members=false should still work
			{
				Config: `
					data "coolify_teams" "test" {
						with_members = false
						filter {
							name = "id"
							values = ["0"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "teams.#", "1"),
					resource.TestCheckResourceAttr(resName, "teams.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "teams.0.name", "Root Team"),
					// Members array should be empty when with_members=false
					resource.TestCheckResourceAttr(resName, "teams.0.members.#", "0"),
				),
			},
		},
	})
}

func TestTeamsDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := service.NewTeamsDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test sensitive fields
	teamsAttr := resp.Schema.Attributes["teams"].(schema.SetNestedAttribute)
	sensitiveFields := []string{"discord_webhook_url", "smtp_password", "telegram_token", "resend_api_key"}
	for _, field := range sensitiveFields {
		attr := teamsAttr.NestedObject.Attributes[field].(schema.StringAttribute)
		if !attr.Sensitive {
			t.Errorf("%s field should be marked as sensitive in schema", field)
		}
	}

	// Test filter block
	_, ok := resp.Schema.Blocks["filter"].(schema.ListNestedBlock)
	if !ok {
		t.Error("filter should be a ListNestedBlock")
	}

	// Test with_members attribute exists and is optional bool
	if attr, ok := resp.Schema.Attributes["with_members"].(schema.BoolAttribute); !ok {
		t.Error("with_members should be a BoolAttribute")
	} else {
		if !attr.Optional {
			t.Error("with_members should be optional")
		}
	}
}

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

func TestAccApplicationsDataSource(t *testing.T) {
	resName := "data.coolify_applications.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_applications" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "applications.#"),
					// Check the last application in the list (expecting the first created application, order seems to be id descending)
					resource.TestCheckResourceAttrSet(resName, "applications.0.build_pack"),
					resource.TestCheckResourceAttrSet(resName, "applications.0.created_at"),
					resource.TestCheckResourceAttrSet(resName, "applications.0.name"),
					resource.TestCheckResourceAttrSet(resName, "applications.0.uuid"),
				),
			},
			// Single filter by uuid
			{
				Config: `
				data "coolify_applications" "test" {
					filter {
						name = "uuid"
						values = ["` + acctest.ApplicationUUID + `"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "applications.#", "1"),
					resource.TestCheckResourceAttr(resName, "applications.0.build_pack", "dockerfile"),
					resource.TestCheckResourceAttr(resName, "applications.0.created_at", "2024-11-10T08:59:09Z"),
					resource.TestCheckResourceAttr(resName, "applications.0.name", "dockerfile-"+acctest.ApplicationUUID),
					resource.TestCheckResourceAttr(resName, "applications.0.uuid", acctest.ApplicationUUID),
				),
			},
		},
	})
}

func TestApplicationsDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := service.NewApplicationsDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test filter block
	_, ok := resp.Schema.Blocks["filter"].(schema.ListNestedBlock)
	if !ok {
		t.Error("filter should be a ListNestedBlock")
	}
}

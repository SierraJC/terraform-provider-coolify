package service_test

import (
	"context"
	"testing"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"terraform-provider-coolify/internal/service"
)

func TestApplicationResourceSchema(t *testing.T) {
	ctx := context.Background()
	rs := service.NewApplicationResource()

	req := tfresource.SchemaRequest{}
	resp := &tfresource.SchemaResponse{}

	rs.Schema(ctx, req, resp)

	if resp.Schema.Description == "" {
		t.Error("Expected schema description to be set")
	}

	// Check that source_type is required
	if sourceTypeAttr, ok := resp.Schema.Attributes["source_type"].(schema.StringAttribute); ok {
		if !sourceTypeAttr.Required {
			t.Error("Expected source_type to be required")
		}
	} else {
		t.Error("Expected source_type attribute to exist")
	}

	// Check that project_uuid is required
	if projectUuidAttr, ok := resp.Schema.Attributes["project_uuid"].(schema.StringAttribute); ok {
		if !projectUuidAttr.Required {
			t.Error("Expected project_uuid to be required")
		}
	} else {
		t.Error("Expected project_uuid attribute to exist")
	}

	// Check that server_uuid is required
	if serverUuidAttr, ok := resp.Schema.Attributes["server_uuid"].(schema.StringAttribute); ok {
		if !serverUuidAttr.Required {
			t.Error("Expected server_uuid to be required")
		}
	} else {
		t.Error("Expected server_uuid attribute to exist")
	}

	// Check that environment_name is required
	if envNameAttr, ok := resp.Schema.Attributes["environment_name"].(schema.StringAttribute); ok {
		if !envNameAttr.Required {
			t.Error("Expected environment_name to be required")
		}
	} else {
		t.Error("Expected environment_name attribute to exist")
	}
}


package service

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type useStateForUnknownUnlessNull struct{}

func (m useStateForUnknownUnlessNull) Description(ctx context.Context) string {
	return "Handles Optional+Computed fields: marks as Unknown on create when null, preserves state on update"
}

func (m useStateForUnknownUnlessNull) MarkdownDescription(ctx context.Context) string {
	return "Handles Optional+Computed fields: marks as Unknown on create when null, preserves state on update"
}

func (m useStateForUnknownUnlessNull) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	if !req.PlanValue.IsNull() {
		return
	}

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		resp.PlanValue = types.StringUnknown()
		return
	}

	resp.PlanValue = req.StateValue
}

func UseStateForUnknownUnlessNullString() planmodifier.String {
	return useStateForUnknownUnlessNull{}
}

type useStateForUnknownUnlessNullInt64 struct{}

func (m useStateForUnknownUnlessNullInt64) Description(ctx context.Context) string {
	return "Handles Optional+Computed fields: marks as Unknown on create when null, preserves state on update"
}

func (m useStateForUnknownUnlessNullInt64) MarkdownDescription(ctx context.Context) string {
	return "Handles Optional+Computed fields: marks as Unknown on create when null, preserves state on update"
}

func (m useStateForUnknownUnlessNullInt64) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.ConfigValue.IsNull() {
		return
	}

	if !req.PlanValue.IsNull() {
		return
	}

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		resp.PlanValue = types.Int64Unknown()
		return
	}

	resp.PlanValue = req.StateValue
}

func UseStateForUnknownUnlessNullInt64() planmodifier.Int64 {
	return useStateForUnknownUnlessNullInt64{}
}

type useStateForUnknownUnlessNullBool struct{}

func (m useStateForUnknownUnlessNullBool) Description(ctx context.Context) string {
	return "Handles Optional+Computed fields: marks as Unknown on create when null, preserves state on update"
}

func (m useStateForUnknownUnlessNullBool) MarkdownDescription(ctx context.Context) string {
	return "Handles Optional+Computed fields: marks as Unknown on create when null, preserves state on update"
}

func (m useStateForUnknownUnlessNullBool) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	if !req.PlanValue.IsNull() {
		return
	}

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		resp.PlanValue = types.BoolUnknown()
		return
	}

	resp.PlanValue = req.StateValue
}

func UseStateForUnknownUnlessNullBool() planmodifier.Bool {
	return useStateForUnknownUnlessNullBool{}
}


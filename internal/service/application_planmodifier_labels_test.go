package service

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalize_empty(t *testing.T) {
	got := normalize("")
	if len(got) != 0 {
		t.Errorf("normalize(\"\") = %v, want empty slice", got)
	}
}

func TestNormalize_whitespaceOnly(t *testing.T) {
	got := normalize("  \n  \n")
	if len(got) != 0 {
		t.Errorf("normalize(whitespace) = %v, want empty slice", got)
	}
}

func TestNormalize_plaintextSorts(t *testing.T) {
	got := normalize("k2=v2\nk1=v1")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != len(want) {
		t.Fatalf("normalize len = %d, want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("normalize[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestNormalize_plaintextTrimsPerLine(t *testing.T) {
	got := normalize("  k1=v1  \n\t k2=v2\t \n")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize = %v, want %v", got, want)
	}
}

func TestNormalize_dropsEmptyLines(t *testing.T) {
	got := normalize("k1=v1\n\n\nk2=v2\n\n")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize = %v, want %v", got, want)
	}
}

func TestNormalize_mixedLineEndings(t *testing.T) {
	got := normalize("k1=v1\r\nk2=v2\r\n")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize = %v, want %v (handles CRLF)", got, want)
	}
}

func TestIsBase64Labels_validBase64WithEqLines(t *testing.T) {
	b64 := "azE9djEKazI9djI=" // base64("k1=v1\nk2=v2")
	if !isBase64Labels(b64) {
		t.Error("isBase64Labels(b64 of k1=v1\\nk2=v2) = false, want true")
	}
}

func TestIsBase64Labels_validBase64NoEqLines(t *testing.T) {
	b64 := "SGVsbG8=" // base64("Hello") — no '=' lines
	if isBase64Labels(b64) {
		t.Error("isBase64Labels(b64 of Hello) = true, want false (no '=' line)")
	}
}

func TestIsBase64Labels_corrupt(t *testing.T) {
	if isBase64Labels("!!!notb64!!!") {
		t.Error("isBase64Labels(!!!notb64!!!) = true, want false")
	}
}

func TestNormalize_base64InputDecodes(t *testing.T) {
	b64 := "azE9djEKazI9djI=" // base64("k1=v1\nk2=v2")
	got := normalize(b64)
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize(b64) = %v, want %v", got, want)
	}
}

func TestNormalize_corruptBase64FallsBackToPlaintext(t *testing.T) {
	got := normalize("!!!notb64!!!")
	want := []string{"!!!notb64!!!"}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("normalize(corrupt) = %v, want %v (plaintext fallback)", got, want)
	}
}

func TestNormalize_base64LookingButNotLabels(t *testing.T) {
	// "SGVsbG8=" = base64("Hello") — decodes OK but no '=' lines
	got := normalize("SGVsbG8=")
	want := []string{"SGVsbG8="}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("normalize(SGVsbG8=) = %v, want %v (treat as plaintext)", got, want)
	}
}

func TestFilterCertResolver_dropsLetsencrypt(t *testing.T) {
	in := []string{
		"k1=v1",
		"traefik.http.routers.foo.tls.certresolver=letsencrypt",
		"k2=v2",
	}
	got := filterCertResolver(in)
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("filterCertResolver = %v, want %v", got, want)
	}
}

func TestFilterCertResolver_preservesCustomCA(t *testing.T) {
	in := []string{
		"traefik.http.routers.foo.tls.certresolver=custom-ca",
	}
	got := filterCertResolver(in)
	if len(got) != 1 || got[0] != in[0] {
		t.Errorf("filterCertResolver dropped custom-ca: %v, want %v", got, in)
	}
}

func TestFilterCertResolver_dropsMultipleLEAcrossRouters(t *testing.T) {
	in := []string{
		"k1=v1",
		"traefik.http.routers.foo.tls.certresolver=letsencrypt",
		"traefik.http.routers.bar.tls.certresolver=letsencrypt",
		"traefik.http.routers.baz.tls.certresolver=letsencrypt",
	}
	got := filterCertResolver(in)
	if len(got) != 1 || got[0] != "k1=v1" {
		t.Errorf("filterCertResolver = %v, want [k1=v1] (3 LE lines dropped)", got)
	}
}

func TestFilterCertResolver_caseInsensitiveKeyExactValue(t *testing.T) {
	in := []string{
		"TRAEFIK.HTTP.ROUTERS.FOO.TLS.CERTRESOLVER=letsencrypt",
		"traefik.http.routers.bar.tls.certresolver=LETSENCRYPT",
	}
	got := filterCertResolver(in)
	// First line: key matches case-insensitively, value exact → drop.
	// Second line: value "LETSENCRYPT" is NOT exact-match "letsencrypt" → preserve.
	if len(got) != 1 || got[0] != in[1] {
		t.Errorf("filterCertResolver = %v, want [%q] (uppercase value preserved)", got, in[1])
	}
}

func TestNormalize_stateWithLEEqualsConfigClean(t *testing.T) {
	state := "k1=v1\ntraefik.http.routers.foo.tls.certresolver=letsencrypt"
	config := "k1=v1"
	ns := normalize(state)
	nc := normalize(config)
	if len(ns) != len(nc) || ns[0] != nc[0] {
		t.Errorf("normalize(state)=%v, normalize(config)=%v — should be equal", ns, nc)
	}
}

func TestSemanticEqual_table(t *testing.T) {
	b64 := base64StdEncodingHelper

	tests := []struct {
		name      string
		state     string
		config    string
		wantEqual bool
	}{
		{"both empty", "", "", true},
		{"whitespace only", "  \n  \n", "", true},
		{"identical plaintext", "k1=v1\nk2=v2", "k1=v1\nk2=v2", true},
		{"reverse order plaintext", "k2=v2\nk1=v1", "k1=v1\nk2=v2", true},
		{
			"state base64, config plaintext, same content",
			b64("k1=v1\nk2=v2"), "k1=v1\nk2=v2", true,
		},
		{
			"state has LE addition",
			"k1=v1\ntraefik.http.routers.foo.tls.certresolver=letsencrypt", "k1=v1", true,
		},
		{
			"state base64 + LE, config plaintext clean",
			b64("k1=v1\ntraefik.http.routers.foo.tls.certresolver=letsencrypt"), "k1=v1", true,
		},
		{"real diff added line", "k1=v1", "k1=v1\nk2=v2", false},
		{"real diff changed value", "k1=v1", "k1=v2", false},
		{
			"custom certresolver preserved",
			"k1=v1\ntraefik.http.routers.foo.tls.certresolver=custom-ca", "k1=v1", false,
		},
		{"corrupt base64 falls back", "!!!notb64!!!", "!!!notb64!!!", true},
		{"base64-looking but not labels", "SGVsbG8=", "SGVsbG8=", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := semanticEqual(tc.state, tc.config)
			if got != tc.wantEqual {
				t.Errorf("semanticEqual(%q, %q) = %v, want %v",
					tc.state, tc.config, got, tc.wantEqual)
			}
		})
	}
}

func TestCoolifyLabelsSemanticEqual_semanticallyEqualSetsPlanToState(t *testing.T) {
	mod := CoolifyLabelsSemanticEqual()
	state := types.StringValue("k1=v1\ntraefik.http.routers.foo.tls.certresolver=letsencrypt")
	config := types.StringValue("k1=v1")

	req := planmodifier.StringRequest{
		StateValue:  state,
		PlanValue:   config,
		ConfigValue: config,
	}
	resp := &planmodifier.StringResponse{
		PlanValue: config,
	}

	mod.PlanModifyString(context.Background(), req, resp)

	if resp.PlanValue.ValueString() != state.ValueString() {
		t.Errorf("PlanValue = %q, want %q (state when semantically equal)",
			resp.PlanValue.ValueString(), state.ValueString())
	}
}

func TestCoolifyLabelsSemanticEqual_realDiffLeavesPlanAlone(t *testing.T) {
	mod := CoolifyLabelsSemanticEqual()
	state := types.StringValue("k1=v1")
	config := types.StringValue("k1=v1\nk2=v2")

	req := planmodifier.StringRequest{
		StateValue:  state,
		PlanValue:   config,
		ConfigValue: config,
	}
	resp := &planmodifier.StringResponse{
		PlanValue: config,
	}

	mod.PlanModifyString(context.Background(), req, resp)

	if resp.PlanValue.ValueString() != config.ValueString() {
		t.Errorf("PlanValue = %q, want %q (config preserved when real diff)",
			resp.PlanValue.ValueString(), config.ValueString())
	}
}

func TestCoolifyLabelsSemanticEqual_nullStateEarlyReturn(t *testing.T) {
	mod := CoolifyLabelsSemanticEqual()
	state := types.StringNull()
	config := types.StringValue("k1=v1")

	req := planmodifier.StringRequest{
		StateValue:  state,
		PlanValue:   config,
		ConfigValue: config,
	}
	resp := &planmodifier.StringResponse{
		PlanValue: config,
	}

	mod.PlanModifyString(context.Background(), req, resp)

	if resp.PlanValue.ValueString() != config.ValueString() {
		t.Errorf("PlanValue = %q, want %q (untouched on null state)",
			resp.PlanValue.ValueString(), config.ValueString())
	}
}

func TestCoolifyLabelsSemanticEqual_unknownPlanEarlyReturn(t *testing.T) {
	mod := CoolifyLabelsSemanticEqual()
	state := types.StringValue("k1=v1")
	planVal := types.StringUnknown()

	req := planmodifier.StringRequest{
		StateValue:  state,
		PlanValue:   planVal,
		ConfigValue: types.StringNull(),
	}
	resp := &planmodifier.StringResponse{
		PlanValue: planVal,
	}

	mod.PlanModifyString(context.Background(), req, resp)

	if !resp.PlanValue.IsUnknown() {
		t.Errorf("PlanValue = %v, want unknown (untouched on unknown plan)",
			resp.PlanValue)
	}
}

func TestCoolifyLabelsSemanticEqual_descriptionNotEmpty(t *testing.T) {
	mod := CoolifyLabelsSemanticEqual()
	if mod.Description(context.Background()) == "" {
		t.Error("Description() = empty, want non-empty")
	}
	if mod.MarkdownDescription(context.Background()) == "" {
		t.Error("MarkdownDescription() = empty, want non-empty")
	}
}

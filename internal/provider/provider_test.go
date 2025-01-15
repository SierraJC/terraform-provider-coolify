package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"terraform-provider-coolify/internal/acctest"
	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/consts"
	"terraform-provider-coolify/internal/provider"
)

const (
	testAccPrivateKeyUUID  = "ys4g88w"
	testAccServerUUID      = "rg8ks8c"
	testAccProjectUUID     = "uoswco88w8swo40k48o8kcwk"
	testAccEnvironmentName = "production"
	testAccApplicationUUID = "mc8gw00wscww4gskgk0gwgw0"
	testAccServiceUUID     = "i0800ok00gcww840kk8sok0s"
)

func testAccPreCheck(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	variables := []string{
		consts.ENV_KEY_ENDPOINT,
		consts.ENV_KEY_TOKEN,
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func providerConfigDynamicValue(config map[string]interface{}) (tfprotov6.DynamicValue, error) {
	providerConfigTypes := map[string]tftypes.Type{
		"endpoint": tftypes.String,
		"token":    tftypes.String,
		"retry":    tftypes.Object{},
	}
	providerConfigObjectType := tftypes.Object{AttributeTypes: providerConfigTypes}

	providerConfigObjectValue := tftypes.NewValue(providerConfigObjectType, map[string]tftypes.Value{
		"endpoint": tftypes.NewValue(tftypes.String, config["endpoint"]),
		"token":    tftypes.NewValue(tftypes.String, config["token"]),
		"retry":    tftypes.NewValue(tftypes.Object{}, nil),
	})

	value, err := tfprotov6.NewDynamicValue(providerConfigObjectType, providerConfigObjectValue)
	if err != nil {
		err = fmt.Errorf("failed to create dynamic value: %w", err)
	}

	return value, err
}

func TestProtocol6ProviderServerSchemaVersion(t *testing.T) {
	t.Parallel()

	providerServer, err := acctest.TestAccProtoV6ProviderFactories["coolify"]()
	require.NotNil(t, providerServer)
	require.NoError(t, err)

	resp, err := providerServer.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	require.NotNil(t, resp.Provider)
	require.NoError(t, err)
	assert.Empty(t, resp.Diagnostics)

	assert.EqualValues(t, 0, resp.Provider.Version)
}

func TestProtocol6ProviderServerConfigure(t *testing.T) {
	testAccPreCheck(t)

	accEndpoint := os.Getenv(consts.ENV_KEY_ENDPOINT)
	accToken := os.Getenv(consts.ENV_KEY_TOKEN)

	tests := map[string]struct {
		config          map[string]interface{}
		env             map[string]string
		expectedSuccess bool
	}{
		"config: endpoint": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
			},
			expectedSuccess: false,
		},
		"config: token": {
			config: map[string]interface{}{
				"token": accToken,
			},
			expectedSuccess: false,
		},
		"config: endpoint,token": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
				"token":    accToken,
			},
			expectedSuccess: true,
		},
		"config: endpoint(invalid),token": {
			config: map[string]interface{}{
				"endpoint": "url://an invalid url %/",
				"token":    accToken,
			},
			expectedSuccess: false,
		},
		"config: endpoint,token(invalid)": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
				"token":    "invalid_token",
			},
			expectedSuccess: false,
		},
		"env: endpoint": {
			env: map[string]string{
				consts.ENV_KEY_ENDPOINT: accEndpoint,
			},
			expectedSuccess: false,
		},
		"env: endpoint,token": {
			env: map[string]string{
				consts.ENV_KEY_ENDPOINT: accEndpoint,
				consts.ENV_KEY_TOKEN:    accToken,
			},
			expectedSuccess: true,
		},
		"config: endpoint env: token": {
			config: map[string]interface{}{
				"endpoint": accEndpoint,
			},
			env: map[string]string{
				consts.ENV_KEY_TOKEN: accToken,
			},
			expectedSuccess: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv(consts.ENV_KEY_ENDPOINT, "")
			t.Setenv(consts.ENV_KEY_TOKEN, "")
			for key, value := range test.env {
				t.Setenv(key, value)
			}

			providerServer, err := acctest.TestAccProtoV6ProviderFactories["coolify"]()
			require.NotNil(t, providerServer)
			require.NoError(t, err)

			providerConfigValue, err := providerConfigDynamicValue(test.config)
			require.NotNil(t, providerConfigValue)
			require.NoError(t, err)

			resp, err := providerServer.ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{
				Config: &providerConfigValue,
			})
			require.NotNil(t, resp)
			require.NoError(t, err)

			if test.expectedSuccess {
				assert.Empty(t, resp.Diagnostics)
			} else {
				assert.NotEmpty(t, resp.Diagnostics)
			}
		})
	}
}

func TestGetRetryConfig(t *testing.T) {
	testCases := map[string]struct {
		input    *provider.RetryConfigModel
		expected api.RetryConfig
	}{
		"nil config": {
			input: nil,
			expected: api.RetryConfig{
				MaxAttempts: consts.DEFAULT_RETRY_ATTEMPTS,
				MinWait:     consts.DEFAULT_RETRY_MIN_WAIT,
				MaxWait:     consts.DEFAULT_RETRY_MAX_WAIT,
			},
		},
		"empty config": {
			input: &provider.RetryConfigModel{},
			expected: api.RetryConfig{
				MaxAttempts: consts.DEFAULT_RETRY_ATTEMPTS,
				MinWait:     consts.DEFAULT_RETRY_MIN_WAIT,
				MaxWait:     consts.DEFAULT_RETRY_MAX_WAIT,
			},
		},
		"partial config": {
			input: &provider.RetryConfigModel{
				Attempts: types.Int64Value(5),
			},
			expected: api.RetryConfig{
				MaxAttempts: 5,
				MinWait:     consts.DEFAULT_RETRY_MIN_WAIT,
				MaxWait:     consts.DEFAULT_RETRY_MAX_WAIT,
			},
		},
		"full config": {
			input: &provider.RetryConfigModel{
				Attempts: types.Int64Value(6),
				MinWait:  types.Int64Value(2),
				MaxWait:  types.Int64Value(20),
			},
			expected: api.RetryConfig{
				MaxAttempts: 6,
				MinWait:     2,
				MaxWait:     20,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := provider.GetRetryConfig(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

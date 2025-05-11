package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"terraform-provider-coolify/internal/consts"
	"terraform-provider-coolify/internal/provider"
)

const (
	PrivateKeyUUID  = "ys4g88w"
	ServerUUID      = "rg8ks8c"
	ProjectUUID     = "uoswco88w8swo40k48o8kcwk"
	EnvironmentName = "production"
	ApplicationUUID = "mc8gw00wscww4gskgk0gwgw0"
	ServiceUUID     = "joscs8wg0owkggsc4gwkw84k"
	DestinationUUID = "kgso0w8"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"coolify": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
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

// MARK: Helper functions

const testAccNamePrefix = "tf-acc"

func GetRandomResourceName(resType string) string {
	generatedIdentifier := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	return fmt.Sprintf("%s-%s-%s", testAccNamePrefix, resType, generatedIdentifier)
}

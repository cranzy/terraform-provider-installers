package installers_test

import (
	"testing"

	installers "github.com/Mirantis/terraform-provider-mirantis/mirantis/installers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := installers.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = installers.Provider()
}

package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/helpers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/testclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

// lintignore:AT001
func (td TestData) DataSourceTest(t *testing.T, steps []resource.TestStep) {
	// DataSources don't need a check destroy - however since this is a wrapper function
	// and not matching the ignore pattern `XXX_data_source_test.go`, this needs to be explicitly opted out
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}

	td.runAcceptanceTest(t, testCase)
}

func (td TestData) ResourceTest(t *testing.T, testResource types.TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := testclient.Build()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	td.runAcceptanceTest(t, testCase)
}

func RunTestsInSequence(t *testing.T, tests map[string]map[string]func(t *testing.T)) {
	for group, m := range tests {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	if testclient.EnableBinaryTesting {
		testCase.ProviderFactories = map[string]terraform.ResourceProviderFactory{
			"azurerm": func() (terraform.ResourceProvider, error) {
				azurerm := provider.TestAzureProvider()
				return azurerm, nil
			},
		}
	} else {
		testCase.Providers = testclient.SupportedProviders
	}

	resource.ParallelTest(t, testCase)
}

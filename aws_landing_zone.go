package aws_landing_zone

import (
	"github.com/tflz/aws-landing-zone/config/organization_config"
)

func Build() {
	// Load the organization configuration
	org_config, err := organization_config.LoadOrganizationConfig()
	if err != nil {
		panic(err)
	}

	// Print the home region
	println("Home Region:", org_config.HomeRegion)

	for _, ou := range org_config.OrganizationalUnits {
		switch ou.Parent.(type) {
		case organization_config.RootOrganizationalUnitReference:
			println("Organization Unit:", ou.Alias, "Name:", ou.Name, "Parent: ROOT")
		case organization_config.AliasOrganizationalUnitReference:
			println("Organization Unit:", ou.Alias, "Name:", ou.Name, "Parent Alias:", ou.Parent.(organization_config.AliasOrganizationalUnitReference).Alias)
		case organization_config.RawOrganizationalUnitReference:
			println("Organization Unit:", ou.Alias, "Name:", ou.Name, "Parent ID:", ou.Parent.(organization_config.RawOrganizationalUnitReference).ParentId)
		default:
			println("Organization Unit:", ou.Alias, "Name:", ou.Name, "Parent: unknown")
		}
	}
}

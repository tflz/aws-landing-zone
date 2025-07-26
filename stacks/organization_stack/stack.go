package organization_stack

import (
	"github.com/aws/jsii-runtime-go"
	organizations "github.com/cdktf/cdktf-provider-aws-go/aws/v21/organizationsorganization"
	organizational_units "github.com/cdktf/cdktf-provider-aws-go/aws/v21/organizationsorganizationalunit"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v21/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"github.com/tflz/aws-landing-zone/config/organization_config"
)

type OrganizationStack struct {
	Organization        organizations.OrganizationsOrganization
	OrganizationalUnits map[string]organizational_units.OrganizationsOrganizationalUnit

	terraform_stack     cdktf.TerraformStack
	organization_config *organization_config.OrganizationConfig
}

func NewStack(app cdktf.App, org_config *organization_config.OrganizationConfig) *OrganizationStack {
	stack_name := "organization"
	tf_stack := cdktf.NewTerraformStack(app, &stack_name)

	awsprovider.NewAwsProvider(tf_stack, jsii.String("management-account"), &awsprovider.AwsProviderConfig{
		Region: jsii.String(org_config.HomeRegion),
	})

	stack := &OrganizationStack{
		terraform_stack:     tf_stack,
		organization_config: org_config,

		OrganizationalUnits: make(map[string]organizational_units.OrganizationsOrganizationalUnit),
	}

	stack.createOrganization()
	stack.createOrganizationalUnits()

	return stack
}

func (o *OrganizationStack) createOrganization() {
	// Create the organization resource
	o.Organization = organizations.NewOrganizationsOrganization(o.terraform_stack, jsii.String("organization"), &organizations.OrganizationsOrganizationConfig{
		FeatureSet: jsii.String("ALL"),
	})
}

func (o *OrganizationStack) createOrganizationalUnits() {
	for _, ou := range o.GetSortedOrganizationalUnits() {
		var parent_id *string
		switch parent := ou.Parent.(type) {
		case organization_config.RootOrganizationalUnitReference:
			parent_id = jsii.String("root")
		case organization_config.RawOrganizationalUnitReference:
			parent_id = jsii.String(parent.ParentId)
		case organization_config.AliasOrganizationalUnitReference:
			parent_id = o.OrganizationalUnits[parent.Alias].Id()
		default:
			panic("Unknown parent type for organizational unit: " + ou.Alias)
		}

		o.OrganizationalUnits[ou.Alias] = organizational_units.NewOrganizationsOrganizationalUnit(o.terraform_stack, jsii.String(ou.Alias), &organizational_units.OrganizationsOrganizationalUnitConfig{
			Name:     jsii.String(ou.Name),
			ParentId: parent_id,
		})
	}
}

func (o *OrganizationStack) GetSortedOrganizationalUnits() []organization_config.OrganizationalUnit {
	graph := make(map[string][]organization_config.OrganizationalUnit)
	queue := make([]organization_config.OrganizationalUnit, 0)

	for _, ou := range o.organization_config.OrganizationalUnits {
		switch ou.Parent.(type) {
		case organization_config.RootOrganizationalUnitReference, organization_config.RawOrganizationalUnitReference:
			// Root or raw reference, no parent
			queue = append(queue, ou)
		case organization_config.AliasOrganizationalUnitReference:
			// Alias reference, add to graph
			parent_alias := ou.Parent.(organization_config.AliasOrganizationalUnitReference).Alias
			graph[parent_alias] = append(graph[parent_alias], ou)
		}
	}

	var sorted []organization_config.OrganizationalUnit

	for len(queue) > 0 {
		current := queue[0]

		// Remove current from the queue and append children to the queue
		queue = append(queue[1:], graph[current.Alias]...)

		// Add current to the sorted list
		sorted = append(sorted, current)
	}

	return sorted
}

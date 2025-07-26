package aws_landing_zone

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"github.com/tflz/aws-landing-zone/config/organization_config"
	"github.com/tflz/aws-landing-zone/stacks/organization_stack"
)

func Build() {
	// Load the organization configuration
	org_config, err := organization_config.LoadOrganizationConfig()
	if err != nil {
		panic(err)
	}

	app := cdktf.NewApp(nil)
	organization_stack.NewStack(app, org_config)

	app.Synth()
}

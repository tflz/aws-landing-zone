package aws_landing_zone

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"github.com/tflz/aws-landing-zone/config_store"
	"github.com/tflz/aws-landing-zone/stacks/organization_stack"
)

func Build() {
	// Load the organization configuration
	config_store, err := config_store.NewConfigStore()
	if err != nil {
		panic(err)
	}

	app := cdktf.NewApp(nil)
	organization_stack.NewStack(app, config_store)

	app.Synth()
}

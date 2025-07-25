package organization_config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type rawOrganizationConfig struct {
	HomeRegion          string                `hcl:"home_region"`
	OrganizationalUnits []rawOrganizationUnit `hcl:"organizational_unit,block"`

	ctx hcl.EvalContext
}

type rawOrganizationUnit struct {
	Alias  string         `hcl:"alias,label"`
	Name   string         `hcl:"name"`
	Parent hcl.Expression `hcl:"parent"`
}

func parse(filename string) (*rawOrganizationConfig, error) {
	parser := hclparse.NewParser()
	configFile, err := parser.ParseHCLFile(filename)
	if err != nil {
		return nil, err
	}

	var config rawOrganizationConfig
	diag := gohcl.DecodeBody(configFile.Body, nil, &config)
	if diag.HasErrors() {
		return nil, diag
	}

	return &config, nil
}

func (c *rawOrganizationConfig) generateContext() hcl.EvalContext {
	var variables = make(map[string]cty.Value)

	variables["root"] = cty.ObjectVal(map[string]cty.Value{
		"__ref": cty.StringVal("root"),
	})

	organizationalUnitMap := make(map[string]cty.Value)
	for _, ou := range c.OrganizationalUnits {
		organizationalUnitMap[ou.Alias] = cty.ObjectVal(map[string]cty.Value{
			"__ref": cty.StringVal("organizational_unit"),
			"alias": cty.StringVal(ou.Alias),
		})
	}

	variables["organizational_unit"] = cty.ObjectVal(organizationalUnitMap)

	return hcl.EvalContext{
		Variables: variables,
	}
}

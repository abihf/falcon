package parser

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type EnumDirectiveContext struct {
	Name   string
	Config *graphql.EnumConfig
	Ast    *ast.EnumDefinition
}

type EnumValueDirectiveContext struct {
	Name   string
	Config *graphql.EnumValueConfig
	Ast    *ast.EnumValueDefinition
}

func (p *parser) parseEnum(def *ast.EnumDefinition) (*graphql.Enum, error) {
	var enumDesc string

	values := make(graphql.EnumValueConfigMap)
	for _, valDef := range def.Values {
		valName := valDef.Name.Value
		var valDesc string
		if valDef.Description != nil {
			valDesc = valDef.Description.Value
		}
		values[valName] = &graphql.EnumValueConfig{
			Value:       valName,
			Description: valDesc,
		}

		if len(valDef.Directives) > 0 {
			dirContext := &EnumValueDirectiveContext{
				Name:   valName,
				Config: values[valName],
				Ast:    valDef,
			}
			if err := p.processDirectives(valDef.Directives, valDef.Kind, dirContext); err != nil {
				return nil, err
			}
		}
	}

	config := graphql.EnumConfig{
		Name:        def.Name.Value,
		Description: enumDesc,
		Values:      values,
	}

	if len(def.Directives) > 0 {
		dirContext := &EnumDirectiveContext{
			Name:   def.Name.Value,
			Config: &config,
			Ast:    def,
		}
		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}

	return graphql.NewEnum(config), nil
}

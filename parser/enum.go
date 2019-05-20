package parser

import (
	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func (p *parser) parseEnum(def *ast.EnumDefinition) (*graphql.Enum, error) {
	values := make(graphql.EnumValueConfigMap)
	for _, valDef := range def.Values {
		valName := valDef.Name.Value
		values[valName] = &graphql.EnumValueConfig{
			Value:       valName,
			Description: stringValue(valDef.Description),
		}

		if len(valDef.Directives) > 0 {
			dirContext := &directives.EnumValueDirectiveContext{
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
		Description: stringValue(def.Description),
		Values:      values,
	}

	if len(def.Directives) > 0 {
		dirContext := &directives.EnumDirectiveContext{
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

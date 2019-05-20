package parser

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type ScalarDirectiveContext struct {
	Name   string
	Config *graphql.ScalarConfig
	Ast    *ast.ScalarDefinition
}

func (p *parser) parseScalar(def *ast.ScalarDefinition) (*graphql.Scalar, error) {
	name := def.Name.Value

	serializeFn, ok := p.resolver[name+".serialize"]
	if !ok {
		return nil, fmt.Errorf("No Serializer for scalar %s", name)
	}
	parseValueFn, ok := p.resolver[name+".parseValue"]
	if !ok {
		return nil, fmt.Errorf("No value parser scalar %s", name)
	}
	parseLiteralFn, ok := p.resolver[name+".parseLiteral"]
	if !ok {
		return nil, fmt.Errorf("No literal parser scalar %s", name)
	}

	config := graphql.ScalarConfig{
		Name:         def.Name.Value,
		Description:  stringValue(def.Description),
		Serialize:    serializeFn.(graphql.SerializeFn),
		ParseValue:   parseValueFn.(graphql.ParseValueFn),
		ParseLiteral: parseLiteralFn.(graphql.ParseLiteralFn),
	}
	if len(def.Directives) > 0 {
		dirContext := &ScalarDirectiveContext{
			Name:   def.Name.Value,
			Config: &config,
			Ast:    def,
		}
		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}

	return graphql.NewScalar(config), nil
}

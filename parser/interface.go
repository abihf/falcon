package parser

import (
	"fmt"

	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func (p *parser) parseInterface(def *ast.InterfaceDefinition) (*graphql.Interface, error) {
	resolverName := def.Name.Value

	var resolver graphql.ResolveTypeFn
	if r, ok := p.resolver[resolverName]; ok {
		resolver, ok = r.(graphql.ResolveTypeFn)
		if !ok {
			return nil, fmt.Errorf("Invalid resolver for %s, should be `graphql.ResolveTypeFn`", resolverName)
		}
	}

	fields := make(graphql.Fields)
	config := graphql.InterfaceConfig{
		Name:        def.Name.Value,
		Description: stringValue(def.Description),
		Fields:      fields,
		ResolveType: resolver,
	}

	if len(def.Directives) > 0 {
		dirContext := &directives.InterfaceDirectiveContext{
			Name:   def.Name.Value,
			Config: &config,
			Ast:    def,
		}

		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}

	parsed := graphql.NewInterface(config)
	p.types[def.Name.Value] = parsed

	for _, fieldDef := range def.Fields {
		var err error
		fields[fieldDef.Name.Value], err = p.parseField(def.Name.Value, fieldDef)
		if err != nil {
			return nil, err
		}
	}
	return parsed, nil

}

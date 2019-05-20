package parser

import (
	"fmt"

	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func (p *parser) parseDirective(def *ast.DirectiveDefinition) (*graphql.Directive, error) {
	var locations []string
	args := make(graphql.FieldConfigArgument)

	for _, loc := range def.Locations {
		locations = append(locations, loc.Value)
	}

	for _, argDef := range def.Arguments {
		argType, err := p.convertType(argDef.Type)
		if err != nil {
			return nil, err
		}
		args[argDef.Name.Value] = &graphql.ArgumentConfig{
			Type:         argType,
			Description:  stringValue(argDef.Description),
			DefaultValue: astValue(argDef.DefaultValue),
		}
	}

	config := graphql.DirectiveConfig{
		Name:        def.Name.Value,
		Description: stringValue(def.Description),
		Locations:   locations,
		Args:        args,
	}
	return graphql.NewDirective(config), nil
}

func (p *parser) processDirectives(directives []*ast.Directive, kind string, context interface{}) error {
	for _, dir := range directives {
		if err := p.processDirective(dir, kind, context); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) processDirective(dir *ast.Directive, kind string, context interface{}) error {
	resolverName := "@" + dir.Name.Value
	resolver, ok := p.resolver[resolverName]
	if !ok {
		// don't do anything
		return nil
	}

	fn, ok := resolver.(directives.Visitor)
	if !ok {
		fn, ok = resolver.(func(param *directives.VisitorParam) error)
	}
	if !ok {
		return fmt.Errorf("Invalid directive resolver function for %s", resolverName)
	}

	args := make(map[string]interface{})
	for _, arg := range dir.Arguments {
		args[arg.Name.Value] = arg.Value.GetValue()
	}

	param := &directives.VisitorParam{
		Name:    dir.Name.Value,
		Kind:    kind,
		Context: context,
		Args:    args,
	}
	return fn(param)
}

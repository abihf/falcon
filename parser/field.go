package parser

import (
	"fmt"
	"reflect"

	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func (p *parser) parseField(parentName string, def *ast.FieldDefinition) (*graphql.Field, error) {
	fieldName := def.Name.Value
	resolverName := parentName + "." + fieldName
	var resolver graphql.FieldResolveFn = graphql.DefaultResolveFn
	if r, ok := p.resolver[resolverName]; ok {
		resolver, ok = r.(graphql.FieldResolveFn)
		if !ok {
			resolver, ok = r.(func(graphql.ResolveParams) (interface{}, error))
		}
		if !ok {
			return nil, fmt.Errorf(
				"Invalid resolver for %s, expected `func(graphql.ResolveParams)(interface{}, error)` got `%s`",
				resolverName,
				reflect.TypeOf(r).String(),
			)
		}
	}

	fieldType, err := p.convertType(def.Type)
	if err != nil {
		return nil, err
	}

	args := make(graphql.FieldConfigArgument)
	for _, argDef := range def.Arguments {
		argName := argDef.Name.Value
		argType, err := p.convertType(argDef.Type)
		if err != nil {
			return nil, err
		}

		args[argName] = &graphql.ArgumentConfig{
			Description:  stringValue(argDef.Description),
			DefaultValue: astValue(argDef.DefaultValue),
			Type:         argType,
		}

		if len(argDef.Directives) > 0 {
			dirContext := &directives.FieldArgDirectiveContext{
				Name:   argName,
				Field:  fieldName,
				Parent: parentName,
				Config: args[argName],
				Ast:    argDef,
			}
			if err := p.processDirectives(argDef.Directives, argDef.Kind, dirContext); err != nil {
				return nil, err
			}
		}
	}
	field := &graphql.Field{
		Name:        fieldName,
		Type:        fieldType,
		Resolve:     resolver,
		Args:        args,
		Description: stringValue(def.Description),
	}

	if len(def.Directives) > 0 {
		dirContext := &directives.FieldDirectiveContext{
			Name:   fieldName,
			Parent: parentName,
			Config: field,
			Ast:    def,
		}
		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}

	return field, nil
}

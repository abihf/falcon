package falcon

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type FieldDirectiveContext struct {
	Name   string
	Parent string
	Config *graphql.Field
	Ast    *ast.FieldDefinition
}

type FieldArgDirectiveContext struct {
	Name   string
	Field  string
	Parent string
	Config *graphql.ArgumentConfig
	Ast    *ast.InputValueDefinition
}

func (p *parser) parseField(parentName string, def *ast.FieldDefinition) (*graphql.Field, error) {
	fieldName := def.Name.Value
	resolverName := parentName + "." + fieldName
	var resolver graphql.FieldResolveFn
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

	var fieldDesc string
	if def.Description != nil {
		fieldDesc = def.Description.Value
	}

	args := make(graphql.FieldConfigArgument)
	for _, argDef := range def.Arguments {
		argName := argDef.Name.Value
		argType, err := p.convertType(argDef.Type)
		if err != nil {
			return nil, err
		}

		var argDesc string
		if argDef.Description != nil {
			argDesc = argDef.Description.Value
		}
		var defaultValue interface{}
		if argDef.DefaultValue != nil {
			defaultValue = argDef.DefaultValue.GetValue()
		}

		args[argName] = &graphql.ArgumentConfig{
			Description:  argDesc,
			DefaultValue: defaultValue,
			Type:         argType,
		}

		if len(argDef.Directives) > 0 {
			dirContext := &FieldArgDirectiveContext{
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
		Description: fieldDesc,
	}

	if len(def.Directives) > 0 {
		dirContext := &FieldDirectiveContext{
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

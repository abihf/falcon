package parser

import (
	"fmt"
	"reflect"

	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func (p *parser) parseObject(def *ast.ObjectDefinition) (*graphql.Object, error) {
	var typeOfResolver graphql.IsTypeOfFn
	typeOfResolverName := def.Name.Value + ".isTypeOf"
	if r, ok := p.resolver[typeOfResolverName]; ok {
		typeOfResolver, ok = r.(graphql.IsTypeOfFn)
		if !ok {
			typeOfResolver, ok = r.(func(graphql.IsTypeOfParams) bool)
		}
		if !ok {
			return nil, fmt.Errorf(
				"Invalid resolver for %s. got %s",
				typeOfResolverName,
				reflect.TypeOf(r).String(),
			)
		}
	}

	var interfaces []*graphql.Interface
	for _, ifaceName := range def.Interfaces {
		iface, err := p.getParsed(ifaceName.Name.Value, kinds.InterfaceDefinition)
		if err != nil {
			return nil, fmt.Errorf("Invalid interface %s", err.Error())
		}
		interfaces = append(interfaces, iface.(*graphql.Interface))
	}

	fields := make(graphql.Fields)
	config := graphql.ObjectConfig{
		Name:        def.Name.Value,
		Description: stringValue(def.Description),
		Fields:      fields,
		Interfaces:  interfaces,
		IsTypeOf:    typeOfResolver,
	}

	if len(def.Directives) > 0 {
		dirContext := &directives.ObjectDirectiveContext{
			Name:   def.Name.Value,
			Config: &config,
			Ast:    def,
		}
		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}

	// TODO: store parsed types first
	obj := graphql.NewObject(config)
	p.types[def.Name.Value] = obj

	for _, fieldDef := range def.Fields {
		var err error
		fields[fieldDef.Name.Value], err = p.parseField(def.Name.Value, fieldDef)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

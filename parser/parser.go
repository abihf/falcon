package parser

import (
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"

	gqlParser "github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type parser struct {
	resolver    map[string]interface{}
	definitions map[string]ast.Definition
	types       map[string]graphql.Type
}

func Create(resolver map[string]interface{}) *parser {
	// insert common resolver
	if _, ok := resolver["@deprecated"]; !ok {
		// resolver["@deprecated"] = CreateDeprecationCallback()
	}

	// populate definitions map
	definitions := make(map[string]ast.Definition)
	return &parser{
		resolver:    resolver,
		definitions: definitions,
		types:       make(map[string]graphql.Type),
	}
}

func (p *parser) ParseSource(name string, body []byte) error {
	doc, err := gqlParser.Parse(gqlParser.ParseParams{
		Source: source.NewSource(&source.Source{Name: name, Body: body}),
		Options: gqlParser.ParseOptions{
			NoLocation: false,
			NoSource:   false,
		},
	})

	if err != nil {
		return err
	}

	for _, d := range doc.Definitions {
		switch d.(type) {
		case *ast.SchemaDefinition:
			def := d.(*ast.SchemaDefinition)
			p.definitions["_schema"] = def

		case *ast.ObjectDefinition:
			def := d.(*ast.ObjectDefinition)
			p.definitions[def.Name.Value] = def

		case *ast.InterfaceDefinition:
			def := d.(*ast.InterfaceDefinition)
			p.definitions[def.Name.Value] = def

		case *ast.EnumDefinition:
			def := d.(*ast.EnumDefinition)
			p.definitions[def.Name.Value] = def

		default:
			return graphql.NewLocatedError(
				fmt.Sprintf("Unknown definition type %s", d.GetKind()),
				[]ast.Node{d},
			)
		}
	}

	return nil
}

func (p *parser) getParsed(name string, expected string) (graphql.Type, error) {
	if parsed, ok := p.types[name]; ok {
		return parsed, nil
	}

	def, ok := p.definitions[name]
	if !ok {
		return nil, fmt.Errorf("Can not found %s", name)
	}
	if expected != "" && expected != def.GetKind() {
		return nil, graphql.NewLocatedError(
			fmt.Sprintf("Invalid type for %s: expected %s got %s", name, expected, def.GetKind()),
			[]ast.Node{def},
		)
	}

	parsed, err := p.parseDefinition(def)
	if err != nil {
		return nil, err
	}
	p.types[name] = parsed
	return parsed, nil
}

func (p *parser) parseDefinition(d ast.Definition) (graphql.Type, error) {
	switch d.(type) {
	case *ast.ObjectDefinition:
		return p.parseObject(d.(*ast.ObjectDefinition))
	case *ast.InterfaceDefinition:
		return p.parseInterface(d.(*ast.InterfaceDefinition))
	case *ast.EnumDefinition:
		return p.parseEnum(d.(*ast.EnumDefinition))
	}
	return nil, graphql.NewLocatedError(
		fmt.Sprintf("Can not parse definition of type %s"),
		[]ast.Node{d},
	)
}

func (p *parser) convertType(t ast.Type) (graphql.Output, error) {
	switch t.(type) {
	case *ast.NonNull:
		o, err := p.convertType(t.(*ast.NonNull).Type)
		if err != nil {
			return nil, err
		}
		return graphql.NewNonNull(o), nil

	case *ast.List:
		o, err := p.convertType(t.(*ast.List).Type)
		if err != nil {
			return nil, err
		}
		return graphql.NewList(o), nil

	case *ast.Named:
		name := t.(*ast.Named).Name.Value
		switch name {
		case "String":
			return graphql.String, nil
		case "ID":
			return graphql.ID, nil
		case "Int":
			return graphql.Int, nil
		case "Float":
			return graphql.Float, nil
		case "Boolean":
			return graphql.Boolean, nil
		default:
			parsed, err := p.getParsed(name, "")
			if err != nil {
				return nil, gqlerrors.NewSyntaxError(t.GetLoc().Source, t.GetLoc().Start, err.Error())
			}
			return parsed, nil
		}
	}
	return nil, graphql.NewLocatedError(
		fmt.Sprintf("Can not parse type %s", t.String()),
		[]ast.Node{t},
	)
}

package falcon

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type parser struct {
	resolver    Resolver
	definitions map[string]ast.Definition
	types       map[string]graphql.Type
}

func createParser(doc *ast.Document, resolver Resolver) *parser {
	// insert common resolver
	if _, ok := resolver["@deprecated"]; !ok {
		resolver["@deprecated"] = CreateDeprecationCallback()
	}

	// populate definitions map
	definitions := make(map[string]ast.Definition)
	for _, d := range doc.Definitions {
		switch d.(type) {
		case *ast.SchemaDefinition:
			def := d.(*ast.SchemaDefinition)
			definitions["_schema"] = def

		case *ast.ObjectDefinition:
			def := d.(*ast.ObjectDefinition)
			definitions[def.Name.Value] = def

		case *ast.InterfaceDefinition:
			def := d.(*ast.InterfaceDefinition)
			definitions[def.Name.Value] = def

		case *ast.EnumDefinition:
			def := d.(*ast.EnumDefinition)
			definitions[def.Name.Value] = def
		}
	}

	return &parser{
		resolver:    resolver,
		definitions: definitions,
		types:       make(map[string]graphql.Type),
	}
}

func (p *parser) getParsed(name string, expected string) (graphql.Type, error) {
	if parsed, ok := p.types[name]; ok {
		return parsed, nil
	}

	def, ok := p.definitions[name]
	if !ok {
		return nil, fmt.Errorf("%s not found", name)
	}
	if expected != "" && expected != def.GetKind() {
		return nil, fmt.Errorf("Invalid type for %s, expected %s got %s. {%s}", name, expected, def.GetKind(), locToString(def.GetLoc()))
	}

	parsed, err := p.parse(def)
	if err != nil {
		return nil, err
	}
	p.types[name] = parsed
	return parsed, nil
}

func (p *parser) parse(d ast.Definition) (graphql.Type, error) {
	switch d.(type) {
	case *ast.ObjectDefinition:
		return p.parseObject(d.(*ast.ObjectDefinition))
	case *ast.InterfaceDefinition:
		return p.parseInterface(d.(*ast.InterfaceDefinition))
	case *ast.EnumDefinition:
		return p.parseEnum(d.(*ast.EnumDefinition))
	}
	return nil, fmt.Errorf("Can not parse definition of type %s\nat %s", d.GetKind(), locToString(d.GetLoc()))
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
			return p.getParsed(name, "")
		}
	}
	return nil, fmt.Errorf("Can not parse type %s\nat %s", t.String(), locToString(t.GetLoc()))
}

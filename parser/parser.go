package parser

import (
	"fmt"

	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql/gqlerrors"

	gqlParser "github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type namedTypeDefinition interface {
	GetName() *ast.Name
	GetKind() string
}

type parser struct {
	resolver    map[string]interface{}
	definitions map[string]ast.Node
	types       map[string]graphql.Type
	typeNames   []string
	directives  []*graphql.Directive
}

func Create(resolver map[string]interface{}) *parser {
	// insert common resolver
	if _, ok := resolver["@deprecated"]; !ok {
		resolver["@deprecated"] = directives.CreateDeprecationCallback()
	}

	// populate definitions map
	definitions := make(map[string]ast.Node)
	return &parser{
		resolver:    resolver,
		definitions: definitions,
		types:       make(map[string]graphql.Type),
		directives:  graphql.SpecifiedDirectives,
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
		switch def := d.(type) {
		case *ast.SchemaDefinition:
			p.definitions["_schema"] = def

		case *ast.ObjectDefinition:
			p.definitions[def.Name.Value] = d
			p.typeNames = append(p.typeNames, def.GetName().Value)

		case *ast.DirectiveDefinition:
			directive, err := p.parseDirective(def)
			if err != nil {
				return err
			}
			p.directives = append(p.directives, directive)

		case namedTypeDefinition:
			p.definitions[def.GetName().Value] = d

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

func (p *parser) parseDefinition(node ast.Node) (graphql.Type, error) {
	switch def := node.(type) {
	case *ast.ObjectDefinition:
		return p.parseObject(def)
	case *ast.InterfaceDefinition:
		return p.parseInterface(def)
	case *ast.EnumDefinition:
		return p.parseEnum(def)
	case *ast.InputObjectDefinition:
		return p.parseInputObject(def)
	case *ast.ScalarDefinition:
		return p.parseScalar(def)
	default:
		return nil, graphql.NewLocatedError(
			fmt.Sprintf("Can not parse definition of type %s", node.GetKind()),
			[]ast.Node{node},
		)
	}
}

func (p *parser) convertType(astType ast.Type) (graphql.Output, error) {
	switch t := astType.(type) {
	case *ast.NonNull:
		o, err := p.convertType(t.Type)
		if err != nil {
			return nil, err
		}
		return graphql.NewNonNull(o), nil

	case *ast.List:
		o, err := p.convertType(t.Type)
		if err != nil {
			return nil, err
		}
		return graphql.NewList(o), nil

	case *ast.Named:
		name := t.Name.Value
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
		fmt.Sprintf("Can not parse type %s", astType.String()),
		[]ast.Node{astType},
	)
}

package parser

import (
	"github.com/abihf/falcon-graphql/directives"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func (p *parser) parseInputObject(def *ast.InputObjectDefinition) (*graphql.InputObject, error) {

	fields := make(graphql.InputObjectFieldMap)
	for _, fieldDef := range def.Fields {
		inputValue, err := p.parseInputValue(fieldDef)
		if err != nil {
			return nil, err
		}
		fields[fieldDef.Name.Value] = inputValue
	}

	config := graphql.InputObjectConfig{
		Name:        def.Name.Value,
		Description: stringValue(def.Description),
		Fields:      fields,
	}

	if len(def.Directives) > 0 {
		dirContext := &directives.InputObjectDirectiveContext{
			Name:   def.Name.Value,
			Config: &config,
			Ast:    def,
		}
		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}

	return graphql.NewInputObject(config), nil
}

func (p *parser) parseInputValue(def *ast.InputValueDefinition) (*graphql.InputObjectField, error) {
	inputType, err := p.convertType(def.Type)
	if err != nil {
		return nil, err
	}

	field := &graphql.InputObjectField{
		PrivateDescription: stringValue(def.Description),
		DefaultValue:       astValue(def.DefaultValue),
		Type:               inputType,
	}
	if len(def.Directives) > 0 {
		dirContext := &directives.InputValueDirectiveContext{
			Name:  def.Name.Value,
			Field: field,
			Ast:   def,
		}
		if err := p.processDirectives(def.Directives, def.Kind, dirContext); err != nil {
			return nil, err
		}
	}
	return field, nil
}

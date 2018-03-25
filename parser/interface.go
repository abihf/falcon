package parser

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type InterfaceDirectiveContext struct {
	Name   string
	Config *graphql.InterfaceConfig
	Ast    *ast.InterfaceDefinition
}

func (p *parser) parseInterface(def *ast.InterfaceDefinition) (*graphql.Interface, error) {
	resolverName := def.Name.Value

	var resolver graphql.ResolveTypeFn
	if r, ok := p.resolver[resolverName]; ok {
		resolver, ok = r.(graphql.ResolveTypeFn)
		if !ok {
			return nil, fmt.Errorf("Invalid resolver for %s, should be `graphql.ResolveTypeFn`", resolverName)
		}
	}

	var description string
	if def.Description != nil {
		description = def.Description.Value
	}

	fields := make(graphql.Fields)
	config := graphql.InterfaceConfig{
		Name:        def.Name.Value,
		Description: description,
		Fields:      fields,
		ResolveType: resolver,
	}

	if len(def.Directives) > 0 {
		dirContext := &InterfaceDirectiveContext{
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

// func defaultResolveType(p *parser, param graphql.ResolveTypeParams) *graphql.Object {
// 	var typeName string
// 	if val, ok := param.Value.(map[string]interface{}); ok {
// 		tn, ok := val["__typename"]
// 		if !ok {
// 			return nil
// 		}
// 		typeName = tn.(string)
// 	} else {
// 		sourceVal := reflect.ValueOf(param.Value)
// 		if sourceVal.Kind() == reflect.Ptr {
// 			sourceVal = sourceVal.Elem()
// 		}
// 		if !sourceVal.IsValid() {
// 			return nil
// 		}

// 		t := sourceVal.Type()
// 		if t.Kind() == reflect.Struct {
// 			fieldCount := t.NumField()
// 			fieldIndex := -1
// 			for index := 0; fieldIndex < 0 && index < fieldCount; index++ {
// 				field := t.Field(index)
// 				if field.Type.Kind() == reflect.String && strings.EqualFold(field.Name, "TypeName") {
// 					fieldIndex = index
// 				}
// 				// field.Tag.Get
// 			}
// 			if fieldIndex >= 0 {
// 				typeName = sourceVal.Field(fieldIndex).Interface().(string)
// 			}
// 		}
// 	}

// 	if typeName != "" {
// 		obj, err := p.getParsed(typeName, kinds.ObjectDefinition)
// 		if err != nil {
// 			// just return nil
// 			return nil
// 		}

// 		return obj.(*graphql.Object)
// 	}
// 	return nil
// }

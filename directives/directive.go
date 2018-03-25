package directives

import (
	"github.com/abihf/falcon-graphql/parser"
	"github.com/graphql-go/graphql/language/kinds"
)

type Visitor interface {
	Object(*parser.ObjectDirectiveContext, map[string]interface{}) error
	Field(*parser.FieldDirectiveContext, map[string]interface{}) error
	FieldArg(*parser.FieldArgDirectiveContext, map[string]interface{}) error
	Interface(*parser.InterfaceDirectiveContext, map[string]interface{}) error
	Enum(*parser.EnumDirectiveContext, map[string]interface{}) error
	EnumValue(*parser.EnumValueDirectiveContext, map[string]interface{}) error
}

func CreateCallback(visitor Visitor) parser.DirectiveCallback {
	return func(param *parser.DirectiveCallbackParam) error {
		switch param.Kind {
		case kinds.ObjectDefinition:
			visitor.Object(param.Context.(*parser.ObjectDirectiveContext), param.Args)
		case kinds.FieldDefinition:
			visitor.Field(param.Context.(*parser.FieldDirectiveContext), param.Args)
		case kinds.InputValueDefinition:
			visitor.FieldArg(param.Context.(*parser.FieldArgDirectiveContext), param.Args)
		case kinds.InterfaceDefinition:
			visitor.Interface(param.Context.(*parser.InterfaceDirectiveContext), param.Args)
		case kinds.EnumDefinition:
			visitor.Enum(param.Context.(*parser.EnumDirectiveContext), param.Args)
		case kinds.EnumValueDefinition:
			visitor.EnumValue(param.Context.(*parser.EnumValueDirectiveContext), param.Args)
		}
		return nil
	}
}

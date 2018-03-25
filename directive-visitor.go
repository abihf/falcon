package falcon

import (
	"github.com/graphql-go/graphql/language/kinds"
)

type Visitor interface {
	Enum(*EnumDirectiveContext, map[string]interface{}) error
	EnumValue(*EnumValueDirectiveContext, map[string]interface{}) error
}

func CreateCallback(visitor Visitor) DirectiveCallback {
	return func(param *DirectiveCallbackParam) error {
		switch param.Kind {
		case kinds.EnumDefinition:
			visitor.Enum(param.Context.(*EnumDirectiveContext), param.Args)
		case kinds.EnumValueDefinition:
			visitor.EnumValue(param.Context.(*EnumValueDirectiveContext), param.Args)
		}
		return nil
	}
}

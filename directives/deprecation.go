package directives

import (
	"github.com/graphql-go/graphql"
)

type deprecation struct{ DefaultVisitor }

func CreateDeprecationCallback() Visitor {
	return CreateCallback(&deprecation{})
}

func (d *deprecation) Field(dirContext *FieldDirectiveContext, args map[string]interface{}) error {
	dirContext.Config.DeprecationReason = d.getDeprecationReason(args)
	return nil
}

func (d *deprecation) EnumValue(dirContext *EnumValueDirectiveContext, args map[string]interface{}) error {
	dirContext.Config.DeprecationReason = d.getDeprecationReason(args)
	return nil
}

func (d *deprecation) getDeprecationReason(args map[string]interface{}) string {
	reason, ok := args["reason"]
	if !ok {
		return graphql.DefaultDeprecationReason
	}

	str, ok := reason.(string)
	if !ok {
		return graphql.DefaultDeprecationReason // should throw error
	}

	return str
}

package directives

import (
	"github.com/abihf/falcon-graphql/parser"
	"github.com/graphql-go/graphql"
)

type deprecation struct{}

func CreateDeprecationCallback() parser.DirectiveCallback {
	return CreateCallback(&deprecation{})
}

func (d *deprecation) Object(dirContext *parser.ObjectDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (d *deprecation) Field(dirContext *parser.FieldDirectiveContext, args map[string]interface{}) error {
	dirContext.Config.DeprecationReason = d.getDeprecationReason(args)
	return nil
}

func (d *deprecation) FieldArg(dirContext *parser.FieldArgDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (d *deprecation) Interface(dirContext *parser.InterfaceDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (d *deprecation) Enum(dirContext *parser.EnumDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (d *deprecation) EnumValue(dirContext *parser.EnumValueDirectiveContext, args map[string]interface{}) error {
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

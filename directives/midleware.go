package directives

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/kinds"
)

type midleware struct {
	fn MidlewareFn
}

type MidlewareFn func(param graphql.ResolveParams, args map[string]interface{}, next graphql.FieldResolveFn) (interface{}, error)

func CreateMidleware(fn MidlewareFn) Visitor {
	return func(param *VisitorParam) error {
		if param.Kind == kinds.FieldDefinition {
			dirContext := param.Context.(*FieldDirectiveContext)
			next := dirContext.Config.Resolve
			args := param.Args
			dirContext.Config.Resolve = func(resolveParam graphql.ResolveParams) (interface{}, error) {
				return fn(resolveParam, args, next)
			}
		}
		return nil
	}
}

package parser

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
)

type DirectiveCallbackParam struct {
	Name    string
	Kind    string
	Context interface{}
	Args    map[string]interface{}
}
type DirectiveCallback func(param *DirectiveCallbackParam) error

func (p *parser) processDirectives(directives []*ast.Directive, kind string, context interface{}) error {
	for _, dir := range directives {
		if err := p.processDirective(dir, kind, context); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) processDirective(dir *ast.Directive, kind string, context interface{}) error {
	resolverName := "@" + dir.Name.Value
	resolver, ok := p.resolver[resolverName]
	if !ok {
		// don't do anything
		return nil
	}

	fn, ok := resolver.(DirectiveCallback)
	if !ok {
		fn, ok = resolver.(func(param *DirectiveCallbackParam) error)
	}
	if !ok {
		return fmt.Errorf("Invalid directive resolver function for %s", resolverName)
	}

	args := make(map[string]interface{})
	for _, arg := range dir.Arguments {
		args[arg.Name.Value] = arg.Value.GetValue()
	}

	param := &DirectiveCallbackParam{
		Name:    dir.Name.Value,
		Kind:    kind,
		Context: context,
		Args:    args,
	}
	return fn(param)
}

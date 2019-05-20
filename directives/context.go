package directives

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type EnumDirectiveContext struct {
	Name   string
	Config *graphql.EnumConfig
	Ast    *ast.EnumDefinition
}

type EnumValueDirectiveContext struct {
	Name   string
	Config *graphql.EnumValueConfig
	Ast    *ast.EnumValueDefinition
}

type FieldDirectiveContext struct {
	Name   string
	Parent string
	Config *graphql.Field
	Ast    *ast.FieldDefinition
}

type FieldArgDirectiveContext struct {
	Name   string
	Field  string
	Parent string
	Config *graphql.ArgumentConfig
	Ast    *ast.InputValueDefinition
}

type InputObjectDirectiveContext struct {
	Name   string
	Config *graphql.InputObjectConfig
	Ast    *ast.InputObjectDefinition
}

type InputValueDirectiveContext struct {
	Name  string
	Field *graphql.InputObjectField
	Ast   *ast.InputValueDefinition
}

type InterfaceDirectiveContext struct {
	Name   string
	Config *graphql.InterfaceConfig
	Ast    *ast.InterfaceDefinition
}

type ObjectDirectiveContext struct {
	Name   string
	Config *graphql.ObjectConfig
	Ast    *ast.ObjectDefinition
}

type ScalarDirectiveContext struct {
	Name   string
	Config *graphql.ScalarConfig
	Ast    *ast.ScalarDefinition
}

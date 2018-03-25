package parser

import (
	"github.com/graphql-go/graphql/language/ast"
)

func locToString(loc *ast.Location) string {
	return "" // string(loc.Source.Body[loc.Start:loc.End])
}

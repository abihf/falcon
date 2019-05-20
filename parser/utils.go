package parser

import (
	"github.com/graphql-go/graphql/language/ast"
)

func locToString(loc *ast.Location) string {
	return "" // string(loc.Source.Body[loc.Start:loc.End])
}

func stringValue(sv *ast.StringValue) string {
	if sv != nil {
		return sv.Value
	}
	return ""
}

func astValue(av ast.Value) interface{} {
	if av != nil {
		return av.GetValue()
	}
	return nil
}

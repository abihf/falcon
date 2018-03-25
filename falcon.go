package falcon

import (
	"github.com/graphql-go/graphql"
	gqlParser "github.com/graphql-go/graphql/language/parser"
)

type Resolver map[string]interface{}

func CreateSchema(schemaString string, resolver Resolver) (*graphql.Schema, error) {
	doc, err := gqlParser.Parse(gqlParser.ParseParams{
		Source: schemaString,
		Options: gqlParser.ParseOptions{
			NoLocation: false,
			NoSource:   false,
		},
	})
	if err != nil {
		return nil, err
	}
	p := createParser(doc, resolver)
	return p.getSchema()
}

package falcon

import (
	"fmt"
	"reflect"

	"github.com/abihf/falcon-graphql/parser"
	"github.com/graphql-go/graphql"
)

type Resolver map[string]interface{}

func CreateSchema(source interface{}, resolver Resolver) (*graphql.Schema, error) {
	p := parser.Create(resolver)
	switch source.(type) {
	case string:
		if err := p.ParseSource("GraphQL", []byte(source.(string))); err != nil {
			return nil, err
		}

	case map[string][]byte:
		for name, body := range source.(map[string][]byte) {
			if err := p.ParseSource(name, body); err != nil {
				return nil, err
			}
		}

	default:
		return nil, fmt.Errorf(
			"Invalid source type.\nExepcted string or map[string][]byte, got %s",
			reflect.TypeOf(source).String(),
		)
	}

	return p.GetSchema()
}

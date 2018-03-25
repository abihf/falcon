package falcon

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

func (p *parser) getSchema() (*graphql.Schema, error) {
	d, ok := p.definitions["_schema"]
	if !ok {
		return nil, errors.New("Can not found schema definition")
	}
	return p.parseSchema(d.(*ast.SchemaDefinition))
}

func (p *parser) parseSchema(def *ast.SchemaDefinition) (*graphql.Schema, error) {
	var queryObject *graphql.Object
	var mutationObject *graphql.Object
	var subscriptionObject *graphql.Object

	for _, operation := range def.OperationTypes {
		parsed, err := p.getParsed(operation.Type.Name.Value, kinds.ObjectDefinition)
		if err != nil {
			return nil, fmt.Errorf("Invalid query operation type for {%s} caused by %s", locToString(operation.GetLoc()), err.Error())
		}
		obj := parsed.(*graphql.Object)
		switch operation.Operation {
		case ast.OperationTypeQuery:
			queryObject = obj
		case ast.OperationTypeMutation:
			mutationObject = obj
		case ast.OperationTypeSubscription:
			subscriptionObject = obj
		}
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        queryObject,
		Mutation:     mutationObject,
		Subscription: subscriptionObject,
	})
	return &schema, err
}

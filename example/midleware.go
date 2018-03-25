package main

import "github.com/graphql-go/graphql"

func appendMidleware(
	param graphql.ResolveParams,
	args map[string]interface{},
	next graphql.FieldResolveFn,
) (interface{}, error) {
	suffix := args["suffix"].(string)
	res, err := next(param)
	if err != nil {
		return nil, err
	}
	return res.(string) + suffix, nil
}

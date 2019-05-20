package resolver

import "github.com/graphql-go/graphql"

func init() {
	resolvers.Add("QueryRoot", "node", queryNode)
}

func queryNode(param graphql.ResolveParams) (interface{}, error) {
	return &user{
		TypeName: "User",
		ID:       "1",
		Name:     "Abi",
	}, nil
}

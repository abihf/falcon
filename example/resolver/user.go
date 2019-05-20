package resolver

import (
	"github.com/graphql-go/graphql"
)

func init() {
	resolvers.RegisterType("User", &user{})
	resolvers.Add("QueryRoot", "me", queryUser)
	resolvers.Add("User", "friends", getUserFriends)
}

type user struct {
	TypeName string
	ID       string `graphql:"id"`
	Name     string `graphql:"name"`
}

func queryUser(param graphql.ResolveParams) (interface{}, error) {
	return &user{
		TypeName: "User",
		ID:       "1",
		Name:     "Abi",
	}, nil
}

func getUserFriends(param graphql.ResolveParams) (interface{}, error) {
	return []*user{
		{
			TypeName: "User",
			ID:       "1",
			Name:     "Abi",
		},
	}, nil
}

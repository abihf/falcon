package resolver

import (
	"github.com/abihf/falcon-graphql"
)

var resolvers = falcon.Resolver{}

func Get() falcon.Resolver {
	return resolvers
}

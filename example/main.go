package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/abihf/falcon-graphql"
	"github.com/abihf/falcon-graphql/example/resolver"
	"github.com/graphql-go/graphql"
)

type user struct {
	Name string `json:"name"`
}

func main() {
	schemaString := `
schema {
	query: QueryRoot
}

type QueryRoot {
	me: User
	node(id: ID!): Node
}

type User implements Node {
	id: ID!
	name: String!
	size(unit: LongUnit = METER): Float
	friends: [User!]!
}

interface Node {
	id: ID!
}

enum LongUnit {
	INCH @deprecated(reason: "Use other")
	CENTI_METER
	METER
}
	`
	resolvers := resolver.Get()
	schema, err := falcon.CreateSchema(schemaString, resolvers)
	if err != nil {
		log.Fatalf("Failed to parse schema\n%s", err.Error())
	}

	// Query
	query := `
		query GetNode{
			node(id: "1") {
				id
				...SomeFragment
			}
		}

		fragment SomeFragment on User {
			__typename
			name
			friends {
				name
			}
		}
	`
	params := graphql.Params{
		Schema:        *schema,
		RequestString: query,
	}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("result: %s \n", rJSON)
}

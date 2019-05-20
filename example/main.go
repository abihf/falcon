package main

import (
	"log"
	"net/http"

	"github.com/abihf/falcon-graphql"
	"github.com/abihf/falcon-graphql/example/resolver"
	"github.com/graphql-go/handler"
)

type user struct {
	Name string `json:"name"`
}

func main() {
	resolvers := resolver.Get()
	resolvers.ApplyDirectiveMidleware("append", appendMidleware)
	schema, err := falcon.CreateSchema(map[string][]byte{
		"1.gql": schemaBody1,
		"2.gql": schemaBody2,
	}, resolvers)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", handler.New(&handler.Config{
		Schema:   schema,
		GraphiQL: true,
		Pretty:   true,
	}))
}

var schemaBody1 = []byte(`
schema {
  query: QueryRoot
}

type QueryRoot {
	"Get cuurrent user"
	me: User
	
	"random node?"
  node(id: ID!): Node
}
`)

var schemaBody2 = []byte(`
directive @deprecated(
  reason: String = "No longer supported"
) on FIELD_DEFINITION | ENUM_VALUE

"tore user information"
type User implements Node {
	# user id
	id: ID!

	# user name
	name: String! @append(suffix:"Hafshin")

	# friend list
  friends: [User!]!

	size(unit: LongUnit = METER): Float @deprecated(reason: "Why?")
}

interface Node {
  id: ID!
}

enum LongUnit {
  INCH @deprecated(reason: "Use other")
  CENTI_METER
  METER
}
`)

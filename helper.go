package falcon

import (
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

// RegisterType is helper function to register a struct.
// Example:
//  type User struct {
//		ID   string `graphql:"id"`
//    Name string `graphql:"name"`
//  }
func (r Resolver) RegisterType(objectName string, v interface{}) {
	r[objectName+".isTypeOf"] = CreateTypeChecker(v)

	t := reflect.TypeOf(v).Elem()
	fieldCount := t.NumField()
	for index := 0; index < fieldCount; index++ {
		(func(i int) {
			field := t.Field(i)
			tagParam := strings.SplitN(field.Tag.Get("graphql"), ",", 1)
			if len(tagParam) > 0 && tagParam[0] != "" {
				r[objectName+"."+tagParam[0]] =
					func(param graphql.ResolveParams) (interface{}, error) {
						return reflect.ValueOf(param.Source).Elem().Field(i).Interface(), nil
					}
			}
		})(index)
	}
}

func CreateTypeChecker(v interface{}) graphql.IsTypeOfFn {
	t := reflect.TypeOf(v)
	return func(param graphql.IsTypeOfParams) bool {
		return reflect.TypeOf(param.Value).ConvertibleTo(t)
	}
}

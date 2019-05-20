package directives

type VisitorParam struct {
	Name    string
	Kind    string
	Context interface{}
	Args    map[string]interface{}
}
type Visitor func(param *VisitorParam) error

type VisitorObject interface {
	Object(c *ObjectDirectiveContext, args map[string]interface{}) error
	Field(c *FieldDirectiveContext, args map[string]interface{}) error
	FieldArg(c *FieldArgDirectiveContext, args map[string]interface{}) error
	Interface(c *InterfaceDirectiveContext, args map[string]interface{}) error
	Enum(c *EnumDirectiveContext, args map[string]interface{}) error
	EnumValue(c *EnumValueDirectiveContext, args map[string]interface{}) error
	InputObject(c *InputObjectDirectiveContext, args map[string]interface{}) error
	InputValue(c *InputValueDirectiveContext, args map[string]interface{}) error
	Scalar(c *ScalarDirectiveContext, args map[string]interface{}) error
}

func CreateCallback(visitor VisitorObject) Visitor {
	return func(param *VisitorParam) error {
		args := param.Args
		switch ctx := param.Context.(type) {
		case *ObjectDirectiveContext:
			return visitor.Object(ctx, args)
		case *FieldDirectiveContext:
			return visitor.Field(ctx, args)
		case *FieldArgDirectiveContext:
			return visitor.FieldArg(ctx, args)
		case *InterfaceDirectiveContext:
			return visitor.Interface(ctx, args)
		case *EnumDirectiveContext:
			return visitor.Enum(ctx, args)
		case *EnumValueDirectiveContext:
			return visitor.EnumValue(ctx, args)
		case *InputObjectDirectiveContext:
			return visitor.InputObject(ctx, args)
		case *InputValueDirectiveContext:
			return visitor.InputValue(ctx, args)
		case *ScalarDirectiveContext:
			return visitor.Scalar(ctx, args)
		default:
			return nil
		}
	}
}

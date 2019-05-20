package directives

type DefaultVisitor struct{}

var _ VisitorObject = &DefaultVisitor{}

func (v *DefaultVisitor) Object(c *ObjectDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) Field(c *FieldDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) FieldArg(c *FieldArgDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) Interface(c *InterfaceDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) Enum(c *EnumDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) EnumValue(c *EnumValueDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) InputObject(c *InputObjectDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) InputValue(c *InputValueDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (v *DefaultVisitor) Scalar(c *ScalarDirectiveContext, args map[string]interface{}) error {
	return nil
}

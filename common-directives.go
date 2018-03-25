package falcon

type deprecationVisitor struct{}

func CreateDeprecationCallback() DirectiveCallback {
	return CreateCallback(&deprecationVisitor{})
}

func (d *deprecationVisitor) Enum(dirContext *EnumDirectiveContext, args map[string]interface{}) error {
	return nil
}

func (d *deprecationVisitor) EnumValue(dirContext *EnumValueDirectiveContext, args map[string]interface{}) error {
	dirContext.Config.DeprecationReason = d.getDeprecationReason(args)
	return nil
}

func (d *deprecationVisitor) getDeprecationReason(args map[string]interface{}) string {
	reason, ok := args["reason"]
	if !ok {
		return "deprecated"
	}

	str, ok := reason.(string)
	if !ok {
		return "deprecated" // should throw error
	}

	return str
}

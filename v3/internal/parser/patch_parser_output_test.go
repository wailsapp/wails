package parser

func patchParserOutput(p *Project) {
	// Patch Package data in BoundMethods
	for _, packageData := range p.BoundMethods {
		for _, boundMethods := range packageData {
			for _, boundMethod := range boundMethods {
				for _, param := range boundMethod.Inputs {
					patchParameterType(param.Type)
				}

				for _, param := range boundMethod.Outputs {
					patchParameterType(param.Type)
				}
			}
		}
	}

	// Patch Package data in Models
	for _, packageData := range p.Models {
		for _, structDef := range packageData {
			for _, field := range structDef.Fields {
				patchParameterType(field.Type)
			}
		}
	}
}

func patchParameterType(t *ParameterType) {
	if t == nil {
		return
	}

	t.Package = nil
	patchParameterType(t.MapKey)
	patchParameterType(t.MapValue)
}

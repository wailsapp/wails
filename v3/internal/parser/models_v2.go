package parser

import "go/types"

func (p *Parameter) Models() (models []*types.Named) {
	return p.modelsIn(p.Type())

}

func (p *Parameter) modelsIn(t types.Type) (models []*types.Named) {
	// TODO: break cyclic references
	project := p.Parent().Parent()

	for {
		switch x := t.(type) {
		case *types.Basic:
			return
		case *types.Slice:
			t = x.Elem()
		case *types.Map:
			t = x.Elem()
		case *types.Named:
			models = append(models, x)
			models = append(models, p.modelsInNamed(x)...)
			return
		case *types.Struct:
			models = append(models, types.NewNamed(types.NewTypeName(0, nil, project.anonymousStructID(x), nil), x, nil))
			models = append(models, p.modelsInStruct(x)...)
			return
		default:
			return
		}

	}
}

func (p *Parameter) modelsInNamed(n *types.Named) (models []*types.Named) {
	switch x := n.Underlying().(type) {
	case *types.Struct:
		models = append(models, p.modelsInStruct(x)...)
	}
	return
}

func (p *Parameter) modelsInStruct(s *types.Struct) (models []*types.Named) {
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		models = append(models, p.modelsIn(field.Type())...)
	}
	return
}

func (m *BoundMethod) Models() (models []*types.Named) {
	for _, param := range m.JSInputs() {
		models = append(models, param.Models()...)
	}
	for _, param := range m.JSOutputs() {
		models = append(models, param.Models()...)
	}
	return
}

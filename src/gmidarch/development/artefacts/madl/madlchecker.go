package madl

import (
	"github.com/gfads/midarch/src/gmidarch/development/repositories/architectural"
	"github.com/gfads/midarch/src/shared"
)

type MADLChecker interface {
	SyntaxCheck(MADL)
	SemanticCheck(MADL, architectural.ArchitecturalRepository)
}

type MADLCheckerImpl struct{}

func NewMADLChecker() MADLChecker {

	return MADLCheckerImpl{}
}

func (MADLCheckerImpl) SyntaxCheck(m MADL) {

	// Check if all components/connectors were declared
	for a := range m.Attachments {

		if !m.IsInComponents(m.Attachments[a].C1) {
			shared.ErrorHandler(shared.GetFunction(), "Component '"+m.Attachments[a].C1.Id+"' was not Declared!!")
		}

		if !m.IsInConnectors(m.Attachments[a].T) {
			shared.ErrorHandler(shared.GetFunction(), "Connector '"+m.Attachments[a].T.Id+"' was not Declared!!")
		}

		if !m.IsInComponents(m.Attachments[a].C2) {
			shared.ErrorHandler(shared.GetFunction(), "Component '"+m.Attachments[a].C2.Id+"' was not Declared!!")
		}
	}

	// Check if all components were used
	for c := range m.Components {
		if !m.IsComponentInAttachments(m.Components[c]) {
			shared.ErrorHandler(shared.GetFunction(), "Component '"+m.Components[c].Id+"' declared, but not Used!!")
		}
	}

	// Check if all connectors were used
	for t := range m.Connectors {
		if !m.IsConnectorInAttachments(m.Connectors[t]) {
			shared.ErrorHandler(shared.GetFunction(), "Connector '"+m.Connectors[t].Id+"' declared, but not Used!!")
		}
	}
}

func (c MADLCheckerImpl) SemanticCheck(m MADL, archRepo architectural.ArchitecturalRepository) {

	arm := architectural.NewArchitecturalRepositoryManager()

	// Check if all component/connectors types used in the madl exist in architectural repositories
	for c := range m.Components {
		if !arm.TypeExist(m.Components[c].TypeName) {
			shared.ErrorHandler(shared.GetFunction(), "Component type '"+m.Components[c].TypeName+"' does not exist in the architectural repositories!!")
		}
	}
	for c := range m.Connectors {
		if !arm.TypeExist(m.Connectors[c].TypeName) {
			shared.ErrorHandler(shared.GetFunction(), "Connector type '"+m.Connectors[c].TypeName+"' does not exist in the architectural repositories!!")
		}
	}

	// Check arity of connectors
	for i := range m.Connectors {
		lArityMADL := m.CountArity(m.Connectors[i].Id, shared.LEFT_ARITY)
		rArityMADL := m.CountArity(m.Connectors[i].Id, shared.RIGHT_ARITY)

		lArityRepo, rArityRepo := arm.GetConnectorDefaultArities(m.Connectors[i].TypeName)

		if lArityMADL > lArityRepo || rArityMADL > rArityRepo {
			shared.ErrorHandler(shared.GetFunction(), "Arity of connector '"+m.Connectors[i].Id+"' is incompatible with arity's type")
		}
	}

	// Check duplicity of component's ids
	for i := 0; i < len(m.Components); i++ {
		for j := i + 1; j < len(m.Components); j++ {
			if m.Components[i].Id == m.Components[j].Id {
				shared.ErrorHandler(shared.GetFunction(), "Component '"+m.Components[i].Id+"' is duplicated!!")
			}
		}
	}

	// Check duplicty of connector'' ids
	for i := 0; i < len(m.Connectors); i++ {
		for j := i + 1; j < len(m.Connectors); j++ {
			if m.Connectors[i].Id == m.Connectors[j].Id {
				shared.ErrorHandler(shared.GetFunction(), "Connector '"+m.Connectors[i].Id+"' is duplicated!!")
			}

		}
	}
}

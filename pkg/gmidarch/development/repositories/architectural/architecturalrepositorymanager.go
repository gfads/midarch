package architectural

import (
	"github.com/gfads/midarch/pkg/shared"
)

type ArchitecturalRepositoryManager interface {
	GetRepository() ArchitecturalRepository
	GetConnectorDefaultArities(string) (int, int)
	TypeExist(string) bool
}

type ArchitecturalRepositoryManagerImpl struct {
	Repository ArchitecturalRepository
}

func NewArchitecturalRepositoryManager(businessComponents map[string]interface{}) ArchitecturalRepositoryManager {
	var arm ArchitecturalRepositoryManager

	// Create a repositories to be managed
	armImpl := ArchitecturalRepositoryManagerImpl{}
	armImpl.Repository = LoadArchitecturalRepository(businessComponents)

	arm = armImpl

	return arm
}

func (armImpl ArchitecturalRepositoryManagerImpl) GetRepository() ArchitecturalRepository {
	return armImpl.Repository
}

func (armImpl ArchitecturalRepositoryManagerImpl) TypeExist(t string) bool {
	foundComp := false
	foundConn := false

	compLibrary := armImpl.Repository.CompLibrary
	for i := range compLibrary {
		if t == i {
			foundComp = true
		}
	}

	connLibrary := armImpl.Repository.ConnLibrary
	for i := range connLibrary {
		if t == i {
			foundConn = true
		}
	}

	if foundComp || foundConn {
		return true
	} else {
		return false
	}
}

func (armImpl ArchitecturalRepositoryManagerImpl) GetConnectorDefaultArities(t string) (int, int) { // TODO
	lArity := 0
	rArity := 0
	foundConn := false

	connLibrary := armImpl.Repository.ConnLibrary
	for i := range connLibrary {
		if t == i {
			foundConn = true
			lArity = connLibrary[i].DefaultLeftArity
			rArity = connLibrary[i].DefaultRightArity
		}
	}

	if !foundConn {
		shared.ErrorHandler(shared.GetFunction(), "Connector type'"+t+"' does not exist. Not possible to define its arities!")
	}

	return lArity, rArity
}

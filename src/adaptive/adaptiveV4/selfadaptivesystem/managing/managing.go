package managing

import (
	"adaptive/adaptiveV4/selfadaptivesystem/managed"
	"sync"
)

type ManagingSystem interface {
	Start(*sync.WaitGroup)
}
type ManagingSystemImpl struct {
	Managed managed.Managed
	Mapek   MAPEK
}

func NewManagingSystem(ms managed.Managed, mapek *MAPEK) ManagingSystem {
	var r ManagingSystem

	// Configure managed system
	if mapek.M != nil {
		mapek.M.SetManaged(ms)
	}
	if mapek.E != nil {
		mapek.E.SetManaged(ms)
	}

	// Configure MAPEK
	if mapek.M != nil {
		mapek.M.SetMAPEK(*mapek)
	}
	if mapek.A != nil {
		mapek.A.SetMAPEK(*mapek)
	}
	if mapek.P != nil {
		mapek.P.SetMAPEK(*mapek)
	}
	if mapek.E != nil {
		mapek.E.SetMAPEK(*mapek)
	}

	r = &ManagingSystemImpl{Managed: ms, Mapek: *mapek}

	return r
}

func (ms *ManagingSystemImpl) Start(wg *sync.WaitGroup) {

	if ms.Mapek.M != nil {
		go ms.Mapek.M.Start()
	}
	if ms.Mapek.A != nil {
		go ms.Mapek.A.Start()
	}
	wg.Done()
}

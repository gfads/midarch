package managing

import (
	"adaptive/adaptiveV5/selfadaptivesystem/managed"
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

	// Configure managed system (Monitor and Executor only as they are the ones to interact with the managed system)
	if mapek.M != nil {
		mapek.M.SetManaged(ms)
	}
	if mapek.E != nil {
		mapek.E.SetManaged(ms)
	}

	// Configure MAPEK of all managing system's components
	if mapek.M != nil {   // Monitor
		mapek.M.SetMAPEK(*mapek)
	}
	if mapek.A != nil {   // Analyser
		mapek.A.SetMAPEK(*mapek)
	}
	if mapek.P != nil {   // Planner
		mapek.P.SetMAPEK(*mapek)
	}
	if mapek.E != nil {   // Executor
		mapek.E.SetMAPEK(*mapek)
	}

	r = &ManagingSystemImpl{Managed: ms, Mapek: *mapek}

	return r
}

func (ms *ManagingSystemImpl) Start(wg *sync.WaitGroup) {

	// Start Monitor
	if ms.Mapek.M != nil {
		go ms.Mapek.M.Start()
	}

	// Start Analyser
	if ms.Mapek.A != nil {
		go ms.Mapek.A.Start()
	}

	// Start Planner
	if ms.Mapek.P != nil {
		go ms.Mapek.P.Start()
	}

	// Start Executor
	if ms.Mapek.E != nil {
		go ms.Mapek.E.Start()
	}

	wg.Done()
}

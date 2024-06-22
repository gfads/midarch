package adaptive

import (
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
)

var isFirstTimeEvolutiveProtocol = true

// @Type: EvolutiveProtocol
// @Behaviour: Behaviour = I_Hasnewprotocol -> InvR.e1 -> Behaviour
type EvolutiveProtocol struct{}

func (EvolutiveProtocol) I_Hasnewprotocol(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	if isFirstTimeEvolutiveProtocol {
		time.Sleep(shared.FIRST_MONITOR_TIME) // only first time
		isFirstTimeEvolutiveProtocol = false
	} else {
		time.Sleep(shared.MONITOR_TIME)
	}

	// return from this point if no new components detected
	if len(shared.ListOfComponentsToAdaptTo) == 0 {
		*reset = true
		return
	}

	*msg = messages.SAMessage{Payload: shared.ListOfComponentsToAdaptTo}

	shared.ListOfComponentsToAdaptTo = []string{}
}

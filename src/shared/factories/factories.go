package factories

import (
	"gmidarch/development/components"
	"shared"
)

func LocateNaming() components.Namingproxy {
	p := components.Namingproxy{Host: shared.NAMING_HOST,Port: shared.NAMING_PORT}
	return p
}

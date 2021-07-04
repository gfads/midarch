package madl

import (
	"gmidarch/development/components/component"
	"gmidarch/development/connectors"
)

type Attachment struct {
	C1 component.Component
	T  connectors.Connector
	C2 component.Component
}


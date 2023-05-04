package madl

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/components/component"
	"github.com/gfads/midarch/pkg/gmidarch/development/connectors"
)

type Attachment struct {
	C1 component.Component
	T  connectors.Connector
	C2 component.Component
}

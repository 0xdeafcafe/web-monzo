package models

import (
	"github.com/0xdeafcafe/gomonzo"
	"github.com/maxwellhealth/bongo"
)

// Context is the application context
type Context struct {
	Mongo *bongo.Connection
	Monzo *gomonzo.GoMonzo
}

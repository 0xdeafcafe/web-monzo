package models

import (
	"github.com/0xdeafcafe/gomonzo"
	"github.com/jinzhu/gorm"
)

// Context is the application context
type Context struct {
	DB    *gorm.DB         `json:"-"`
	Monzo *gomonzo.GoMonzo `json:"-"`
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Audit contains the base info of any model in the database
type Audit struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// Init set's a models base information
func (audit *Audit) Init() {
	now := time.Now().UTC()

	audit.ID = uuid.NewV4().String()
	audit.CreatedAt = now
	audit.UpdatedAt = now
}

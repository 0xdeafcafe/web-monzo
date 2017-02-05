package handlers

import (
	"github.com/0xdeafcafe/web-monzo/models"
	"gopkg.in/gin-gonic/gin.v1"
)

// HomeHandler ..
type HomeHandler struct {
	Context *models.Context
}

// Index is the index route of the Home Handler
func (handler HomeHandler) Index(c *gin.Context) {
	//cc := c.(models.Context)
}

// NewHomeHandler creates a new HomeHandler and registers the reqired routes
func NewHomeHandler(r *gin.Engine, c *models.Context) {
	handler := new(HomeHandler)
	handler.Context = c

	r.GET("/", handler.Index)
}

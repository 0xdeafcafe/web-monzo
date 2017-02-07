package handlers

import (
	"net/http"

	"fmt"

	"github.com/0xdeafcafe/web-monzo/models"
	"github.com/gin-contrib/sessions"
	"gopkg.in/gin-gonic/gin.v1"
)

// HomeHandler ..
type HomeHandler struct {
	Context *models.Context
}

// Index is the index route of the Home Handler
func (hndlr HomeHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	webSession := models.GetValidWebSession(hndlr.Context.Mongo, session.Get("webSessionID"), c.ClientIP())
	_, webSession = webSession.Refresh(hndlr.Context.Mongo, hndlr.Context.Monzo, session)
	monzoAccounts, _, err := hndlr.Context.Monzo.ListAccounts(webSession.ToToken())
	if err != nil {
		panic(err)
	}
	if len(monzoAccounts.Accounts) <= 0 {
		panic("no accounts")
	}

	primaryAccount := monzoAccounts.Accounts[0]
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/accounts/%s/transactions", primaryAccount.ID))
	return
}

// NewHomeHandler creates a new HomeHandler and registers the reqired routes
func NewHomeHandler(r *gin.Engine, c *models.Context) {
	handler := new(HomeHandler)
	handler.Context = c

	r.GET("/", handler.Index)
}

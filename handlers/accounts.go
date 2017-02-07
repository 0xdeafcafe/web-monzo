package handlers

import (
	"net/http"
	"time"

	monzoModels "github.com/0xdeafcafe/gomonzo/models"
	"github.com/0xdeafcafe/web-monzo/models"
	"github.com/gin-contrib/sessions"
	"gopkg.in/gin-gonic/gin.v1"
)

// AccountsHandler ..
type AccountsHandler struct {
	Context *models.Context
}

// ListTransactions ..
func (hndlr AccountsHandler) ListTransactions(c *gin.Context) {
	accountID := c.Param("account_id")
	session := sessions.Default(c)
	webSession := models.GetValidWebSession(hndlr.Context.Mongo, session.Get("webSessionID"), c.ClientIP())
	if webSession == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	token := webSession.ToToken()
	monzoAccounts, _, err := hndlr.Context.Monzo.ListAccounts(token)
	if err != nil {
		panic(err)
	}

	account := new(monzoModels.Account)
	for _, acc := range monzoAccounts.Accounts {
		if acc.ID == accountID {
			account = &acc
		}
	}
	if account == nil {
		panic("account not found")
	}

	mostRecentTx := models.MostRecentTransaction(hndlr.Context.Mongo, account.ID)
	since := time.Date(1994, 8, 18, 0, 0, 0, 0, time.UTC)
	if mostRecentTx != nil {
		since = mostRecentTx.Created
	}
	monzoTransactions, _, err := hndlr.Context.Monzo.ListTransactionsSince(token, account.ID, since)
	if err != nil {
		panic(err)
	}
	for _, transaction := range monzoTransactions.Transactions {
		models.AddTransaction(hndlr.Context.Mongo, &transaction)
	}

	c.HTML(http.StatusOK, "accounts/transactions", gin.H{
		"title":        "Home",
		"accounts":     &monzoAccounts.Accounts,
		"account":      &account,
		"transactions": models.GetTransactionsByAccountID(hndlr.Context.Mongo, account.ID, 100),
		"flash":        session.Flashes(),
	})
}

// NewAccountsHandler ..
func NewAccountsHandler(e *gin.Engine, c *models.Context) {
	ctrl := new(AccountsHandler)
	ctrl.Context = c

	e.GET("accounts/:account_id/transactions", ctrl.ListTransactions)
}

package handlers

import (
	"net/http"
	"time"

	"fmt"

	monzoModels "github.com/0xdeafcafe/gomonzo/models"
	"github.com/0xdeafcafe/web-monzo/models"
	"github.com/gin-contrib/sessions"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2/bson"
)

// AccountsHandler ..
type AccountsHandler struct {
	Context *models.Context
}

// ListTransactions ..
func (hndlr AccountsHandler) ListTransactions(c *gin.Context) {
	session := sessions.Default(c)
	webSession, err := models.GetWebSession(hndlr.Context.Mongo, hndlr.Context.Monzo, session, c.ClientIP())
	if err != nil {
		session.AddFlash(err)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	token := webSession.ToToken()
	monzoAccounts, _, err := hndlr.Context.Monzo.ListAccounts(token)
	if err != nil {
		panic(err)
	}

	accountID := c.Param("account_id")
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
		"title":        "Transactions",
		"accounts":     &monzoAccounts.Accounts,
		"account":      &account,
		"transactions": models.GetTransactionsByAccountID(hndlr.Context.Mongo, account.ID, 100),
		"flash":        session.Flashes(),
	})
}

// ViewTransaction ..
func (hndlr AccountsHandler) ViewTransaction(c *gin.Context) {
	session := sessions.Default(c)
	webSession, err := models.GetWebSession(hndlr.Context.Mongo, hndlr.Context.Monzo, session, c.ClientIP())
	if err != nil {
		session.AddFlash(err)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	token := webSession.ToToken()
	monzoAccounts, _, err := hndlr.Context.Monzo.ListAccounts(token)
	if err != nil {
		panic(err)
	}

	accountID := c.Param("account_id")
	account := new(monzoModels.Account)
	for _, acc := range monzoAccounts.Accounts {
		if acc.ID == accountID {
			account = &acc
		}
	}
	if account == nil {
		panic("account not found")
	}

	monzoTransaction, monzoError, err := hndlr.Context.Monzo.GetTransaction(token, c.Param("transaction_id"))
	if err != nil {
		if err.Error() == "not_found.transaction_not_found" {
			session.AddFlash(models.NewWarningFlash(monzoError.Message, monzoError.Code))
			session.Save()
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/accounts/%s/transactions", accountID))
			return
		}

		panic(err)
	}

	models.AddTransaction(hndlr.Context.Mongo, monzoTransaction)
	var transaction models.Transaction
	hndlr.Context.Mongo.Collection(models.TransactionCollectionName).FindOne(bson.M{"transactionid": monzoTransaction.ID}, &transaction)
	fmt.Println(transaction)
	c.HTML(http.StatusOK, "accounts/transaction", gin.H{
		"title":       fmt.Sprintf("%s - Transaction", transaction.FriendlyName()),
		"accounts":    &monzoAccounts.Accounts,
		"account":     &account,
		"transaction": &transaction,
		"flash":       session.Flashes(),
	})
}

// NewAccountsHandler ..
func NewAccountsHandler(e *gin.Engine, c *models.Context) {
	ctrl := new(AccountsHandler)
	ctrl.Context = c

	e.GET("accounts/:account_id/transactions", ctrl.ListTransactions)
	e.GET("accounts/:account_id/transactions/:transaction_id", ctrl.ViewTransaction)
}

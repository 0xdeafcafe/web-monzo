package models

import (
	"hash/crc32"
	"time"

	"gopkg.in/mgo.v2/bson"

	"fmt"

	"strings"

	"math"

	"github.com/0xdeafcafe/gomonzo/models"
	raven "github.com/getsentry/raven-go"
	"github.com/maxwellhealth/bongo"
)

const (
	// TransactionCollectionName is the name of the MongoDB collection containing the Transaction info
	TransactionCollectionName = "transactions"
)

// Transaction ..
type Transaction struct {
	bongo.DocumentBase `bson:",inline"`

	TransactionID     string
	AccountID         string
	Description       string
	Merchant          *models.Merchant
	IsLoad            bool
	Currency          string
	Amount            int64
	AccountBalance    int64
	Metadata          map[string]string
	Notes             string
	DeclineReason     string
	LocalAmount       int64
	LocalCurrency     string
	Scheme            string
	DedupeID          string
	Originator        bool
	IncludeInSpending bool
	Created           time.Time
	Updated           *time.Time
	Settled           string
}

// Update modifies an existing transaction setting the fields that are allowed to change
func (transaction Transaction) Update(connection *bongo.Connection, newTransaction models.Transaction) {
	transaction.Description = newTransaction.Description
	transaction.Merchant = newTransaction.Merchant
	transaction.Currency = newTransaction.Currency
	transaction.Amount = newTransaction.Amount
	transaction.Metadata = newTransaction.Metadata
	transaction.Notes = newTransaction.Notes
	transaction.DeclineReason = newTransaction.DeclineReason
	transaction.IncludeInSpending = newTransaction.IncludeInSpending
	transaction.Updated = newTransaction.Updated
	transaction.Settled = newTransaction.Settled

	err := connection.Collection(TransactionCollectionName).Save(&transaction)
	if err != nil {
		raven.CaptureError(err, nil, nil)
	}
}

// AmountInteger ..
func (transaction Transaction) AmountInteger() string {
	amountPrecise := float64(transaction.Amount) / 100
	amount := math.Abs(amountPrecise) // to remove potential negative
	return fmt.Sprintf("%d.", int64(amount))
}

// AmountFractional ..
func (transaction Transaction) AmountFractional() string {
	amount := float64(transaction.Amount) / 100
	amountStr := fmt.Sprintf("%.2f", amount)
	return strings.Split(amountStr, ".")[1]
}

// FriendlyName ..
func (transaction Transaction) FriendlyName() string {
	if transaction.Merchant != nil {
		return transaction.Merchant.Name
	}
	return transaction.Description
}

// HasLogo ..
func (transaction Transaction) HasLogo() bool {
	return transaction.Merchant != nil
}

// LogoOrHex ..
func (transaction Transaction) LogoOrHex() string {
	if transaction.HasLogo() {
		return transaction.Merchant.Logo
	}

	// Create RGB mapping
	hash := crc32.ChecksumIEEE([]byte(transaction.Description))
	r := (hash & 0xFF0000) >> 16
	g := (hash & 0x00FF00) >> 8
	b := (hash & 0x0000FF) >> 16

	// Mix colours into aesthetically pleasing pallet
	r = (r + 0xff) / 2
	g = (g + 0xff) / 2
	b = (b + 0xff) / 2

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// HasPhysicalLocation ..
func (transaction Transaction) HasPhysicalLocation() bool {
	return transaction.Merchant != nil && transaction.Merchant.Address != nil
}

// NewTransaction creates a new transaction from a Monzo Transaction.
func NewTransaction(trans *models.Transaction) *Transaction {
	return &Transaction{
		TransactionID:     trans.ID,
		AccountID:         trans.AccountID,
		Description:       trans.Description,
		Merchant:          trans.Merchant,
		IsLoad:            trans.IsLoad,
		Currency:          trans.Currency,
		Amount:            trans.Amount,
		AccountBalance:    trans.AccountBalance,
		Metadata:          trans.Metadata,
		Notes:             trans.Notes,
		DeclineReason:     trans.DeclineReason,
		LocalAmount:       trans.LocalAmount,
		LocalCurrency:     trans.LocalCurrency,
		Scheme:            trans.Scheme,
		DedupeID:          trans.DedupeID,
		Originator:        trans.Originator,
		IncludeInSpending: trans.IncludeInSpending,
		Created:           trans.Created,
		Updated:           trans.Updated,
		Settled:           trans.Settled,
	}
}

// AddTransaction ..
func AddTransaction(connection *bongo.Connection, transaction *models.Transaction) {
	var existingTransaction Transaction
	found := true
	err := connection.Collection(TransactionCollectionName).FindOne(bson.M{"transactionid": transaction.ID}, &existingTransaction)
	if err != nil {
		found = false

		if dnfError, ok := err.(*bongo.DocumentNotFoundError); !ok {
			raven.CaptureError(dnfError, nil, nil)
		}
	}

	if !found {
		newTransaction := NewTransaction(transaction)
		connection.Collection(TransactionCollectionName).Save(newTransaction)
		return
	}

	// Check if we gotta update info - two ifs for readability
	update := false
	if existingTransaction.Updated == nil && transaction.Updated != nil {
		update = true
	}
	if existingTransaction.Updated != nil && transaction.Updated != nil && transaction.Updated.After(*existingTransaction.Updated) {
		update = true
	}

	if update {
		existingTransaction.Update(connection, *transaction)
	}
}

// MostRecentTransaction ..
func MostRecentTransaction(connection *bongo.Connection, accountID string) *Transaction {
	mostRecentTx := connection.Collection(TransactionCollectionName).Find(bson.M{"accountid": accountID})
	mostRecentTx.Query.Limit(1)
	mostRecentTx.Query.Sort("-created")

	var transaction Transaction
	if mostRecentTx.Next(&transaction) {
		return &transaction
	}

	return nil
}

// GetTransactionsByAccountID ..
func GetTransactionsByAccountID(connection *bongo.Connection, accountID string, count int) []Transaction {
	transactionsSet := connection.Collection(TransactionCollectionName).Find(bson.M{"accountid": accountID})
	transactionsSet.Query.Sort("-created")
	if count > -1 {
		transactionsSet.Query.Limit(count)
	}

	var transactions []Transaction
	var transaction Transaction
	for transactionsSet.Next(&transaction) {
		transactions = append(transactions, transaction)
	}

	return transactions
}

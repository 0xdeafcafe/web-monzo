package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"errors"

	"github.com/0xdeafcafe/gomonzo"
	"github.com/0xdeafcafe/gomonzo/models"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/sessions"
	"github.com/maxwellhealth/bongo"
)

const (
	// WebSessionCollectionName is the name of the MongoDB collection containing the WebSession info
	WebSessionCollectionName = "web_sessions"
)

// WebSession defines a session
type WebSession struct {
	bongo.DocumentBase `bson:",inline"`

	UserID   string
	ClientID string
	IP       string

	TokenType    string
	AccessToken  string `gorm:"type:varchar(500)"`
	RefreshToken string `gorm:"type:varchar(500)"`

	ExpiresIn int
	ExpiresAt int64
	Revoked   bool
}

// ToToken creates a Monzo Token from a WebSession
func (webSession WebSession) ToToken() *models.Token {
	return &models.Token{
		UserID:       webSession.UserID,
		ClientID:     webSession.ClientID,
		TokenType:    webSession.TokenType,
		AccessToken:  webSession.AccessToken,
		RefreshToken: webSession.RefreshToken,
		ExpiresIn:    webSession.ExpiresIn,
		ExpiresAt:    time.Unix(0, webSession.ExpiresAt),
	}
}

// Refresh refreshes a Monzo Token if it needs to be refreshed
func (webSession *WebSession) Refresh(connection *bongo.Connection, monzo *gomonzo.GoMonzo, session sessions.Session) (bool, *WebSession, error) {
	token, refreshed, _, err := monzo.RefreshAuthenticationIfNeeded(webSession.ToToken())
	if err != nil {
		return false, nil, err
	}

	if !refreshed {
		return false, webSession, nil
	}

	newWebSession := NewWebSession(connection, token, webSession.IP)
	session.Set("webSessionID", newWebSession.Id.Hex())
	session.Save()
	return false, newWebSession, nil
}

// NewWebSession creates a new session from a Monzo token and an IP address.
func NewWebSession(connection *bongo.Connection, token *models.Token, ip string) *WebSession {
	webSession := &WebSession{
		UserID:   token.UserID,
		ClientID: token.ClientID,
		IP:       ip,

		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,

		ExpiresIn: token.ExpiresIn,
		ExpiresAt: token.ExpiresAt.UnixNano(),
		Revoked:   false,
	}

	err := connection.Collection(WebSessionCollectionName).Save(webSession)
	if err != nil {
		raven.CaptureError(err, nil, nil)
	}
	return webSession
}

// GetWebSession ..
func GetWebSession(connection *bongo.Connection, monzo *gomonzo.GoMonzo, session sessions.Session, ip string) (*WebSession, error) {
	sessionID := session.Get("webSessionID")
	defer session.Save()
	if sessionID == nil {
		return nil, nil
	}

	var webSession *WebSession
	err := connection.Collection(WebSessionCollectionName).FindById(bson.ObjectIdHex(sessionID.(string)), &webSession)
	if err != nil {
		if dnfError, ok := err.(*bongo.DocumentNotFoundError); !ok {
			raven.CaptureError(dnfError, nil, nil)
		}
		return nil, nil
	}

	if webSession.Revoked || webSession.IP != ip {
		session.Delete("webSessionID")
		return nil, errors.New("session_revoked")
	}

	refreshed, newWebSession, err := webSession.Refresh(connection, monzo, session)
	if err != nil {
		session.Delete("webSessionID")
		return nil, err
	}
	if refreshed {
		return newWebSession, nil
	}

	return webSession, nil
}

package models

import "github.com/0xdeafcafe/gomonzo/models"

// Cookie contains links between a cookie and it's content
type Cookie struct {
	Audit

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

// NewCookie creates a new cookie from a Monzo token and an IP address.
func NewCookie(token *models.Token, ip string) *Cookie {
	cookie := &Cookie{
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

	cookie.Init()
	return cookie
}

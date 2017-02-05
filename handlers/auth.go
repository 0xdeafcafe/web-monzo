package handlers

import (
	"net/http"

	"github.com/0xdeafcafe/web-monzo/models"
	"github.com/gin-contrib/sessions"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

// AuthHandler ..
type AuthHandler struct {
	Context *models.Context
}

const authRedirect = "https://auth.getmondo.co.uk/?client_id=%s&redirect_uri=%sauth/callback&response_type=code&state=%s"

// Index ..
func (hndlr AuthHandler) Index(c *gin.Context) {
	session := sessions.Default(c)
	warning := session.Get("_warning")
	c.HTML(http.StatusOK, "auth/index", gin.H{
		"title":     "ay",
		"useLayout": "true",
		"warning":   warning,
	})
}

// Create handles creating a new monzo auth
func (hndlr AuthHandler) Create(c *gin.Context) {
	state := uuid.NewV4().String()
	url := hndlr.Context.Monzo.CreateAuthorizationURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback ..
func (hndlr AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")

	token, monzoError, err := hndlr.Context.Monzo.RequestAccessToken(code)
	if err != nil {
		session := sessions.Default(c)
		errorStr := ""
		if monzoError != nil {
			errorStr = monzoError.ErrorDescription
		} else {
			errorStr = err.Error()
		}
		session.Set("_warning", errorStr)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	cookie := models.NewCookie(token, c.ClientIP())
	hndlr.Context.DB.Create(&cookie)
	c.JSON(http.StatusOK, &cookie)
}

// NewAuthHandler creates a new AuthHandler and registers the reqired routes
func NewAuthHandler(e *gin.Engine, c *models.Context) {
	handler := new(AuthHandler)
	handler.Context = c

	e.GET("/auth", handler.Index)
	e.GET("/auth/create", handler.Create)
	e.GET("/auth/callback", handler.Callback)
}

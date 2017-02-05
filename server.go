package main

import (
	"html/template"
	"log"

	"fmt"

	"github.com/0xdeafcafe/web-monzo/handlers"
	"github.com/0xdeafcafe/web-monzo/helpers"
	"github.com/0xdeafcafe/web-monzo/models"

	"github.com/0xdeafcafe/gomonzo"
	monzoModels "github.com/0xdeafcafe/gomonzo/models"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/sessions"
	"github.com/jinzhu/configor"
	eztemplate "github.com/michelloworld/ez-gin-template"
	"gopkg.in/gin-gonic/gin.v1"
)

var config = struct {
	ConnectionString string `json:"connectionString"`
	ClientID         string `json:"clientId"`
	ClientSecret     string `json:"clientSecret"`
	ApplicationURL   string `json:"applicationUrl"`
	CookieSecret     string `json:"cookieSecret"`

	DSN string `json:"dsn"`
}{}

func main() {
	configor.Load(&config, "config.json")
	raven.SetDSN(config.DSN)

	// Create context
	context := &models.Context{
		Monzo: gomonzo.New(&monzoModels.MonzoOptions{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  fmt.Sprintf("%s/auth/callback", config.ApplicationURL),
		}),
		DB: helpers.NewDatabaseConnection(config.ConnectionString),
	}

	// Create Session Store
	store := sessions.NewCookieStore([]byte(config.CookieSecret))

	// Create renderer
	render := eztemplate.New()
	render.Debug = true
	render.TemplatesDir = "views/"
	render.Layout = "layouts/base"
	render.TemplateFuncMap = template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"notEq": func(a, b interface{}) bool {
			return a != b
		},
	}

	// Create Gin
	r := gin.Default()
	r.HTMLRender = render.Init()
	//r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.Use(sessions.Sessions("session", store))
	r.Static("/static", "./static/")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Register handlers
	handlers.NewHomeHandler(r, context)
	handlers.NewAuthHandler(r, context)

	// Listen for things
	log.Fatal(r.Run(":3000"))
}

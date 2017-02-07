package main

import (
	"html/template"
	"log"

	"fmt"

	"github.com/0xdeafcafe/web-monzo/handlers"
	"github.com/0xdeafcafe/web-monzo/helpers"
	"github.com/0xdeafcafe/web-monzo/models"

	"encoding/gob"

	"github.com/0xdeafcafe/gomonzo"
	monzoModels "github.com/0xdeafcafe/gomonzo/models"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-contrib/sessions"
	"github.com/jinzhu/configor"
	eztemplate "github.com/michelloworld/ez-gin-template"
	"gopkg.in/gin-gonic/gin.v1"
)

var config = struct {
	Mongo struct {
		ConnectionString string `json:"connection_string"`
		DatabaseName     string `json:"database_name"`
	} `json:"mongo"`

	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	ApplicationURL string `json:"application_url"`
	CookieSecret   string `json:"cookie_secret"`

	DSN string `json:"dsn"`
}{}

func main() {
	configor.Load(&config, "config.json")
	raven.SetDSN(config.DSN)

	// Create Gin
	r := gin.Default()
	if gin.Mode() == gin.ReleaseMode {
		r.Use(sentry.Recovery(raven.DefaultClient, false))
	}

	// Create context
	context := &models.Context{
		Monzo: gomonzo.New(&monzoModels.MonzoOptions{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  fmt.Sprintf("%s/auth/callback", config.ApplicationURL),
		}),
		Mongo: helpers.NewDatabaseConnection(config.Mongo.ConnectionString, config.Mongo.DatabaseName),
	}

	// Create Session Store
	store := sessions.NewCookieStore([]byte(config.CookieSecret))

	// Register Gob Types
	gob.Register(&models.Flash{})

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
		"greaterThan": func(a, b int64) bool {
			return a > b
		},
		"lessThanOrEq": func(a, b int64) bool {
			return a <= b
		},
	}

	// Set Renderer
	r.HTMLRender = render.Init()

	// Add Middleware
	r.Use(sessions.Sessions("session", store))

	// Register Static content
	r.Static("/static", "./static/")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Register handlers
	handlers.NewHomeHandler(r, context)
	handlers.NewAuthHandler(r, context)
	handlers.NewAccountsHandler(r, context)

	// Listen for things
	log.Fatal(r.Run(":3000"))
}

package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := flag.Int("port", 1323, "Default port of the service")
	redirectURL := flag.String("redirectUrl", "http://localhost:1323/auth", "Config: redirect URL")
	authURL := flag.String("authUrl", "http://localhost:8000/authorize", "Config: authorization URL")
	tokenURL := flag.String("tokenUrl", "http://localhost:8000/token", "Config: token URL")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	conf := &oauth2.Config{
		ClientID:     "test_client",
		ClientSecret: "aabbccdd",
		RedirectURL:  *redirectURL,
		Scopes: []string{
			"projects.view",
			"users.*",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  *authURL,
			TokenURL: *tokenURL,
		},
	}

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/connect", func(c echo.Context) error {
		url := conf.AuthCodeURL("")
		// TODO: implement interaction with login page / service
		// For now (simplicity) we know that oauth service
		// returns access code for every request

		return c.Redirect(http.StatusSeeOther, url)
	})

	e.GET("/auth", func(c echo.Context) error {
		code := c.QueryParam("code")

		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		fmt.Println("Access Token:", tok.AccessToken)
		// TODO: now we have token (tok), we can issue requests

		return c.String(http.StatusOK, "Connected\n")
	})

	// Start server
	e.Run(standard.New(fmt.Sprintf(":%d", *port)))
}

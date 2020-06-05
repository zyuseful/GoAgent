package main

import (
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	//Echo_Middleware "github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	/*
		e := echo.New()
		// e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("<DOMAIN>")
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())
		e.GET("/", func(c echo.Context) error {
			return c.HTML(http.StatusOK, `
				<h1>Welcome to Echo!</h1>
				<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
			`)
		})
		e.Logger.Fatal(e.StartAutoTLS(":443"))
	*/

	/*
	e := echo.New()
	//certFile, keyFile string
	certFile := "/Users/zys/go/src/myagent/uploadFiles/harbor-ca.crt"
	keyFile := "/Users/zys/go/src/myagent/uploadFiles/harbor-ca.key"

	//e.TLSServer.ListenAndServeTLS(certFile,keyFile)
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
				<h1>Welcome to Echo!</h1>
				<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
			`)
	})
	e.StartTLS(":8999",certFile,keyFile)
	 */

	ipAndPort := ":8090"
	e := echo.New()
	//MiddlewareFunc
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		fmt.Println(e)
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("as")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("as")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		fmt.Println(handlerFunc)
		return nil
	})

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
				<h1>Welcome to Echo!</h1>
				<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
			`)
	})
	e.Start(ipAndPort)
}
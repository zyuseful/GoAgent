package main

import (
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myagent/src/core/perception"
	"myagent/src/core/structure"

	//Echo_Middleware "github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main1() {
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

//func main() {
//	fmt.Println(2 ^ 1)
//}
func main2() {
	node := perception.CreatePNode()

	fmt.Printf("%b\n", node.GetState())
	node.SetPNodeState(false, true, true)
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

	node.SetPNodeState(true, true, false)
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

	node.SetThisPNodeActive()
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

	node.SetThisPNodeDeaded()
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

	node.SetCheckComePNode()
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

	node.SetNoCheckComePNode()
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

	node.SetCheckComePNodeActive()
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())
	node.SetCheckComePNodeDeaded()
	fmt.Println("是否检查 come :", node.CheckComePNode())
	fmt.Println("节点是否存活   :", node.CheckThisPNodeIsActive())
	fmt.Printf("%b\n", node.GetState())

}

func main3() {
	list := structure.ArrayList{}
	list.Add("A")
	list.Add("B")
	list.Add("C")
	list.Add("D")

	list.Print()
	fmt.Println(list.Size())

	list.AppendTo(1, "c", "d", "e", "f")
	list.Print()
	fmt.Println(list.Size())
	//
	//
	list.AppendTo(list.Size(), "c1", "d1", "e1", "f1")
	list.Print()
	fmt.Println(list.Size())

	list.AppendTo(list.Size()-1, "c11", "d11", "e11", "f11")
	//list.AppendTo(8,"c1","d1","e1","f1")
	list.Print()
	fmt.Println(list.Size())

	list.Add("AA", "AB", "AC", "AD")
	list.Print()
	fmt.Println(list.Size())
}

func main() {
	list := structure.ArrayList{}
	list.Add("A")
	list.Add("B")
	list.Add("C")
	list.Add("D")

	list.Remove(4)
	list.Print()
}

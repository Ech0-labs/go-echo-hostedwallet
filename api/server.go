package api

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ltcsuite/ltcd/rpcclient"
)

type Server struct {
	client *rpcclient.Client
	http   *echo.Echo
}

func InitClientFromEnv() (*rpcclient.Client, error) {
	host := os.Getenv("RPC_HOST")
	user := os.Getenv("RPC_USER")
	pass := os.Getenv("RPC_PASS")

	conf := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	return rpcclient.New(conf, nil)
}

func InitFromEnv() (*Server, error) {
	client, err := InitClientFromEnv()
	if err != nil {
		return nil, err
	}

	e := echo.New()

	e.GET("/echo", GetEcho(e, client))
	e.POST("/echo", PostEcho(e, client))

	e.GET("/addrs", GetAddr(e, client))

	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on port %s\n", port)
	address := fmt.Sprintf("localhost:%v", port)
	fmt.Println(address)

	go e.Logger.Fatal(e.Start(address))

	return &Server{client, e}, nil
}

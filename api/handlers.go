package api

import (
	"net/http"

	"github.com/Ech0-labs/go-echo-prototype/flow"
	"github.com/labstack/echo/v4"
	"github.com/ltcsuite/ltcd/rpcclient"
)

func GetEcho(e *echo.Echo, c *rpcclient.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "GET /echo")
	}
}

func PostEcho(e *echo.Echo, client *rpcclient.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		msg := c.FormValue("msg")
		if msg == "" {
			return c.String(http.StatusNoContent, "no message provided")
		}
		hash, err := flow.Send(client, msg)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, hash.String())
	}
}

package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Cat struct {
	Name	string	`json:"name"`
	Type	string	`json:"type"`
}

func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "yallo from the web side!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nand his type is: %s\n", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to lets us know if you want json or string data",
	})
}

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s",err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s",err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your cat: %#v", cat)
	return c.String(http.StatusOK, "we got you cat!")
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "horay you are on the secret admin main page!")
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "BlueBot/1.0")

		return next(c)
	}
}

func main() {
	fmt.Println("yallo for the first time")

	e := echo.New()

	e.Use(ServerHeader)

	g := e.Group("/admin")

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method} ${host}${path} ${latency)human}` + "\n",
	}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/main", mainAdmin)

	e.GET("/", yallo)
	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCat)

	e.Start(":8080")
}
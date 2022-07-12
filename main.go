package main

import (
	"mas-kusa-api/db"
	"mas-kusa-api/env"
	"mas-kusa-api/routers"
	"mas-kusa-api/schedules"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	env.InitEnv()
	if err := db.InitDB(); err != nil {
		return
	}
	if err := schedules.InitSchedule(); err != nil {
		return
	}

	db.Psql.AutoMigrate(&db.User{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")

	routers.SetRouter(e)

	e.Logger.Fatal(e.Start(":" + env.Port))
}

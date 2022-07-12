package routers

import (
	"fmt"
	"image/png"
	"mas-kusa-api/controllers"
	"mas-kusa-api/env"
	"mas-kusa-api/utils"
	"os"
	"strings"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: env.SigningKey,
}

func SetRouter(e *echo.Echo) {
	e.GET("", hello)
	e.POST("/users/signup", registerMastodonInfo)
	e.POST("/users/signin", generateJWT)
	auth := e.Group("")
	auth.Use(middleware.JWTWithConfig(Config))
	auth.GET("/generate", generateTempMasKusa)
}

func userInfoFromToken(c echo.Context) (string, string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return claims.Instance, claims.Token
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, SimpleJson{Message: "OK"})
}

func registerMastodonInfo(c echo.Context) error {
	r := new(Register)
	if err := c.Bind(r); err != nil {
		return err
	}
	if instance := r.Instance; instance[len(instance)-1:] == "/" {
		r.Instance = instance[:len(instance)-1]
	}
	if err := utils.CheckAccount(r.Instance, r.Token); err != nil {
		return c.JSON(http.StatusBadRequest, SimpleJson{err.Error()})
	}
	if acct, err := utils.GetUserName(r.Instance, r.Token); err != nil {
		return c.JSON(http.StatusBadRequest, SimpleJson{err.Error()})
	} else if err := controllers.AddUser(acct, r.Instance, r.Token); err != nil {
		return c.JSON(http.StatusBadRequest, SimpleJson{err.Error()})
	}
	return c.JSON(http.StatusOK, r)
}

func generateJWT(c echo.Context) error {
	u := new(Register)
	if err := c.Bind(u); err != nil {
		return err
	}
	if !controllers.IsAlreadyExistUser(u.Instance, u.Token) {
		return c.JSON(http.StatusNotFound, "this user is not found")
	}

	claims := &jwtCustomClaims{
		u.Instance,
		u.Token,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(env.SigningKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Token{Token: t})
}

func generateTempMasKusa(c echo.Context) error {
	var (
		instance string
		token    string
	)
	instance, token = userInfoFromToken(c)

	if _, err := controllers.GetUserInfo(token); err != nil {
		return err
	}

	now := time.Now()
	baseInstance := strings.Replace(strings.Replace(strings.Replace(instance, "https://", "", -1), "http://", "", -1), "/", "", -1)
	savingPath := "static/" + baseInstance + "/" + token[:15] + "-" + now.Format("20060102") + ".png"
	urlImagePath := env.ServerUrl + ":" + env.Port + "/" + savingPath
	if _, err := os.Stat(baseInstance); err != nil {
		os.Mkdir(baseInstance, 0777)
	}

	if _, err := os.Stat(savingPath); err == nil {
		fmt.Println("file is exist.")
		return c.JSON(http.StatusOK, ImagePath{Path: urlImagePath, Refresh: false})
	}

	baseDate, tootList := utils.CountToot(instance, token, true)

	weekDayNum := int(baseDate.Weekday())
	wholeTootCounter := [][7]int{}
	weekCount := [7]int{}
	for _, toot := range tootList {
		weekCount[weekDayNum] = toot
		weekDayNum++
		if weekDayNum >= 7 {
			wholeTootCounter = append(wholeTootCounter, weekCount)
			weekCount = [7]int{}
			weekDayNum = 0
		}
	}
	wholeTootCounter = append(wholeTootCounter, weekCount)

	baseImage := utils.GenKusa(wholeTootCounter)
	if imagePath, err := os.Create(savingPath); err != nil {
		return err
	} else {
		png.Encode(imagePath, baseImage)
	}
	return c.JSON(http.StatusOK, ImagePath{Path: urlImagePath, Refresh: true})
}
